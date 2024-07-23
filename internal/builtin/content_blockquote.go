package builtin

import (
	"context"
	"strings"

	"github.com/hashicorp/hcl/v2"
	"github.com/zclconf/go-cty/cty"

	"github.com/blackstork-io/fabric/pkg/diagnostics"
	"github.com/blackstork-io/fabric/plugin"
	"github.com/blackstork-io/fabric/plugin/dataspec"
	"github.com/blackstork-io/fabric/plugin/dataspec/constraint"
)

func makeBlockQuoteContentProvider() *plugin.ContentProvider {
	return &plugin.ContentProvider{
		ContentFunc: genBlockQuoteContent,
		Args: &dataspec.RootSpec{
			Attrs: []*dataspec.AttrSpec{{
				Name:        "value",
				Type:        cty.String,
				ExampleVal:  cty.StringVal("Text to be formatted as a quote"),
				Constraints: constraint.RequiredNonNull,
			}},
		},
		Doc: "Formats text as a block quote",
	}
}

func genBlockQuoteContent(ctx context.Context, params *plugin.ProvideContentParams) (*plugin.ContentResult, diagnostics.Diag) {
	value := params.Args.GetAttrVal("value")
	text, err := genTextContentText(value.AsString(), params.DataContext)
	if err != nil {
		return nil, diagnostics.Diag{{
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
