package builtin

import (
	"bytes"
	"context"
	"fmt"
	"strings"
	"text/template"

	"github.com/Masterminds/sprig/v3"
	"github.com/hashicorp/hcl/v2"
	"github.com/zclconf/go-cty/cty"

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
				Constraints: constraint.RequiredNonNull,
			},
		},
		Doc: `
			Produces a table.

			This content provider assumes that ` + "`query_result`" + ` is a list of objects representing rows,
			and uses the configured ` + "`value`" + ` go templates (see below) to display each row.

			NOTE: ` + "`header`" + ` templates are executed with the whole context availible, while ` + "`value`" + `
			templates are executed on each item of the ` + "`query_result`" + ` list.
		`,
	}
}

func genTableContent(ctx context.Context, params *plugin.ProvideContentParams) (*plugin.ContentResult, diagnostics.Diag) {
	headers, values, err := parseTableContentArgs(params)
	if err != nil {
		return nil, diagnostics.Diag{{
			Severity: hcl.DiagError,
			Summary:  "Failed to parse arguments",
			Detail:   err.Error(),
		}}
	}
	result, err := renderTableContent(headers, values, params.DataContext)
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
	if arr.IsNull() {
		return nil, nil, fmt.Errorf("columns is required")
	}
	if len(arr.AsValueSlice()) == 0 {
		return nil, nil, fmt.Errorf("columns must not be empty")
	}
	for _, val := range arr.AsValueSlice() {
		obj := val.AsValueMap()
		var (
			header cty.Value
			value  cty.Value
			ok     bool
		)
		if header, ok = obj["header"]; !ok || header.IsNull() {
			return nil, nil, fmt.Errorf("missing header in table cell")
		}
		if value, ok = obj["value"]; !ok || value.IsNull() {
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

func renderTableContent(headers, values []tableCellTmpl, datactx plugin.MapData) (string, error) {
	hstr := make([]string, len(headers))
	vstr := [][]string{}
	for i, header := range headers {
		var buf bytes.Buffer
		err := header.Execute(&buf, datactx.Any())
		if err != nil {
			return "", fmt.Errorf("failed to render header: %w", err)
		}
		hstr[i] = strings.TrimSpace(
			strings.ReplaceAll(buf.String(), "\n", " "),
		)
	}
	if datactx == nil {
		return "", fmt.Errorf("data context is nil")
	}
	if queryResult, ok := datactx["query_result"]; ok && queryResult != nil {
		queryResult, ok := queryResult.(plugin.ListData)
		if !ok {
			return "", fmt.Errorf("query_result is not an array")
		}
		for _, row := range queryResult {
			rowstr := make([]string, len(values))
			for i, value := range values {
				var buf bytes.Buffer
				err := value.Execute(&buf, row.Any())
				if err != nil {
					return "", fmt.Errorf("failed to render value: %w", err)
				}
				rowstr[i] = strings.TrimSpace(
					strings.ReplaceAll(buf.String(), "\n", " "),
				)
			}
			vstr = append(vstr, rowstr)
		}
	}
	var buf bytes.Buffer
	buf.WriteByte('|')
	for _, header := range hstr {
		buf.WriteString(header)
		buf.WriteByte('|')
	}
	buf.WriteByte('\n')
	buf.WriteByte('|')
	for range hstr {
		buf.WriteString("---")
		buf.WriteByte('|')
	}
	buf.WriteByte('\n')
	for _, row := range vstr {
		buf.WriteByte('|')
		for _, value := range row {
			buf.WriteString(value)
			buf.WriteByte('|')
		}
		buf.WriteByte('\n')
	}
	return buf.String(), nil
}
