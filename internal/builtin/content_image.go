package builtin

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/hcldec"
	"github.com/zclconf/go-cty/cty"

	"github.com/blackstork-io/fabric/plugin"
)

func makeImageContentProvider() *plugin.ContentProvider {
	return &plugin.ContentProvider{
		ContentFunc: genImageContent,
		Args: hcldec.ObjectSpec{
			"src": &hcldec.AttrSpec{
				Name:     "src",
				Type:     cty.String,
				Required: true,
			},
			"alt": &hcldec.AttrSpec{
				Name:     "alt",
				Type:     cty.String,
				Required: false,
			},
		},
	}
}

func genImageContent(ctx context.Context, params *plugin.ProvideContentParams) (*plugin.Content, hcl.Diagnostics) {
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
	return &plugin.Content{
		Markdown: fmt.Sprintf("![%s](%s)", altStr, srcStr),
	}, nil
}
