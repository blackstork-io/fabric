package elastic

import (
	"context"
	"log/slog"

	"github.com/hashicorp/hcl/v2"
	"github.com/zclconf/go-cty/cty"

	"github.com/blackstork-io/fabric/pkg/diagnostics"
	"github.com/blackstork-io/fabric/plugin"
	"github.com/blackstork-io/fabric/plugin/dataspec"
	"github.com/blackstork-io/fabric/plugin/dataspec/constraint"
)

const (
	defaultSearchResultsSize   = 1000
	maxSimpleSearchResultsSize = 10000
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
				Name:   "api_key_str",
				Type:   cty.String,
				Secret: true,
			},
			&dataspec.AttrSpec{
				Name:   "api_key",
				Type:   cty.List(cty.String),
				Secret: true,
			},
			&dataspec.AttrSpec{
				Name: "basic_auth_username",
				Type: cty.String,
			},
			&dataspec.AttrSpec{
				Name:   "basic_auth_password",
				Type:   cty.String,
				Secret: true,
			},
			&dataspec.AttrSpec{
				Name:   "bearer_auth",
				Type:   cty.String,
				Secret: true,
			},
			&dataspec.AttrSpec{
				Name:   "ca_certs",
				Type:   cty.String,
				Secret: true,
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
				Type: cty.DynamicPseudoType,
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

func fetchElasticSearchData(ctx context.Context, params *plugin.RetrieveDataParams) (plugin.Data, diagnostics.Diag) {
	client, err := makeSearchClient(params.Config)
	if err != nil {
		return nil, diagnostics.Diag{{
			Severity: hcl.DiagError,
			Summary:  "Failed to create elasticsearch client",
			Detail:   err.Error(),
		}}
	}
	var diags diagnostics.Diag
	if (params.Args.GetAttr("only_hits").IsNull() || params.Args.GetAttr("only_hits").True()) &&
		!params.Args.GetAttr("aggs").IsNull() {
		if params.Args.GetAttr("query").IsNull() && params.Args.GetAttr("query_string").IsNull() {
			return nil, diagnostics.Diag{{
				Severity: hcl.DiagError,
				Summary:  "Invalid arguments",
				Detail:   "Aggregations are not supported without a query or query_string",
			}}
		}
		diags.Append(&hcl.Diagnostic{
			Severity: hcl.DiagWarning,
			Summary:  "Aggregations are not supported",
			Detail:   "Aggregations are not supported when only_hits is true",
		})
	}

	index := params.Args.GetAttr("index")
	if index.IsNull() {
		return nil, diagnostics.Diag{{
			Severity: hcl.DiagError,
			Summary:  "Invalid arguments",
			Detail:   "Index value is required",
		}}
	}

	id := params.Args.GetAttr("id")
	var data plugin.Data
	if !id.IsNull() {
		data, err = getByID(client.Get, params.Args)
	} else {
		var size int = defaultSearchResultsSize
		if sizeArg := params.Args.GetAttr("size"); !sizeArg.IsNull() {
			size64, _ := sizeArg.AsBigFloat().Int64()
			size = int(size64)
			if size < 0 {
				return nil, diagnostics.Diag{{
					Severity: hcl.DiagError,
					Summary:  "Invalid arguments",
					Detail:   "Size value must be greater or equal 0",
				}}
			}
		}
		if size <= maxSimpleSearchResultsSize {
			slog.DebugContext(ctx, "Sending normal search request", "size", size)
			data, err = search(client.Search, params.Args, size)
		} else {
			slog.DebugContext(ctx,"Starting a scroll search request", "size", size)
			data, err = searchWithScroll(client, params.Args, size)
		}
	}
	if err != nil {
		diags.Append(&hcl.Diagnostic{
			Severity: hcl.DiagError,
			Summary:  "Failed to fetch data",
			Detail:   err.Error(),
		})
	} else {
		slog.DebugContext(ctx, "Returning data received from Elasticsearch")
	}
	return data, diags
}
