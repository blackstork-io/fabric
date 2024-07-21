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

func makeListContentProvider() *plugin.ContentProvider {
	return &plugin.ContentProvider{
		ContentFunc: genListContent,
		Args: &dataspec.RootSpec{
			Attrs: []*dataspec.AttrSpec{
				{
					Name:        "item_template",
					Type:        cty.String,
					Constraints: constraint.NonNull,
					DefaultVal:  cty.StringVal("{{.}}"),
					ExampleVal:  cty.StringVal(`[{{.Title}}]({{.URL}})`),
					Doc:         "Go template for the item of the list",
				},
				{
					Name:       "format",
					Type:       cty.String,
					DefaultVal: cty.StringVal("unordered"),
					OneOf: []cty.Value{
						cty.StringVal("unordered"),
						cty.StringVal("ordered"),
						cty.StringVal("tasklist"),
					},
				},
				{
					Name:        "items",
					Type:        dataquery.DelayedEvalType.CtyType(),
					Constraints: constraint.RequiredMeaningful,
					ExampleVal: cty.ListVal([]cty.Value{
						cty.StringVal("First item"),
						cty.StringVal("Second item"),
						cty.StringVal("Third item"),
					}),
					Doc: "List of items to render.",
				},
			},
		},
		Doc: "Produces a list of items",
	}
}

func genListContent(ctx context.Context, params *plugin.ProvideContentParams) (*plugin.ContentResult, diagnostics.Diag) {
	format := params.Args.GetAttr("format").AsString()

	tmpl, err := template.New("item").Funcs(sprig.FuncMap()).Parse(params.Args.GetAttr("item_template").AsString())
	if err != nil {
		return nil, diagnostics.Diag{{
			Severity: hcl.DiagError,
			Summary:  "Failed to parse template",
			Detail:   err.Error(),
		}}
	}

	items := dataquery.DelayedEvalType.MustFromCty(params.Args.GetAttr("items")).Result()
	if items == nil {
		return nil, diagnostics.Diag{{
			Severity: hcl.DiagError,
			Summary:  "Failed to parse arguments",
			Detail:   "Data is nil",
		}}
	}
	itemsList, ok := items.(plugin.ListData)
	if !ok {
		return nil, diagnostics.Diag{{
			Severity: hcl.DiagError,
			Summary:  "Failed to parse arguments",
			Detail:   "Data must be a list",
		}}
	}

	result, err := renderListContent(format, tmpl, itemsList)
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

func renderListContent(format string, tmpl *template.Template, items plugin.ListData) (string, error) {
	var buf bytes.Buffer
	var tmpBuf bytes.Buffer
	for i, item := range items {
		tmpBuf.Reset()
		err := tmpl.Execute(&tmpBuf, item.Any())
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
		buf.Write(bytes.TrimSpace(bytes.ReplaceAll(tmpBuf.Bytes(), []byte("\n"), []byte(" "))))
		buf.WriteString("\n")
	}
	return buf.String(), nil
}
