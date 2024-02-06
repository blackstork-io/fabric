package plugin

import (
	"context"

	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/hcldec"
	"github.com/zclconf/go-cty/cty"
)

type DataSources map[string]*DataSource

func (ds DataSources) Validate() hcl.Diagnostics {
	if ds == nil {
		return nil
	}
	var diags hcl.Diagnostics
	for name, source := range ds {
		if source == nil {
			diags = append(diags, &hcl.Diagnostic{
				Severity: hcl.DiagError,
				Summary:  "Incomplete DataSourceSchema",
				Detail:   "DataSource '" + name + "' not loaded",
			})
		} else {
			diags = append(diags, source.Validate()...)
		}
	}
	return diags
}

type DataSource struct {
	DataFunc RetrieveDataFunc
	Args     hcldec.Spec
	Config   hcldec.Spec
}

func (ds *DataSource) Validate() hcl.Diagnostics {
	var diags hcl.Diagnostics
	if ds.DataFunc == nil {
		diags = append(diags, &hcl.Diagnostic{
			Severity: hcl.DiagError,
			Summary:  "Incomplete DataSourceSchema",
			Detail:   "DataSource function not loaded",
		})
	}
	if ds.Args == nil {
		diags = append(diags, &hcl.Diagnostic{
			Severity: hcl.DiagError,
			Summary:  "Incomplete DataSourceSchema",
			Detail:   "Missing args schema",
		})
	}
	return diags
}

func (ds *DataSource) Execute(ctx context.Context, params *RetrieveDataParams) (Data, hcl.Diagnostics) {
	if ds == nil {
		return nil, hcl.Diagnostics{{
			Severity: hcl.DiagError,
			Summary:  "Missing DataSource schema",
		}}
	}
	if ds.DataFunc == nil {
		return nil, hcl.Diagnostics{{
			Severity: hcl.DiagError,
			Summary:  "Incomplete DataSourceSchema",
			Detail:   "datasource function not loaded",
		}}
	}
	return ds.DataFunc(ctx, params)
}

type RetrieveDataParams struct {
	Config cty.Value
	Args   cty.Value
}

type RetrieveDataFunc func(ctx context.Context, params *RetrieveDataParams) (Data, hcl.Diagnostics)
