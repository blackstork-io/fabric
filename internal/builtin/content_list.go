package builtin

import (
	"bytes"
	"context"
	"errors"
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

const (
	listQueryResultKey = "query_result"
)

func makeListContentProvider() *plugin.ContentProvider {
	return &plugin.ContentProvider{
		ContentFunc: genListContent,
		Args: dataspec.ObjectSpec{
			&dataspec.AttrSpec{
				Name:        "item_template",
				Type:        cty.String,
				Constraints: constraint.RequiredNonNull,
				ExampleVal:  cty.StringVal(`[{{.Title}}]({{.URL}})`),
				Doc:         "Go template for the item of the list",
			},
			&dataspec.AttrSpec{
				Name:       "format",
				Type:       cty.String,
				DefaultVal: cty.StringVal("unordered"),
				OneOf: []cty.Value{
					cty.StringVal("unordered"),
					cty.StringVal("ordered"),
					cty.StringVal("tasklist"),
				},
			},
		},
		Doc: "Produces a list of items",
	}
}

func genListContent(ctx context.Context, params *plugin.ProvideContentParams) (*plugin.ContentResult, diagnostics.Diag) {
	format, tmpl, err := parseListContentArgs(params)
	if err != nil {
		return nil, diagnostics.Diag{{
			Severity: hcl.DiagError,
			Summary:  "Failed to parse template",
			Detail:   err.Error(),
		}}
	}
	result, err := renderListContent(format, tmpl, params.DataContext)
	if err != nil {
		return nil, diagnostics.Diag{{
			Severity: hcl.DiagError,
			Summary:  "Failed to render template",
			Detail:   err.Error(),
		}}
	}
	return &plugin.ContentResult{
		Content: &plugin.ContentElement{
			Markdown: result,
		},
	}, nil
}

func parseListContentArgs(params *plugin.ProvideContentParams) (string, *template.Template, error) {
	itemTemplate := params.Args.GetAttr("item_template")
	format := params.Args.GetAttr("format").AsString()

	tmpl, err := template.New("item").Funcs(sprig.FuncMap()).Parse(itemTemplate.AsString())
	return format, tmpl, err
}

func renderListContent(format string, tmpl *template.Template, datactx plugin.MapData) (string, error) {
	if datactx == nil {
		return "", errors.New("data context is required")
	}
	queryResult, ok := datactx[listQueryResultKey]
	if !ok || queryResult == nil {
		return "", errors.New("query_result is required in data context")
	}
	items, ok := queryResult.(plugin.ListData)
	if !ok {
		return "", errors.New("query_result must be an array")
	}
	var buf bytes.Buffer
	for i, item := range items {
		tmpbuf := bytes.Buffer{}
		err := tmpl.Execute(&tmpbuf, item.Any())
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
