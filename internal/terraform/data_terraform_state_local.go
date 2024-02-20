package terraform

import (
	"context"
	"os"

	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/hcldec"
	"github.com/zclconf/go-cty/cty"

	"github.com/blackstork-io/fabric/plugin"
)

func makeTerraformStateLocalDataSource() *plugin.DataSource {
	return &plugin.DataSource{
		Config: nil,
		Args: hcldec.ObjectSpec{
			"path": &hcldec.AttrSpec{
				Name:     "path",
				Type:     cty.String,
				Required: true,
			},
		},
		DataFunc: fetchTerraformStateLocalData,
	}
}

func fetchTerraformStateLocalData(ctx context.Context, params *plugin.RetrieveDataParams) (plugin.Data, hcl.Diagnostics) {
	path := params.Args.GetAttr("path")
	if path.IsNull() || path.AsString() == "" {
		return nil, hcl.Diagnostics{{
			Severity: hcl.DiagError,
			Summary:  "Failed to parse arguments",
			Detail:   "path is required",
		}}
	}
	data, err := readTerraformStateFile(path.AsString())
	if err != nil {
		return nil, hcl.Diagnostics{{
			Severity: hcl.DiagError,
			Summary:  "Failed to read terraform state",
			Detail:   err.Error(),
		}}
	}
	return data, nil
}

func readTerraformStateFile(fp string) (plugin.Data, error) {
	data, err := os.ReadFile(fp)
	if err != nil {
		return nil, err
	}
	return plugin.UnmarshalJSONData(data)
}
