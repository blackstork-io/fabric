package builtin

import (
	"context"
	"io"
	"os"

	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/hcldec"
	"github.com/zclconf/go-cty/cty"

	"github.com/blackstork-io/fabric/plugin"
)

func makeTXTDataSource() *plugin.DataSource {
	return &plugin.DataSource{
		DataFunc: fetchTXTData,
		Args: hcldec.ObjectSpec{
			"path": &hcldec.AttrSpec{
				Name:     "path",
				Type:     cty.String,
				Required: true,
			},
		},
	}
}

func fetchTXTData(ctx context.Context, params *plugin.RetrieveDataParams) (plugin.Data, hcl.Diagnostics) {
	path := params.Args.GetAttr("path")
	if path.IsNull() || path.AsString() == "" {
		return nil, hcl.Diagnostics{{
			Severity: hcl.DiagError,
			Summary:  "Failed to parse arguments",
			Detail:   "path is required",
		}}
	}
	f, err := os.Open(path.AsString())
	if err != nil {
		return nil, hcl.Diagnostics{{
			Severity: hcl.DiagError,
			Summary:  "Failed to open txt file",
			Detail:   err.Error(),
		}}
	}
	defer f.Close()
	data, err := io.ReadAll(f)
	if err != nil {
		return nil, hcl.Diagnostics{{
			Severity: hcl.DiagError,
			Summary:  "Failed to read txt file",
			Detail:   err.Error(),
		}}
	}
	return plugin.StringData(string(data)), nil
}
