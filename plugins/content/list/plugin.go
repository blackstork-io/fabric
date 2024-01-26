package list

import (
	"bytes"
	"errors"
	"fmt"
	"slices"
	"strings"
	"text/template"

	"github.com/Masterminds/semver/v3"
	"github.com/blackstork-io/fabric/plugininterface/v1"
	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/hcldec"
	"github.com/zclconf/go-cty/cty"
)

var Version = semver.MustParse("0.1.0")

const (
	queryResultKey = "query_result"
	defaultFormat  = "unordered"
)

var allowedFormats = []string{"unordered", "ordered", "tasklist"}

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
				"format": &hcldec.AttrSpec{
					Name:     "format",
					Type:     cty.String,
					Required: false,
				},
			},
		},
	}
}

func (p Plugin) Call(args plugininterface.Args) plugininterface.Result {
	format, tmpl, err := p.parseArgs(args)
	if err != nil {
		return plugininterface.Result{
			Diags: hcl.Diagnostics{{
				Severity: hcl.DiagError,
				Summary:  "Failed to parse template",
				Detail:   err.Error(),
			}},
		}
	}
	result, err := p.render(format, tmpl, args.Context)
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

func (p Plugin) parseArgs(args plugininterface.Args) (string, itemTempl, error) {
	itemTemplate := args.Args.GetAttr("item_template")
	if itemTemplate.IsNull() {
		return "", nil, errors.New("item_template is required")
	}
	format := args.Args.GetAttr("format")
	if format.IsNull() {
		format = cty.StringVal(defaultFormat)
	}
	if !slices.Contains(allowedFormats, format.AsString()) {
		return "", nil, errors.New("invalid format: " + format.AsString())
	}
	tmpl, err := template.New("item").Parse(itemTemplate.AsString())
	return format.AsString(), tmpl, err
}

func (p Plugin) render(format string, tmpl itemTempl, datactx map[string]any) (string, error) {
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
	for i, item := range items {
		tmpbuf := bytes.Buffer{}
		err := tmpl.Execute(&tmpbuf, item)
		if err != nil {
			return "", err
		}
		if format == "unordered" {
			buf.WriteString("* ")
		} else if format == "tasklist" {
			buf.WriteString("* [ ] ")
		} else {
			fmt.Fprintf(&buf, "%d. ", i+1)
		}
		buf.WriteString(strings.TrimSpace(strings.ReplaceAll(tmpbuf.String(), "\n", " ")))
		buf.WriteString("\n")
	}
	return buf.String(), nil
}
