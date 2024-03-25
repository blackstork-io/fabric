package builtin

import (
	"context"
	"fmt"

	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/hcldec"
	"github.com/zclconf/go-cty/cty"

	"github.com/blackstork-io/fabric/plugin"
)

const (
	defaultCodeLanguage = ""
)

func makeCodeContentProvider() *plugin.ContentProvider {
	return &plugin.ContentProvider{
		ContentFunc: genCodeContent,
		Args: hcldec.ObjectSpec{
			"value": &hcldec.AttrSpec{
				Name:     "value",
				Type:     cty.String,
				Required: true,
			},
			"language": &hcldec.AttrSpec{
				Name:     "language",
				Type:     cty.String,
				Required: false,
			},
		},
	}
}

func genCodeContent(ctx context.Context, params *plugin.ProvideContentParams) (*plugin.ContentResult, hcl.Diagnostics) {
	value := params.Args.GetAttr("value")
	if value.IsNull() {
		return nil, hcl.Diagnostics{{
			Severity: hcl.DiagError,
			Summary:  "Failed to parse arguments",
			Detail:   "value is required",
		}}
	}
	lang := params.Args.GetAttr("language")
	if lang.IsNull() {
		lang = cty.StringVal(defaultCodeLanguage)
	}
	text, err := genTextContentText(value.AsString(), params.DataContext)
	if err != nil {
		return nil, hcl.Diagnostics{{
			Severity: hcl.DiagError,
			Summary:  "Failed to render code",
			Detail:   err.Error(),
		}}
	}
	text = fmt.Sprintf("```%s\n%s\n```", lang.AsString(), text)
	return &plugin.ContentResult{
		Content: &plugin.ContentElement{
			Markdown: text,
		},
	}, nil
}
