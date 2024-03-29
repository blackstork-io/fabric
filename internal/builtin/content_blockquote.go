package builtin

import (
	"context"
	"strings"

	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/hcldec"
	"github.com/zclconf/go-cty/cty"

	"github.com/blackstork-io/fabric/plugin"
)

func makeBlockQuoteContentProvider() *plugin.ContentProvider {
	return &plugin.ContentProvider{
		ContentFunc: genBlockQuoteContent,
		Args: hcldec.ObjectSpec{
			"value": &hcldec.AttrSpec{
				Name:     "value",
				Type:     cty.String,
				Required: true,
			},
		},
	}
}

func genBlockQuoteContent(ctx context.Context, params *plugin.ProvideContentParams) (*plugin.ContentResult, hcl.Diagnostics) {
	value := params.Args.GetAttr("value")
	if value.IsNull() {
		return nil, hcl.Diagnostics{{
			Severity: hcl.DiagError,
			Summary:  "Failed to parse arguments",
			Detail:   "value is required",
		}}
	}
	text, err := genTextContentText(value.AsString(), params.DataContext)
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
