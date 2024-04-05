package snyk

import (
	"context"
	"fmt"

	"github.com/blackstork-io/fabric/internal/snyk/client"
	"github.com/blackstork-io/fabric/plugin"
	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/hcldec"
	"github.com/zclconf/go-cty/cty"
)

func makeSnykIssuesDataSource(loader ClientLoadFn) *plugin.DataSource {
	return &plugin.DataSource{
		DataFunc: fetchSnykIssues(loader),
		Config: hcldec.ObjectSpec{
			"api_key": &hcldec.AttrSpec{
				Name:     "api_key",
				Type:     cty.String,
				Required: true,
			},
		},
		Args: hcldec.ObjectSpec{
			"group_id": &hcldec.AttrSpec{
				Name:     "group_id",
				Type:     cty.String,
				Required: false,
			},
			"org_id": &hcldec.AttrSpec{
				Name:     "org_id",
				Type:     cty.String,
				Required: false,
			},
			"scan_item_id": &hcldec.AttrSpec{
				Name:     "scan_item_id",
				Type:     cty.String,
				Required: false,
			},
			"scan_item_type": &hcldec.AttrSpec{
				Name:     "scan_item_type",
				Type:     cty.String,
				Required: false,
			},
			"type": &hcldec.AttrSpec{
				Name:     "type",
				Type:     cty.String,
				Required: false,
			},
			"updated_before": &hcldec.AttrSpec{
				Name:     "updated_before",
				Type:     cty.String,
				Required: false,
			},
			"updated_after": &hcldec.AttrSpec{
				Name:     "updated_after",
				Type:     cty.String,
				Required: false,
			},
			"created_before": &hcldec.AttrSpec{
				Name:     "created_before",
				Type:     cty.String,
				Required: false,
			},
			"created_after": &hcldec.AttrSpec{
				Name:     "created_after",
				Type:     cty.String,
				Required: false,
			},
			"effective_severity_level": &hcldec.AttrSpec{
				Name:     "effective_severity_level",
				Type:     cty.List(cty.String),
				Required: false,
			},
			"status": &hcldec.AttrSpec{
				Name:     "status",
				Type:     cty.List(cty.String),
				Required: false,
			},
			"ignored": &hcldec.AttrSpec{
				Name:     "ignored",
				Type:     cty.Bool,
				Required: false,
			},
			"limit": &hcldec.AttrSpec{
				Name:     "limit",
				Type:     cty.Number,
				Required: false,
			},
		},
	}
}

func fetchSnykIssues(loader ClientLoadFn) plugin.RetrieveDataFunc {
	return func(ctx context.Context, params *plugin.RetrieveDataParams) (plugin.Data, hcl.Diagnostics) {
		client, err := makeClient(loader, params.Config)
		if err != nil {
			return nil, hcl.Diagnostics{
				{
					Severity: hcl.DiagError,
					Summary:  "Failed to create Snyk client",
					Detail:   err.Error(),
				},
			}
		}
		limit := 0
		if arg := params.Args.GetAttr("limit"); !arg.IsNull() {
			n, _ := arg.AsBigFloat().Int64()
			if n <= 0 {
				return nil, hcl.Diagnostics{
					{
						Severity: hcl.DiagError,
						Summary:  "Invalid limit",
						Detail:   "limit must be greater than 0",
					},
				}
			}
			limit = int(n)
		}
		req, err := makeRequest(params.Args)
		if err != nil {
			return nil, hcl.Diagnostics{
				{
					Severity: hcl.DiagError,
					Summary:  "Failed to create Snyk request",
					Detail:   err.Error(),
				},
			}
		}
		var data plugin.ListData
		for {
			res, err := client.ListIssues(ctx, req)
			if err != nil {
				return nil, hcl.Diagnostics{
					{
						Severity: hcl.DiagError,
						Summary:  "Failed to list Snyk issues",
						Detail:   err.Error(),
					},
				}
			}
			for _, v := range res.Data {
				item, err := plugin.ParseDataAny(v)
				if err != nil {
					return nil, hcl.Diagnostics{
						{
							Severity: hcl.DiagError,
							Summary:  "Failed to parse Snyk issue",
							Detail:   err.Error(),
						},
					}
				}
				data = append(data, item)
			}
			if limit > 0 && len(data) >= limit {
				data = data[:limit]
				break
			}
			if links := res.Links; links != nil && links.Next != nil {
				req.StartingAfter = links.Next
			} else {
				break
			}
		}
		return data, nil
	}
}

func makeRequest(args cty.Value) (*client.ListIssuesReq, error) {
	req := &client.ListIssuesReq{
		Limit: pageSize,
	}
	groupID := args.GetAttr("group_id")
	orgID := args.GetAttr("org_id")
	if groupID.IsNull() && orgID.IsNull() {
		return nil, fmt.Errorf("either group_id or org_id must be set")
	}
	if !groupID.IsNull() && !orgID.IsNull() {
		return nil, fmt.Errorf("only one of group_id or org_id is allowed")
	}
	if !groupID.IsNull() {
		req.GroupID = client.String(groupID.AsString())
	}
	if !orgID.IsNull() {
		req.OrgID = client.String(orgID.AsString())
	}
	if arg := args.GetAttr("scan_item_id"); !arg.IsNull() {
		req.ScanItemID = client.String(arg.AsString())
	}
	if arg := args.GetAttr("scan_item_type"); !arg.IsNull() {
		req.ScanItemType = client.String(arg.AsString())
	}
	if arg := args.GetAttr("type"); !arg.IsNull() {
		req.Type = client.String(arg.AsString())
	}
	if arg := args.GetAttr("updated_before"); !arg.IsNull() {
		req.UpdatedBefore = client.String(arg.AsString())
	}
	if arg := args.GetAttr("updated_after"); !arg.IsNull() {
		req.UpdatedAfter = client.String(arg.AsString())
	}
	if arg := args.GetAttr("created_before"); !arg.IsNull() {
		req.CreatedBefore = client.String(arg.AsString())
	}
	if arg := args.GetAttr("created_after"); !arg.IsNull() {
		req.CreatedAfter = client.String(arg.AsString())
	}
	if arg := args.GetAttr("effective_severity_level"); !arg.IsNull() {
		list := []string{}
		for _, v := range arg.AsValueSlice() {
			list = append(list, v.AsString())
		}
		req.EffectiveSeverityLevel = client.StringList(list)
	}
	if arg := args.GetAttr("status"); !arg.IsNull() {
		list := []string{}
		for _, v := range arg.AsValueSlice() {
			list = append(list, v.AsString())
		}
		req.Status = client.StringList(list)
	}
	if arg := args.GetAttr("ignored"); !arg.IsNull() {
		req.Ignored = client.Bool(arg.True())
	}

	return req, nil
}
