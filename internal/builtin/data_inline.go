package builtin

import (
	"context"

	"github.com/hashicorp/hcl/v2"

	"github.com/blackstork-io/fabric/pkg/diagnostics"
	"github.com/blackstork-io/fabric/plugin"
	"github.com/blackstork-io/fabric/plugin/dataspec"
)

func makeInlineDataSource() *plugin.DataSource {
	return &plugin.DataSource{
		DataFunc: fetchInlineData,
		Args: &dataspec.ObjDumpSpec{
			Doc: `
				Arbitrary structure of (possibly nested) blocks and attributes.
				For example:
				  key1 = "value1"
				  nested {
				    blocks {
				      key2 = 42
				    }
				  }
			`,
		},
		Doc: `Creates a queryable key-value map from the block's contents`,
	}
}

func fetchInlineData(ctx context.Context, params *plugin.RetrieveDataParams) (plugin.Data, diagnostics.Diag) {
	if params.Args.IsNull() {
		return nil, diagnostics.Diag{{
			Severity: hcl.DiagError,
			Summary:  "Failed to parse arguments",
			Detail:   "inline data is required",
		}}
	}
	if !params.Args.Type().IsMapType() && !params.Args.Type().IsObjectType() {
		return nil, diagnostics.Diag{{
			Severity: hcl.DiagError,
			Summary:  "Failed to parse arguments",
			Detail:   "inline data must be a map",
		}}
	}
	return plugin.ConvertCtyToData(params.Args), nil
}
