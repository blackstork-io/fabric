package crowdstrike

import (
	"context"

	"github.com/crowdstrike/gofalcon/falcon"
	"github.com/crowdstrike/gofalcon/falcon/client/discover"
	"github.com/hashicorp/hcl/v2"
	"github.com/zclconf/go-cty/cty"

	"github.com/blackstork-io/fabric/pkg/diagnostics"
	"github.com/blackstork-io/fabric/plugin"
	"github.com/blackstork-io/fabric/plugin/dataspec"
	"github.com/blackstork-io/fabric/plugin/dataspec/constraint"
	"github.com/blackstork-io/fabric/plugin/plugindata"
)

func makeFalconDiscoverHostDetailsDataSource(loader ClientLoaderFn) *plugin.DataSource {
	return &plugin.DataSource{
		Doc:      "The `falcon_discover_host_details` data source fetches host details from Falcon Discover Host API.",
		DataFunc: fetchFalconDiscoverHostDetails(loader),
		Config:   makeDataSourceConfig(),
		Args: &dataspec.RootSpec{
			Attrs: []*dataspec.AttrSpec{
				{
					Name:        "size",
					Type:        cty.Number,
					Constraints: constraint.Integer | constraint.RequiredNonNull,
					Doc:         "limit the number of queried items",
				},
				{
					Name: "filter",
					Type: cty.String,
					Doc:  "Host search expression using Falcon Query Language (FQL)",
				},
			},
		},
	}
}

func fetchFalconDiscoverHostDetails(loader ClientLoaderFn) plugin.RetrieveDataFunc {
	return func(ctx context.Context, params *plugin.RetrieveDataParams) (plugindata.Data, diagnostics.Diag) {
		cli, err := loader(makeApiConfig(ctx, params.Config))
		if err != nil {
			return nil, diagnostics.Diag{{
				Severity: hcl.DiagError,
				Summary:  "Unable to create falcon client",
				Detail:   err.Error(),
			}}
		}
		size, _ := params.Args.GetAttrVal("size").AsBigFloat().Int64()
		queryHostParams := discover.NewQueryHostsParams().WithDefaults()
		queryHostParams.SetLimit(&size)
		queryHostParams.SetContext(ctx)
		queryHostsResponse, err := cli.Discover().QueryHosts(queryHostParams)
		if err != nil {
			return nil, diagnostics.Diag{{
				Severity: hcl.DiagError,
				Summary:  "Failed to query Falcon Discover Hosts",
				Detail:   err.Error(),
			}}
		}
		if err = falcon.AssertNoError(queryHostsResponse.GetPayload().Errors); err != nil {
			return nil, diagnostics.Diag{{
				Severity: hcl.DiagError,
				Summary:  "Failed to query Falcon Discover Hosts",
				Detail:   err.Error(),
			}}
		}
		hostIds := queryHostsResponse.GetPayload().Resources

		getHostParams := discover.NewGetHostsParams().WithDefaults()
		getHostParams.SetIds(hostIds)
		getHostParams.SetContext(ctx)
		getHostsResponse, err := cli.Discover().GetHosts(getHostParams)
		if err != nil {
			return nil, diagnostics.Diag{{
				Severity: hcl.DiagError,
				Summary:  "Failed to fetch Falcon Discover Hosts",
				Detail:   err.Error(),
			}}
		}
		if err = falcon.AssertNoError(queryHostsResponse.GetPayload().Errors); err != nil {
			return nil, diagnostics.Diag{{
				Severity: hcl.DiagError,
				Summary:  "Failed to fetch Falcon Discover Hosts",
				Detail:   err.Error(),
			}}
		}

		resources := getHostsResponse.GetPayload().Resources
		data, err := encodeResponse(resources)
		if err != nil {
			return nil, diagnostics.Diag{{
				Severity: hcl.DiagError,
				Summary:  "Failed to parse response",
				Detail:   err.Error(),
			}}
		}
		return data, nil
	}
}
