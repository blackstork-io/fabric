package list

import (
	"bytes"
	"errors"
	"strings"
	"text/template"

	"github.com/Masterminds/semver/v3"
	"github.com/blackstork-io/fabric/plugininterface/v1"
	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/hcldec"
	"github.com/zclconf/go-cty/cty"
)

var Version = semver.MustParse("0.1.0")

const queryResultKey = "query_result"

type Plugin struct{}

type itemTempl = *template.Template

func (Plugin) GetPlugins() []plugininterface.Plugin {
	return []plugininterface.Plugin{
		{
			Namespace:  "blackstork",
			Kind:       "content",
			Name:       "list",
			Version:    plugininterface.Version(*Version),
			ConfigSpec: nil,
			InvocationSpec: &hcldec.ObjectSpec{
				"item_template": &hcldec.AttrSpec{
					Name:     "item_template",
					Type:     cty.String,
					Required: true,
				},
			},
		},
	}
}

func (p Plugin) Call(args plugininterface.Args) plugininterface.Result {
	tmpl, err := p.parseTemplate(args)
	if err != nil {
		return plugininterface.Result{
			Diags: hcl.Diagnostics{{
				Severity: hcl.DiagError,
				Summary:  "Failed to parse template",
				Detail:   err.Error(),
			}},
		}
	}
	result, err := p.render(tmpl, args.Context)
	if err != nil {
		return plugininterface.Result{
			Diags: hcl.Diagnostics{{
				Severity: hcl.DiagError,
				Summary:  "Failed to render template",
				Detail:   err.Error(),
			}},
		}
	}
	return plugininterface.Result{
		Result: result,
	}
}

func (p Plugin) parseTemplate(args plugininterface.Args) (itemTempl, error) {
	itemTemplate := args.Args.GetAttr("item_template")
	if itemTemplate.IsNull() {
		return nil, errors.New("item_template is required")
	}
	return template.New("item").Parse(itemTemplate.AsString())
}

func (p Plugin) render(tmpl itemTempl, datactx map[string]any) (string, error) {
	if datactx == nil {
		return "", errors.New("data context is required")
	}
	queryResult, ok := datactx[queryResultKey]
	if !ok || queryResult == nil {
		return "", errors.New("query_result is required in data context")
	}
	items, ok := queryResult.([]any)
	if !ok {
		return "", errors.New("query_result must be an array")
	}
	var buf bytes.Buffer
	for _, item := range items {
		tmpbuf := bytes.Buffer{}
		err := tmpl.Execute(&tmpbuf, item)
		if err != nil {
			return "", err
		}
		buf.WriteString("* ")
		buf.WriteString(strings.TrimSpace(strings.ReplaceAll(tmpbuf.String(), "\n", " ")))
		buf.WriteString("\n")
	}
	return buf.String(), nil
}
