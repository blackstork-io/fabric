package builtin

import (
	"bytes"
	"context"
	"fmt"
	"text/template"

	"github.com/Masterminds/sprig/v3"
	"github.com/hashicorp/hcl/v2"
	"github.com/zclconf/go-cty/cty"

	"github.com/blackstork-io/fabric/pkg/diagnostics"
	"github.com/blackstork-io/fabric/pkg/utils"
	"github.com/blackstork-io/fabric/plugin"
	"github.com/blackstork-io/fabric/plugin/dataspec"
	"github.com/blackstork-io/fabric/plugin/dataspec/constraint"
	"github.com/blackstork-io/fabric/plugin/plugindata"
)

type tableCellTmpl = *template.Template

func makeTableContentProvider() *plugin.ContentProvider {
	return &plugin.ContentProvider{
		ContentFunc: genTableContent,
		Args: &dataspec.RootSpec{
			Attrs: []*dataspec.AttrSpec{
				{
					Name: "rows",
					Type: cty.List(plugindata.Encapsulated.CtyType()),
					Doc: "A list of objects representing rows in the table.\n" +
						"May be set statically or as a result of one or more queries.",
				},
				{
					Name: "columns",
					Type: cty.List(cty.Object(map[string]cty.Type{
						"header": cty.String,
						"value":  cty.String,
					})),
					Doc: `List of header and value go templates for each column`,
					ExampleVal: cty.ListVal([]cty.Value{
						cty.ObjectVal(map[string]cty.Value{
							"header": cty.StringVal("1st column header template"),
							"value":  cty.StringVal("1st column values template"),
						}),
						cty.ObjectVal(map[string]cty.Value{
							"header": cty.StringVal("2nd column header template"),
							"value":  cty.StringVal("2nd column values template"),
						}),
						cty.ObjectVal(map[string]cty.Value{
							"header": cty.StringVal("..."),
							"value":  cty.StringVal("..."),
						}),
					}),
					Constraints: constraint.RequiredMeaningful,
				},
			},
		},
		Doc: `
			Produces a table.

			Each cell template has access to the data context and the following variables:
			* ` + "`.rows` – the value of `rows` argument" + `
			* ` + "`.row.value` – the current row from `.rows` list" + `
			* ` + "`.row.index` – the current row index" + `
			* ` + "`.col.index` – the current column index" + `

			Header templates have access to the same variables as value templates,
			except for ` + "`.row.value` and `.row.index`",
	}
}

func genTableContent(ctx context.Context, params *plugin.ProvideContentParams) (*plugin.ContentResult, diagnostics.Diag) {
	var rows plugindata.List
	rowsVal := params.Args.GetAttrVal("rows")
	if !rowsVal.IsNull() {
		var err error
		rows, err = utils.FnMapErr(rowsVal.AsValueSlice(), func(v cty.Value) (plugindata.Data, error) {
			data, err := plugindata.Encapsulated.FromCty(v)
			if err != nil {
				return nil, err
			}
			if data == nil {
				return nil, nil
			}
			return *data, nil
		})
		if err != nil {
			return nil, diagnostics.Diag{{
				Severity: hcl.DiagError,
				Summary:  "Failed to parse arguments",
				Detail:   err.Error(),
				Subject:  &params.Args.Attrs["rows"].ValueRange,
			}}
		}
	}

	headers, values, err := parseTableContentArgs(params)
	if err != nil {
		return nil, diagnostics.Diag{{
			Severity: hcl.DiagError,
			Summary:  "Failed to parse arguments",
			Detail:   err.Error(),
		}}
	}
	result, err := renderTableContent(headers, values, params.DataContext, rows)
	if err != nil {
		return nil, diagnostics.Diag{{
			Severity: hcl.DiagError,
			Summary:  "Failed to render table",
			Detail:   err.Error(),
		}}
	}
	return &plugin.ContentResult{
		Content: plugin.NewElementFromMarkdown(result),
	}, nil
}

func parseTableContentArgs(params *plugin.ProvideContentParams) (headers, values []tableCellTmpl, err error) {
	arr := params.Args.GetAttrVal("columns")
	for _, val := range arr.AsValueSlice() {
		obj := val.AsValueMap()
		header := obj["header"]
		if header.IsNull() {
			return nil, nil, fmt.Errorf("missing header in table cell")
		}
		value := obj["value"]
		if value.IsNull() {
			return nil, nil, fmt.Errorf("missing value in table cell")
		}

		headerTmpl, err := template.New("header").Funcs(sprig.FuncMap()).Parse(header.AsString())
		if err != nil {
			return nil, nil, fmt.Errorf("failed to parse header template: %w", err)
		}
		valueTmpl, err := template.New("value").Funcs(sprig.FuncMap()).Parse(value.AsString())
		if err != nil {
			return nil, nil, fmt.Errorf("failed to parse value template: %w", err)
		}
		headers = append(headers, headerTmpl)
		values = append(values, valueTmpl)
	}
	return
}

func renderTableContent(headers, values []tableCellTmpl, dataCtx plugindata.Map, rowsList plugindata.List) (string, error) {
	var buf bytes.Buffer

	data := dataCtx.Any().(map[string]any)

	rows := rowsList.Any().([]any)
	data["rows"] = rows
	col := map[string]any{}
	data["col"] = col
	buf.WriteByte('|')
	var cellBuf bytes.Buffer
	for colIdx, header := range headers {
		cellBuf.Reset()
		col["index"] = colIdx + 1
		err := header.Execute(&cellBuf, data)
		if err != nil {
			return "", fmt.Errorf("failed to render header: %w", err)
		}

		buf.Write(
			bytes.ReplaceAll(
				bytes.TrimSpace(cellBuf.Bytes()),
				[]byte("\n"),
				[]byte(" "),
			),
		)
		buf.WriteByte('|')
	}
	buf.WriteByte('\n')
	buf.WriteByte('|')
	for range headers {
		buf.WriteString("---|")
	}
	buf.WriteString("\n")

	dataRow := map[string]any{}
	data["row"] = dataRow

	for rowIdx, row := range rows {
		buf.WriteByte('|')
		dataRow["index"] = rowIdx + 1
		dataRow["value"] = row
		for colIdx, value := range values {
			cellBuf.Reset()
			col["index"] = colIdx + 1
			err := value.Execute(&cellBuf, data)
			if err != nil {
				return "", fmt.Errorf("failed to render value: %w", err)
			}
			buf.Write(
				bytes.ReplaceAll(
					bytes.TrimSpace(cellBuf.Bytes()),
					[]byte("\n"),
					[]byte(" "),
				),
			)
			buf.WriteByte('|')
		}
		buf.WriteByte('\n')
	}

	return buf.String(), nil
}
