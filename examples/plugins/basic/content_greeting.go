package basic

import (
	"context"
	"fmt"

	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/hcldec"
	"github.com/zclconf/go-cty/cty"

	"github.com/blackstork-io/fabric/plugin"
)

// makeGreetingContentProvider creates a new content provider that prints out a greeting message
func makeGreetingContentProvider() *plugin.ContentProvider {
	return &plugin.ContentProvider{
		// Config is optional, in this case we don't need it
		// We only define the schema for the arguments
		Args: hcldec.ObjectSpec{
			"name": &hcldec.AttrSpec{
				Name:     "name",
				Required: true,
				Type:     cty.String,
			},
		},
		// Optional: We can also define the schema for the config
		ContentFunc: renderGreetingMessage,
	}
}

func renderGreetingMessage(ctx context.Context, params *plugin.ProvideContentParams) (*plugin.ContentResult, hcl.Diagnostics) {
	name := params.Config.GetAttr("name")
	if name.IsNull() || name.AsString() == "" {
		return nil, hcl.Diagnostics{{
			Severity: hcl.DiagError,
			Summary:  "Failed to parse arguments",
			Detail:   "name is required",
		}}
	}
	return &plugin.ContentResult{
		Content: &plugin.ContentElement{
			Markdown: fmt.Sprintf("Hello, %s!", name.AsString()),
		},
	}, nil
}
