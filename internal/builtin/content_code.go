package builtin

import (
	"context"
	"fmt"

	"github.com/hashicorp/hcl/v2"
	"github.com/zclconf/go-cty/cty"

	"github.com/blackstork-io/fabric/plugin"
	"github.com/blackstork-io/fabric/plugin/dataspec"
)

func makeCodeContentProvider() *plugin.ContentProvider {
	return &plugin.ContentProvider{
		ContentFunc: genCodeContent,
		Args: dataspec.ObjectSpec{
			&dataspec.AttrSpec{
				Name:       "value",
				Type:       cty.String,
				Required:   true,
				ExampleVal: cty.StringVal("Text to be formatted as a code block"),
			},
			&dataspec.AttrSpec{
				Name:       "language",
				Type:       cty.String,
				Required:   false,
				ExampleVal: cty.StringVal("python3"),
				DefaultVal: cty.StringVal(""),
				Doc:        `Specifiy the language for syntax highlighting`,
			},
		},
		Doc: "Formats text as code snippet",
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
