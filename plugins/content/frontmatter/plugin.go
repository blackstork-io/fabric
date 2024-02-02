package frontmatter

import (
	"encoding/json"
	"errors"
	"fmt"
	"slices"
	"strings"

	"github.com/Masterminds/semver/v3"
	"github.com/blackstork-io/fabric/plugininterface/v1"
	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/hcldec"
	"github.com/pelletier/go-toml/v2"
	"github.com/zclconf/go-cty/cty"
	"gopkg.in/yaml.v3"
)

var Version = semver.MustParse("0.1.0")

const (
	defaultFormat  = "yaml"
	queryResultKey = "query_result"
)

var allowedFormats = []string{"yaml", "toml", "json"}

type Plugin struct{}

func (Plugin) GetPlugins() []plugininterface.Plugin {
	return []plugininterface.Plugin{
		{
			Namespace:  "blackstork",
			Kind:       "content",
			Name:       "frontmatter",
			Version:    plugininterface.Version(*Version),
			ConfigSpec: nil,
			InvocationSpec: &hcldec.ObjectSpec{
				"format": &hcldec.AttrSpec{
					Name:     "format",
					Type:     cty.String,
					Required: false,
				},
				"content": &hcldec.AttrSpec{
					Name:     "content",
					Type:     cty.Map(cty.DynamicPseudoType),
					Required: false,
				},
			},
		},
	}
}

func (p Plugin) Call(args plugininterface.Args) plugininterface.Result {
	format, m, err := p.parseArgs(args.Args, args.Context)
	if err != nil {
		return plugininterface.Result{
			Diags: hcl.Diagnostics{{
				Severity: hcl.DiagError,
				Summary:  "Failed to parse arguments",
				Detail:   err.Error(),
			}},
		}
	}
	result, err := p.render(format, m)
	if err != nil {
		return plugininterface.Result{
			Diags: hcl.Diagnostics{{
				Severity: hcl.DiagError,
				Summary:  "Failed to render frontmatter",
				Detail:   err.Error(),
			}},
		}
	}
	return plugininterface.Result{
		Result: result,
	}
}

func (p Plugin) parseArgs(args cty.Value, datactx map[string]any) (string, map[string]any, error) {
	format := args.GetAttr("format")
	if format.IsNull() || format.AsString() == "" {
		format = cty.StringVal(defaultFormat)
	}
	if !slices.Contains(allowedFormats, format.AsString()) {
		return "", nil, fmt.Errorf("invalid format: %s", format.AsString())
	}
	var m map[string]any
	if datactx != nil {
		if queryResult, ok := datactx[queryResultKey]; ok {
			if qr, ok := queryResult.(map[string]any); ok {
				m = qr
			} else {
				return "", nil, fmt.Errorf("invalid query result: %T", queryResult)
			}
		}
	}
	content := args.GetAttr("content")
	if m == nil {
		if !content.IsNull() {
			m = p.convertMap(content)
		} else {
			return "", nil, errors.New("query_result and content are nil")
		}
	}
	return format.AsString(), m, nil
}

func (p Plugin) render(format string, m map[string]any) (string, error) {
	switch format {
	case "yaml":
		return p.renderYAML(m)
	case "toml":
		return p.renderTOML(m)
	case "json":
		return p.renderJSON(m)
	default:
		return "", fmt.Errorf("invalid format: %s", format)
	}
}

func (p Plugin) renderYAML(m map[string]any) (string, error) {
	var buf strings.Builder
	buf.WriteString("---\n")
	err := yaml.NewEncoder(&buf).Encode(m)
	if err != nil {
		return "", err
	}
	buf.WriteString("---\n")
	return buf.String(), nil
}

func (p Plugin) renderTOML(m map[string]any) (string, error) {
	var buf strings.Builder
	buf.WriteString("+++\n")
	err := toml.NewEncoder(&buf).Encode(m)
	if err != nil {
		return "", err
	}
	buf.WriteString("+++\n")
	return buf.String(), nil
}

func (p Plugin) renderJSON(m map[string]any) (string, error) {
	var buf strings.Builder
	enc := json.NewEncoder(&buf)
	enc.SetIndent("", "  ")
	err := enc.Encode(m)
	if err != nil {
		return "", err
	}
	return buf.String(), nil
}

func (p Plugin) convert(v cty.Value) any {
	if v.IsNull() {
		return nil
	}
	t := v.Type()
	switch {
	case t == cty.String:
		return v.AsString()
	case t == cty.Number:
		if v.AsBigFloat().IsInt() {
			n, _ := v.AsBigFloat().Int64()
			return n
		} else {
			n, _ := v.AsBigFloat().Float64()
			return n
		}
	case t == cty.Bool:
		return v.True()
	case t.IsMapType() || t.IsObjectType():
		return p.convertMap(v)
	case t.IsListType():
		return p.convertList(v)
	default:
		return nil
	}
}

func (p Plugin) convertList(v cty.Value) []any {
	var result []any
	for _, v := range v.AsValueSlice() {
		result = append(result, p.convert(v))
	}
	return result
}

func (p Plugin) convertMap(v cty.Value) map[string]any {
	result := make(map[string]any)
	for k, v := range v.AsValueMap() {
		result[k] = p.convert(v)
	}
	return result
}
