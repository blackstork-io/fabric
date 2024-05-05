package elastic

import (
	"context"

	"github.com/hashicorp/hcl/v2"
	"github.com/zclconf/go-cty/cty"

	"github.com/blackstork-io/fabric/plugin"
	"github.com/blackstork-io/fabric/plugin/dataspec"
	"github.com/blackstork-io/fabric/plugin/dataspec/constraint"
)

func makeElasticSearchDataSource() *plugin.DataSource {
	return &plugin.DataSource{
		DataFunc: fetchElasticSearchData,
		Config: dataspec.ObjectSpec{
			&dataspec.AttrSpec{
				Name: "base_url",
				Type: cty.String,
			},
			&dataspec.AttrSpec{
				Name: "cloud_id",
				Type: cty.String,
			},
			&dataspec.AttrSpec{
				Name: "api_key_str",
				Type: cty.String,
			},
			&dataspec.AttrSpec{
				Name: "api_key",
				Type: cty.List(cty.String),
			},
			&dataspec.AttrSpec{
				Name: "basic_auth_username",
				Type: cty.String,
			},
			&dataspec.AttrSpec{
				Name: "basic_auth_password",
				Type: cty.String,
			},
			&dataspec.AttrSpec{
				Name: "bearer_auth",
				Type: cty.String,
			},
			&dataspec.AttrSpec{
				Name: "ca_certs",
				Type: cty.String,
			},
		},
		Args: dataspec.ObjectSpec{
			&dataspec.AttrSpec{
				Name:        "index",
				Type:        cty.String,
				Constraints: constraint.RequiredNonNull,
			},
			&dataspec.AttrSpec{
				Name: "id",
				Type: cty.String,
			},
			&dataspec.AttrSpec{
				Name: "query_string",
				Type: cty.String,
			},
			&dataspec.AttrSpec{
				Name: "query",
				Type: cty.Map(cty.DynamicPseudoType),
			},
			&dataspec.AttrSpec{
				Name: "aggs",
				Type: cty.Map(cty.DynamicPseudoType),
			},
			&dataspec.AttrSpec{
				Name: "only_hits",
				Type: cty.Bool,
			},
			&dataspec.AttrSpec{
				Name: "fields",
				Type: cty.List(cty.String),
			},
			&dataspec.AttrSpec{
				Name: "size",
				Type: cty.Number,
			},
		},
	}
}

func fetchElasticSearchData(ctx context.Context, params *plugin.RetrieveDataParams) (plugin.Data, hcl.Diagnostics) {
	client, err := makeSearchClient(params.Config)
	if err != nil {
		return nil, hcl.Diagnostics{{
			Severity: hcl.DiagError,
			Summary:  "Failed to create elasticsearch client",
			Detail:   err.Error(),
		}}
	}
	var diags hcl.Diagnostics
	if (params.Args.GetAttr("only_hits").IsNull() || params.Args.GetAttr("only_hits").True()) &&
		!params.Args.GetAttr("aggs").IsNull() {
		if params.Args.GetAttr("query").IsNull() && params.Args.GetAttr("query_string").IsNull() {
			return nil, hcl.Diagnostics{{
				Severity: hcl.DiagError,
				Summary:  "Invalid arguments",
				Detail:   "Aggregations are not supported without a query or query_string",
			}}
		}
		diags = diags.Append(&hcl.Diagnostic{
			Severity: hcl.DiagWarning,
			Summary:  "Aggregations are not supported",
			Detail:   "Aggregations are not supported when only_hits is true",
		})
	}

	id := params.Args.GetAttr("id")
	var data plugin.Data
	if !id.IsNull() {
		data, err = getByID(client.Get, params.Args)
	} else {
		data, err = search(client.Search, params.Args)
	}
	if err != nil {
		return nil, diags.Extend(hcl.Diagnostics{{
			Severity: hcl.DiagError,
			Summary:  "Failed to get data",
			Detail:   err.Error(),
		}})
	}
	return data, diags
}
