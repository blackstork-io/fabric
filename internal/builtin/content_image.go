package builtin

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/hcl/v2"
	"github.com/zclconf/go-cty/cty"

	"github.com/blackstork-io/fabric/plugin"
	"github.com/blackstork-io/fabric/plugin/dataspec"
)

func makeImageContentProvider() *plugin.ContentProvider {
	return &plugin.ContentProvider{
		ContentFunc: genImageContent,
		Args: dataspec.ObjectSpec{
			&dataspec.AttrSpec{
				Name:       "src",
				Type:       cty.String,
				Required:   true,
				ExampleVal: cty.StringVal("https://example.com/img.png"),
			},
			&dataspec.AttrSpec{
				Name:       "alt",
				Type:       cty.String,
				Required:   false,
				ExampleVal: cty.StringVal("Text description of the image"),
				// Not using empty string as DefaultVal here for semantical meaning
			},
		},
		Doc: "Inserts an image",
	}
}

func genImageContent(ctx context.Context, params *plugin.ProvideContentParams) (*plugin.ContentResult, hcl.Diagnostics) {
	src := params.Args.GetAttr("src")
	if src.IsNull() || src.AsString() == "" {
		return nil, hcl.Diagnostics{{
			Severity: hcl.DiagError,
			Summary:  "Failed to parse arguments",
			Detail:   "src is required",
		}}
	}
	alt := params.Args.GetAttr("alt")
	if alt.IsNull() {
		alt = cty.StringVal("")
	}
	srcStr := strings.TrimSpace(strings.ReplaceAll(src.AsString(), "\n", ""))
	altStr := strings.TrimSpace(strings.ReplaceAll(alt.AsString(), "\n", ""))
	return &plugin.ContentResult{
		Content: &plugin.ContentElement{
			Markdown: fmt.Sprintf("![%s](%s)", altStr, srcStr),
		},
	}, nil
}
