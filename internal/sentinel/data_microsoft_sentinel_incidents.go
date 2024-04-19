package sentinel

import (
	"context"
	"fmt"

	"github.com/hashicorp/hcl/v2"
	"github.com/zclconf/go-cty/cty"

	"github.com/blackstork-io/fabric/internal/sentinel/client"
	"github.com/blackstork-io/fabric/plugin"
	"github.com/blackstork-io/fabric/plugin/dataspec"
)

func makeMicrosoftSentinelIncidentsDataSource(loader ClientLoadFn) *plugin.DataSource {
	return &plugin.DataSource{
		Doc:      "The `microsoft_sentinel_incidents` data source fetches incidents from Microsoft Sentinel.",
		DataFunc: fetchMicrosoftSentinelIncidents(loader),
		Config: dataspec.ObjectSpec{
			&dataspec.AttrSpec{
				Doc:      "The Azure subscription ID",
				Name:     "subscription_id",
				Type:     cty.String,
				Required: true,
			},
			&dataspec.AttrSpec{
				Doc:      "The Azure resource group name",
				Name:     "resource_group_name",
				Type:     cty.String,
				Required: true,
			},
			&dataspec.AttrSpec{
				Doc:      "The Azure workspace name",
				Name:     "workspace_name",
				Type:     cty.String,
				Required: true,
			},
		},
		Args: dataspec.ObjectSpec{
			&dataspec.AttrSpec{
				Doc:      "The filter expression",
				Name:     "filter",
				Type:     cty.String,
				Required: false,
			},
			&dataspec.AttrSpec{
				Doc:      "The maximum number of incidents to return",
				Name:     "limit",
				Type:     cty.Number,
				Required: false,
			},
			&dataspec.AttrSpec{
				Doc:      "The order by expression",
				Name:     "order_by",
				Type:     cty.String,
				Required: false,
			},
		},
	}
}

func fetchMicrosoftSentinelIncidents(loader ClientLoadFn) plugin.RetrieveDataFunc {
	return func(ctx context.Context, params *plugin.RetrieveDataParams) (plugin.Data, hcl.Diagnostics) {
		client, err := makeClient(loader, params.Config)
		if err != nil {
			return nil, hcl.Diagnostics{{
				Severity: hcl.DiagError,
				Summary:  "Unable to create Microsoft Sentinel client",
				Detail:   err.Error(),
			}}
		}
		req, err := parseMicrosoftSentinelIncidentsArgs(params.Config, params.Args)
		if err != nil {
			return nil, hcl.Diagnostics{{
				Severity: hcl.DiagError,
				Summary:  "Failed to parse arguments",
				Detail:   err.Error(),
			}}
		}
		res, err := client.ListIncidents(ctx, req)
		if err != nil {
			return nil, hcl.Diagnostics{{
				Severity: hcl.DiagError,
				Summary:  "Unable to list Microsoft Sentinel incidents",
				Detail:   err.Error(),
			}}
		}

		var data plugin.ListData
		for _, incident := range res.Value {
			item, err := plugin.ParseDataAny(incident)
			if err != nil {
				return nil, hcl.Diagnostics{{
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

func parseMicrosoftSentinelIncidentsArgs(cfg, args cty.Value) (*client.ListIncidentsReq, error) {
	var req client.ListIncidentsReq
	if param := cfg.GetAttr("subscription_id"); !param.IsNull() {
		req.SubscriptionID = param.AsString()
	} else {
		return nil, fmt.Errorf("subscription_id is required")
	}
	if param := cfg.GetAttr("resource_group_name"); !param.IsNull() {
		req.ResourceGroupName = param.AsString()
	} else {
		return nil, fmt.Errorf("resource_group_name is required")
	}
	if param := cfg.GetAttr("workspace_name"); !param.IsNull() {
		req.WorkspaceName = param.AsString()
	} else {
		return nil, fmt.Errorf("workspace_name is required")
	}
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
