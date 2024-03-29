package builtin

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"slices"
	"strings"
	"text/template"

	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/hcldec"
	"github.com/zclconf/go-cty/cty"

	"github.com/blackstork-io/fabric/plugin"
)

const (
	listQueryResultKey = "query_result"
	listDefaultFormat  = "unordered"
)

var listAllowedFormats = []string{"unordered", "ordered", "tasklist"}

func makeListContentProvider() *plugin.ContentProvider {
	return &plugin.ContentProvider{
		ContentFunc: genListContent,
		Args: hcldec.ObjectSpec{
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
	}
}

func genListContent(ctx context.Context, params *plugin.ProvideContentParams) (*plugin.ContentResult, hcl.Diagnostics) {
	format, tmpl, err := parseListContentArgs(params)
	if err != nil {
		return nil, hcl.Diagnostics{{
			Severity: hcl.DiagError,
			Summary:  "Failed to parse template",
			Detail:   err.Error(),
		}}
	}
	result, err := renderListContent(format, tmpl, params.DataContext)
	if err != nil {
		return nil, hcl.Diagnostics{{
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
	if itemTemplate.IsNull() {
		return "", nil, errors.New("item_template is required")
	}
	format := params.Args.GetAttr("format")
	if format.IsNull() {
		format = cty.StringVal(listDefaultFormat)
	}
	if !slices.Contains(listAllowedFormats, format.AsString()) {
		return "", nil, errors.New("invalid format: " + format.AsString())
	}
	tmpl, err := template.New("item").Parse(itemTemplate.AsString())
	return format.AsString(), tmpl, err
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
