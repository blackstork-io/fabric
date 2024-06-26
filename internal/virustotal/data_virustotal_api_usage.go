package virustotal

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/hcl/v2"
	"github.com/zclconf/go-cty/cty"

	"github.com/blackstork-io/fabric/internal/virustotal/client"
	"github.com/blackstork-io/fabric/pkg/diagnostics"
	"github.com/blackstork-io/fabric/plugin"
	"github.com/blackstork-io/fabric/plugin/dataspec"
	"github.com/blackstork-io/fabric/plugin/dataspec/constraint"
)

func makeVirusTotalAPIUsageDataSchema(loader ClientLoadFn) *plugin.DataSource {
	return &plugin.DataSource{
		DataFunc: fetchVirusTotalAPIUsageData(loader),
		Config: dataspec.ObjectSpec{
			&dataspec.AttrSpec{
				Name:        "api_key",
				Type:        cty.String,
				Constraints: constraint.RequiredNonNull,
				Secret:      true,
			},
		},
		Args: dataspec.ObjectSpec{
			&dataspec.AttrSpec{
				Name: "user_id",
				Type: cty.String,
			},
			&dataspec.AttrSpec{
				Name: "group_id",
				Type: cty.String,
			},
			&dataspec.AttrSpec{
				Name: "start_date",
				Type: cty.String,
			},
			&dataspec.AttrSpec{
				Name: "end_date",
				Type: cty.String,
			},
		},
	}
}

func fetchVirusTotalAPIUsageData(loader ClientLoadFn) plugin.RetrieveDataFunc {
	return func(ctx context.Context, params *plugin.RetrieveDataParams) (plugin.Data, diagnostics.Diag) {
		cli, err := makeClient(loader, params.Config)
		if err != nil {
			return nil, diagnostics.Diag{{
				Severity: hcl.DiagError,
				Summary:  "Failed to create client",
				Detail:   err.Error(),
			}}
		}
		args, err := parseAPIUsageArgs(params.Args)
		if err != nil {
			return nil, diagnostics.Diag{{
				Severity: hcl.DiagError,
				Summary:  "Failed to parse arguments",
				Detail:   err.Error(),
			}}
		}
		var data map[string]any
		if args.User != nil {
			req := &client.GetUserAPIUsageReq{
				User: *args.User,
			}
			if args.StartDate != nil {
				req.StartDate = &client.Date{Time: *args.StartDate}
			}
			if args.EndDate != nil {
				req.EndDate = &client.Date{Time: *args.EndDate}
			}

			res, err := cli.GetUserAPIUsage(ctx, req)
			if err != nil {
				return nil, diagnostics.Diag{{
					Severity: hcl.DiagError,
					Summary:  "Failed to fetch data",
					Detail:   err.Error(),
				}}
			}
			data = res.Data
		} else {
			req := &client.GetGroupAPIUsageReq{
				Group: *args.Group,
			}
			if args.StartDate != nil {
				req.StartDate = &client.Date{Time: *args.StartDate}
			}
			if args.EndDate != nil {
				req.EndDate = &client.Date{Time: *args.EndDate}
			}

			res, err := cli.GetGroupAPIUsage(ctx, req)
			if err != nil {
				return nil, diagnostics.Diag{{
					Severity: hcl.DiagError,
					Summary:  "Failed to fetch data",
					Detail:   err.Error(),
				}}
			}
			data = res.Data
		}
		result, err := plugin.ParseDataMapAny(data)
		if err != nil {
			return nil, diagnostics.Diag{{
				Severity: hcl.DiagError,
				Summary:  "Failed to parse data",
				Detail:   err.Error(),
			}}
		}
		return result, nil
	}
}

type apiUsageArgs struct {
	User      *string
	Group     *string
	StartDate *time.Time
	EndDate   *time.Time
}

func parseAPIUsageArgs(args cty.Value) (*apiUsageArgs, error) {
	dst := apiUsageArgs{}

	if args.IsNull() {
		return nil, fmt.Errorf("arguments are null")
	}

	if userID := args.GetAttr("user_id"); !userID.IsNull() {
		userIDStr := userID.AsString()
		dst.User = &userIDStr
	}
	if groupID := args.GetAttr("group_id"); !groupID.IsNull() {
		groupIDStr := groupID.AsString()
		dst.Group = &groupIDStr
	}
	if dst.User == nil && dst.Group == nil {
		return nil, fmt.Errorf("either user_id or group_id must be set")
	}
	if startDate := args.GetAttr("start_date"); !startDate.IsNull() {
		startDateStr := startDate.AsString()
		startDate, err := time.Parse("20060102", startDateStr)
		if err != nil {
			return nil, fmt.Errorf("failed to parse start_date: %w", err)
		}
		dst.StartDate = &startDate
	}
	if endDate := args.GetAttr("end_date"); !endDate.IsNull() {
		endDateStr := endDate.AsString()
		endDate, err := time.Parse("20060102", endDateStr)
		if err != nil {
			return nil, fmt.Errorf("failed to parse end_date: %w", err)
		}
		dst.EndDate = &endDate
	}
	return &dst, nil
}
