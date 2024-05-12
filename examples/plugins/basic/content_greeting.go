package basic

import (
	"context"
	"fmt"

	"github.com/hashicorp/hcl/v2"
	"github.com/zclconf/go-cty/cty"

	"github.com/blackstork-io/fabric/pkg/diagnostics"
	"github.com/blackstork-io/fabric/plugin"
	"github.com/blackstork-io/fabric/plugin/dataspec"
	"github.com/blackstork-io/fabric/plugin/dataspec/constraint"
)

// makeGreetingContentProvider creates a new content provider that prints out a greeting message
func makeGreetingContentProvider() *plugin.ContentProvider {
	return &plugin.ContentProvider{
		// Config is optional, in this case we don't need it
		// We only define the schema for the arguments
		Args: dataspec.ObjectSpec{
			&dataspec.AttrSpec{
				Name:        "name",
				Constraints: constraint.RequiredMeaningfull,
				Doc:         `Name of the user`,
				ExampleVal:  cty.StringVal("John"),
				Type:        cty.String,
			},
		},
		// Optional: We can also define the schema for the config
		ContentFunc: renderGreetingMessage,
	}
}

func renderGreetingMessage(ctx context.Context, params *plugin.ProvideContentParams) (*plugin.ContentResult, diagnostics.Diag) {
	name := params.Config.GetAttr("name")
	if name.IsNull() || name.AsString() == "" {
		return nil, diagnostics.Diag{{
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
