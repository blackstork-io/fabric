package builtin

import (
	"context"
	"strings"

	"github.com/hashicorp/hcl/v2"
	"github.com/zclconf/go-cty/cty"

	"github.com/blackstork-io/fabric/plugin"
	"github.com/blackstork-io/fabric/plugin/dataspec"
)

func makeBlockQuoteContentProvider() *plugin.ContentProvider {
	return &plugin.ContentProvider{
		ContentFunc: genBlockQuoteContent,
		Args: dataspec.ObjectSpec{
			&dataspec.AttrSpec{
				Name:       "value",
				Type:       cty.String,
				ExampleVal: cty.StringVal("Text to be formatted as a quote"),
				Required:   true,
			},
		},
		Doc: "Formats text as a block quote",
	}
}

func genBlockQuoteContent(ctx context.Context, params *plugin.ProvideContentParams) (*plugin.ContentResult, hcl.Diagnostics) {
	value := params.Args.GetAttr("value").AsString()
	text, err := genTextContentText(value, params.DataContext)
	if err != nil {
		return nil, hcl.Diagnostics{{
			Severity: hcl.DiagError,
			Summary:  "Failed to render blockquote",
			Detail:   err.Error(),
		}}
	}
	text = "> " + strings.ReplaceAll(text, "\n", "\n> ")
	return &plugin.ContentResult{
		Content: &plugin.ContentElement{
			Markdown: text,
		},
	}, nil
}
