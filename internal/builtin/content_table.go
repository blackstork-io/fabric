package builtin

import (
	"bytes"
	"context"
	"fmt"
	"text/template"

	"github.com/Masterminds/sprig/v3"
	"github.com/hashicorp/hcl/v2"
	"github.com/zclconf/go-cty/cty"

	"github.com/blackstork-io/fabric/eval/dataquery"
	"github.com/blackstork-io/fabric/pkg/diagnostics"
	"github.com/blackstork-io/fabric/plugin"
	"github.com/blackstork-io/fabric/plugin/dataspec"
	"github.com/blackstork-io/fabric/plugin/dataspec/constraint"
)

type tableCellTmpl = *template.Template

func makeTableContentProvider() *plugin.ContentProvider {
	return &plugin.ContentProvider{
		ContentFunc: genTableContent,
		Args: dataspec.ObjectSpec{
			&dataspec.AttrSpec{
				Name: "rows_var",
				Type: dataquery.DelayedEvalType.CtyType(),
				Doc: "A list of objects representing rows in the table.\n" +
					"May be set statically or as a result of one or more queries.",
			},
			&dataspec.AttrSpec{
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
		Doc: `
			Produces a table.

			Each cell template has access to the data context and the following variables:
			* ` + "`.rows` – the value of `rows_var` attribute" + `
			* ` + "`.row.value` – the current row from `.rows` list" + `
			* ` + "`.row.index` – the current row index" + `
			* ` + "`.col.index` – the current column index" + `

			Header templates have access to the same variables as value templates,
			except for ` + "`.row.value` and `.row.index`",
	}
}

func genTableContent(ctx context.Context, params *plugin.ProvideContentParams) (*plugin.ContentResult, diagnostics.Diag) {
	var rows_var plugin.ListData
	rows_val := params.Args.GetAttr("rows_var")
	if !rows_val.IsNull() {
		res, err := dataquery.DelayedEvalType.FromCty(rows_val)
		if err != nil {
			return nil, diagnostics.FromErr(err, "failed to get rows_var")
		}
		data := res.Result()
		var ok bool
		if data != nil {
			rows_var, ok = data.(plugin.ListData)
			if !ok {
				return nil, diagnostics.Diag{{
					Severity: hcl.DiagError,
					Summary:  "Failed to parse arguments",
					Detail:   fmt.Sprintf("rows_var must be a list, not %T", data),
				}}
			}
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
	result, err := renderTableContent(headers, values, params.DataContext, rows_var)
	if err != nil {
		return nil, diagnostics.Diag{{
			Severity: hcl.DiagError,
			Summary:  "Failed to render table",
			Detail:   err.Error(),
		}}
	}
	return &plugin.ContentResult{
		Content: &plugin.ContentElement{
			Markdown: result,
		},
	}, nil
}

func parseTableContentArgs(params *plugin.ProvideContentParams) (headers, values []tableCellTmpl, err error) {
	arr := params.Args.GetAttr("columns")
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

func renderTableContent(headers, values []tableCellTmpl, dataCtx plugin.MapData, rows_var plugin.ListData) (string, error) {
	var buf bytes.Buffer

	data := dataCtx.Any().(map[string]any)

	rows := rows_var.Any().([]any)
	data["rows"] = rows
	col := map[string]any{}
	data["col"] = col
	buf.WriteByte('|')
	var cellBuf bytes.Buffer
	for col_idx, header := range headers {
		cellBuf.Reset()
		col["index"] = col_idx + 1
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

	for row_idx, row := range rows {
		buf.WriteByte('|')
		dataRow["index"] = row_idx + 1
		dataRow["value"] = row
		for col_idx, value := range values {
			cellBuf.Reset()
			col["index"] = col_idx + 1
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
