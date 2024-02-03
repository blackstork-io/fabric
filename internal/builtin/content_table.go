package builtin

import (
	"bytes"
	"context"
	"fmt"
	"strings"
	"text/template"

	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/hcldec"
	"github.com/zclconf/go-cty/cty"

	"github.com/blackstork-io/fabric/plugin"
)

type tableCellTmpl = *template.Template

func makeTableContentProvider() *plugin.ContentProvider {
	return &plugin.ContentProvider{
		ContentFunc: genTableContent,
		Args: &hcldec.ObjectSpec{
			"columns": &hcldec.AttrSpec{
				Name: "columns",
				Type: cty.List(cty.Object(map[string]cty.Type{
					"header": cty.String,
					"value":  cty.String,
				})),
				Required: true,
			},
		},
	}
}

func genTableContent(ctx context.Context, params *plugin.ProvideContentParams) (*plugin.Content, hcl.Diagnostics) {
	headers, values, err := parseTableContentArgs(params)
	if err != nil {
		return nil, hcl.Diagnostics{{
			Severity: hcl.DiagError,
			Summary:  "Failed to parse arguments",
			Detail:   err.Error(),
		}}
	}
	result, err := renderTableContent(headers, values, params.DataContext)
	if err != nil {
		return nil, hcl.Diagnostics{{
			Severity: hcl.DiagError,
			Summary:  "Failed to render table",
			Detail:   err.Error(),
		}}
	}
	return &plugin.Content{
		Markdown: result,
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
			ok     = false
		)
		if header, ok = obj["header"]; !ok || header.IsNull() {
			return nil, nil, fmt.Errorf("missing header in table cell")
		}
		if value, ok = obj["value"]; !ok || value.IsNull() {
			return nil, nil, fmt.Errorf("missing value in table cell")
		}

		headerTmpl, err := template.New("header").Parse(header.AsString())
		if err != nil {
			return nil, nil, fmt.Errorf("failed to parse header template: %w", err)
		}
		valueTmpl, err := template.New("value").Parse(value.AsString())
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
		err := header.Execute(&buf, datactx)
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
				err := value.Execute(&buf, row)
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
		buf.WriteString("-")
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
