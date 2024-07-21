package microsoft

import (
	"context"
	"fmt"

	"github.com/hashicorp/hcl/v2"
	"github.com/zclconf/go-cty/cty"

	"github.com/blackstork-io/fabric/internal/microsoft/client"
	"github.com/blackstork-io/fabric/pkg/diagnostics"
	"github.com/blackstork-io/fabric/plugin"
	"github.com/blackstork-io/fabric/plugin/dataspec"
	"github.com/blackstork-io/fabric/plugin/dataspec/constraint"
)

func makeMicrosoftSentinelIncidentsDataSource(loader ClientLoadFn) *plugin.DataSource {
	return &plugin.DataSource{
		Doc:      "The `microsoft_sentinel_incidents` data source fetches incidents from Microsoft Sentinel.",
		DataFunc: fetchMicrosoftSentinelIncidents(loader),
		Config: &dataspec.RootSpec{
			Attrs: []*dataspec.AttrSpec{
				{
					Doc:         "The Azure client ID",
					Name:        "client_id",
					Type:        cty.String,
					Constraints: constraint.RequiredNonNull,
				},
				{
					Doc:         "The Azure client secret",
					Name:        "client_secret",
					Type:        cty.String,
					Constraints: constraint.RequiredNonNull,
					Secret:      true,
				},
				{
					Doc:         "The Azure tenant ID",
					Name:        "tenant_id",
					Type:        cty.String,
					Constraints: constraint.RequiredNonNull,
				},
				{
					Doc:         "The Azure subscription ID",
					Name:        "subscription_id",
					Type:        cty.String,
					Constraints: constraint.RequiredNonNull,
				},
				{
					Doc:         "The Azure resource group name",
					Name:        "resource_group_name",
					Type:        cty.String,
					Constraints: constraint.RequiredNonNull,
				},
				{
					Doc:         "The Azure workspace name",
					Name:        "workspace_name",
					Type:        cty.String,
					Constraints: constraint.RequiredNonNull,
				},
			},
		},
		Args: &dataspec.RootSpec{
			Attrs: []*dataspec.AttrSpec{
				{
					Doc:  "The filter expression",
					Name: "filter",
					Type: cty.String,
				},
				{
					Doc:  "The maximum number of incidents to return",
					Name: "limit",
					Type: cty.Number,
				},
				{
					Doc:  "The order by expression",
					Name: "order_by",
					Type: cty.String,
				},
			},
		},
	}
}

func fetchMicrosoftSentinelIncidents(loader ClientLoadFn) plugin.RetrieveDataFunc {
	return func(ctx context.Context, params *plugin.RetrieveDataParams) (plugin.Data, diagnostics.Diag) {
		client, err := makeClient(ctx, loader, params.Config)
		if err != nil {
			return nil, diagnostics.Diag{{
				Severity: hcl.DiagError,
				Summary:  "Unable to create Microsoft Sentinel client",
				Detail:   err.Error(),
			}}
		}
		req, err := parseMicrosoftSentinelIncidentsArgs(params.Config, params.Args)
		if err != nil {
			return nil, diagnostics.Diag{{
				Severity: hcl.DiagError,
				Summary:  "Failed to parse arguments",
				Detail:   err.Error(),
			}}
		}
		res, err := client.ListIncidents(ctx, req)
		if err != nil {
			return nil, diagnostics.Diag{{
				Severity: hcl.DiagError,
				Summary:  "Unable to list Microsoft Sentinel incidents",
				Detail:   err.Error(),
			}}
		}

		var data plugin.ListData
		for _, incident := range res.Value {
			item, err := plugin.ParseDataAny(incident)
			if err != nil {
				return nil, diagnostics.Diag{{
					Severity: hcl.DiagError,
					Summary:  "Unable to parse Microsoft Sentinel incident",
					Detail:   err.Error(),
				}}
			}
			data = append(data, item)
		}
		return data, nil
	}
}

func parseMicrosoftSentinelIncidentsArgs(cfg, args *dataspec.Block) (*client.ListIncidentsReq, error) {
	var req client.ListIncidentsReq

	req.SubscriptionID = cfg.GetAttr("subscription_id").AsString()
	req.ResourceGroupName = cfg.GetAttr("resource_group_name").AsString()
	req.WorkspaceName = cfg.GetAttr("workspace_name").AsString()

	if param := args.GetAttr("filter"); !param.IsNull() {
		req.Filter = client.String(param.AsString())
	}
	if param := args.GetAttr("limit"); !param.IsNull() {
		n, _ := param.AsBigFloat().Int64()
		if n <= 0 {
			return nil, fmt.Errorf("limit must be a positive number")
		}
		req.Top = client.Int(int(n))
	}
	if param := args.GetAttr("order_by"); !param.IsNull() {
		req.OrderBy = client.String(param.AsString())
	}
	return &req, nil
}
