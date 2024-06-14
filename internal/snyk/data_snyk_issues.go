package snyk

import (
	"context"
	"fmt"

	"github.com/hashicorp/hcl/v2"
	"github.com/zclconf/go-cty/cty"

	"github.com/blackstork-io/fabric/internal/snyk/client"
	"github.com/blackstork-io/fabric/pkg/diagnostics"
	"github.com/blackstork-io/fabric/plugin"
	"github.com/blackstork-io/fabric/plugin/dataspec"
	"github.com/blackstork-io/fabric/plugin/dataspec/constraint"
)

func makeSnykIssuesDataSource(loader ClientLoadFn) *plugin.DataSource {
	return &plugin.DataSource{
		Doc:      "The `snyk_issues` data source fetches issues from Snyk.",
		DataFunc: fetchSnykIssues(loader),
		Config: dataspec.ObjectSpec{
			&dataspec.AttrSpec{
				Doc:         "The Snyk API key",
				Name:        "api_key",
				Type:        cty.String,
				Constraints: constraint.RequiredMeaningful,
				Secret:      true,
			},
		},
		Args: dataspec.ObjectSpec{
			&dataspec.AttrSpec{
				Doc:  "The group ID",
				Name: "group_id",
				Type: cty.String,
			},
			&dataspec.AttrSpec{
				Doc:  "The organization ID",
				Name: "org_id",
				Type: cty.String,
			},
			&dataspec.AttrSpec{
				Doc:  "The scan item ID",
				Name: "scan_item_id",
				Type: cty.String,
			},
			&dataspec.AttrSpec{
				Doc:  "The scan item type",
				Name: "scan_item_type",
				Type: cty.String,
			},
			&dataspec.AttrSpec{
				Doc:  "The issue type",
				Name: "type",
				Type: cty.String,
			},
			&dataspec.AttrSpec{
				Doc:  "The updated before date",
				Name: "updated_before",
				Type: cty.String,
			},
			&dataspec.AttrSpec{
				Doc:  "The updated after date",
				Name: "updated_after",
				Type: cty.String,
			},
			&dataspec.AttrSpec{
				Doc:  "The created before date",
				Name: "created_before",
				Type: cty.String,
			},
			&dataspec.AttrSpec{
				Doc:  "The created after date",
				Name: "created_after",
				Type: cty.String,
			},
			&dataspec.AttrSpec{
				Doc:  "The effective severity level",
				Name: "effective_severity_level",
				Type: cty.List(cty.String),
			},
			&dataspec.AttrSpec{
				Doc:  "The status",
				Name: "status",
				Type: cty.List(cty.String),
			},
			&dataspec.AttrSpec{
				Doc:  "The ignored flag",
				Name: "ignored",
				Type: cty.Bool,
			},
			&dataspec.AttrSpec{
				Doc:          "The limit of issues to fetch",
				Name:         "limit",
				Type:         cty.Number,
				MinInclusive: cty.NumberIntVal(1),
			},
		},
	}
}

func fetchSnykIssues(loader ClientLoadFn) plugin.RetrieveDataFunc {
	return func(ctx context.Context, params *plugin.RetrieveDataParams) (plugin.Data, diagnostics.Diag) {
		client, err := makeClient(loader, params.Config)
		if err != nil {
			return nil, diagnostics.Diag{
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
				return nil, diagnostics.Diag{
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
			return nil, diagnostics.Diag{
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
				return nil, diagnostics.Diag{
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
					return nil, diagnostics.Diag{
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
