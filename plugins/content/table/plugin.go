package table

import (
	"bytes"
	"errors"
	"fmt"
	"strings"
	"text/template"

	"github.com/Masterminds/semver/v3"
	"github.com/blackstork-io/fabric/plugininterface/v1"
	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/hcldec"
	"github.com/zclconf/go-cty/cty"
)

var Version = semver.MustParse("0.1.0")

type Plugin struct{}

type cellTmpl = *template.Template

func (Plugin) GetPlugins() []plugininterface.Plugin {
	return []plugininterface.Plugin{
		{
			Namespace:  "blackstork",
			Kind:       "content",
			Name:       "table",
			Version:    plugininterface.Version(*Version),
			ConfigSpec: nil,
			InvocationSpec: &hcldec.ObjectSpec{
				"columns": &hcldec.AttrSpec{
					Name: "columns",
					Type: cty.List(cty.Object(map[string]cty.Type{
						"header": cty.String,
						"value":  cty.String,
					})),
					Required: true,
				},
			},
		},
	}
}

func (p Plugin) Call(args plugininterface.Args) plugininterface.Result {
	headers, values, err := p.parseArgs(args)
	if err != nil {
		return plugininterface.Result{
			Diags: hcl.Diagnostics{{
				Severity: hcl.DiagError,
				Summary:  "Failed to parse arguments",
				Detail:   err.Error(),
			}},
		}
	}
	result, err := p.render(headers, values, args.Context)
	if err != nil {
		return plugininterface.Result{
			Diags: hcl.Diagnostics{{
				Severity: hcl.DiagError,
				Summary:  "Failed to render table",
				Detail:   err.Error(),
			}},
		}
	}
	return plugininterface.Result{
		Result: result,
	}
}

func (p Plugin) parseArgs(args plugininterface.Args) (headers []cellTmpl, values []cellTmpl, err error) {
	arr := args.Args.GetAttr("columns")
	if arr.IsNull() {
		return nil, nil, errors.New("columns is required")
	}
	if len(arr.AsValueSlice()) == 0 {
		return nil, nil, errors.New("columns must not be empty")
	}
	for _, val := range arr.AsValueSlice() {
		obj := val.AsValueMap()
		var (
			header cty.Value
			value  cty.Value
			ok     = false
		)
		if header, ok = obj["header"]; !ok || header.IsNull() {
			return nil, nil, errors.New("missing header in table cell")
		}
		if value, ok = obj["value"]; !ok || value.IsNull() {
			return nil, nil, errors.New("missing value in table cell")
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

func (p Plugin) render(headers, values []cellTmpl, datactx map[string]any) (string, error) {
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
		return "", errors.New("data context is nil")
	}
	if queryResult, ok := datactx["query_result"]; ok && queryResult != nil {
		queryResult, ok := queryResult.([]any)
		if !ok {
			return "", errors.New("query_result is not an array")
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
