package hackerone

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/hcl/v2"
	"github.com/zclconf/go-cty/cty"

	"github.com/blackstork-io/fabric/internal/hackerone/client"
	"github.com/blackstork-io/fabric/plugin"
	"github.com/blackstork-io/fabric/plugin/dataspec"
	"github.com/blackstork-io/fabric/plugin/dataspec/constraint"
)

func makeHackerOneReportsDataSchema(loader ClientLoadFn) *plugin.DataSource {
	return &plugin.DataSource{
		Config: dataspec.ObjectSpec{
			&dataspec.AttrSpec{
				Name:        "api_username",
				Type:        cty.String,
				Constraints: constraint.RequiredNonNull,
			},
			&dataspec.AttrSpec{
				Name:        "api_token",
				Type:        cty.String,
				Constraints: constraint.RequiredNonNull,
			},
		},
		Args: dataspec.ObjectSpec{
			&dataspec.AttrSpec{
				Name: "size",
				Type: cty.Number,
			},
			&dataspec.AttrSpec{
				Name: "page_number",
				Type: cty.Number,
			},
			&dataspec.AttrSpec{
				Name: "sort",
				Type: cty.String,
			},
			&dataspec.AttrSpec{
				Name: "program",
				Type: cty.List(cty.String),
			},
			&dataspec.AttrSpec{
				Name: "inbox_ids",
				Type: cty.List(cty.Number),
			},
			&dataspec.AttrSpec{
				Name: "reporter",
				Type: cty.List(cty.String),
			},
			&dataspec.AttrSpec{
				Name: "assignee",
				Type: cty.List(cty.String),
			},
			&dataspec.AttrSpec{
				Name: "state",
				Type: cty.List(cty.String),
			},
			&dataspec.AttrSpec{
				Name: "id",
				Type: cty.List(cty.Number),
			},
			&dataspec.AttrSpec{
				Name: "weakness_id",
				Type: cty.List(cty.Number),
			},
			&dataspec.AttrSpec{
				Name: "severity",
				Type: cty.List(cty.String),
			},
			&dataspec.AttrSpec{
				Name: "hacker_published",
				Type: cty.Bool,
			},
			&dataspec.AttrSpec{
				Name: "created_at__gt",
				Type: cty.String,
			},
			&dataspec.AttrSpec{
				Name: "created_at__lt",
				Type: cty.String,
			},
			&dataspec.AttrSpec{
				Name: "submitted_at__gt",
				Type: cty.String,
			},
			&dataspec.AttrSpec{
				Name: "submitted_at__lt",
				Type: cty.String,
			},
			&dataspec.AttrSpec{
				Name: "triaged_at__gt",
				Type: cty.String,
			},
			&dataspec.AttrSpec{
				Name: "triaged_at__lt",
				Type: cty.String,
			},
			&dataspec.AttrSpec{
				Name: "triaged_at__null",
				Type: cty.Bool,
			},
			&dataspec.AttrSpec{
				Name: "closed_at__gt",
				Type: cty.String,
			},
			&dataspec.AttrSpec{
				Name: "closed_at__lt",
				Type: cty.String,
			},
			&dataspec.AttrSpec{
				Name: "closed_at__null",
				Type: cty.Bool,
			},
			&dataspec.AttrSpec{
				Name: "disclosed_at__gt",
				Type: cty.String,
			},
			&dataspec.AttrSpec{
				Name: "disclosed_at__lt",
				Type: cty.String,
			},
			&dataspec.AttrSpec{
				Name: "disclosed_at__null",
				Type: cty.Bool,
			},
			&dataspec.AttrSpec{
				Name: "reporter_agreed_on_going_public",
				Type: cty.Bool,
			},
			&dataspec.AttrSpec{
				Name: "bounty_awarded_at__gt",
				Type: cty.String,
			},
			&dataspec.AttrSpec{
				Name: "bounty_awarded_at__lt",
				Type: cty.String,
			},
			&dataspec.AttrSpec{
				Name: "bounty_awarded_at__null",
				Type: cty.Bool,
			},
			&dataspec.AttrSpec{
				Name: "swag_awarded_at__gt",
				Type: cty.String,
			},
			&dataspec.AttrSpec{
				Name: "swag_awarded_at__lt",
				Type: cty.String,
			},
			&dataspec.AttrSpec{
				Name: "swag_awarded_at__null",
				Type: cty.Bool,
			},
			&dataspec.AttrSpec{
				Name: "last_report_activity_at__gt",
				Type: cty.String,
			},
			&dataspec.AttrSpec{
				Name: "last_report_activity_at__lt",
				Type: cty.String,
			},
			&dataspec.AttrSpec{
				Name: "first_program_activity_at__gt",
				Type: cty.String,
			},
			&dataspec.AttrSpec{
				Name: "first_program_activity_at__lt",
				Type: cty.String,
			},
			&dataspec.AttrSpec{
				Name: "first_program_activity_at__null",
				Type: cty.Bool,
			},
			&dataspec.AttrSpec{
				Name: "last_program_activity_at__gt",
				Type: cty.String,
			},
			&dataspec.AttrSpec{
				Name: "last_program_activity_at__lt",
				Type: cty.String,
			},
			&dataspec.AttrSpec{
				Name: "last_program_activity_at__null",
				Type: cty.Bool,
			},
			&dataspec.AttrSpec{
				Name: "last_activity_at__gt",
				Type: cty.String,
			},
			&dataspec.AttrSpec{
				Name: "last_activity_at__lt",
				Type: cty.String,
			},
			&dataspec.AttrSpec{
				Name: "last_public_activity_at__gt",
				Type: cty.String,
			},
			&dataspec.AttrSpec{
				Name: "last_public_activity_at__lt",
				Type: cty.String,
			},
			&dataspec.AttrSpec{
				Name: "keyword",
				Type: cty.String,
			},
			&dataspec.AttrSpec{
				Name: "custom_fields",
				Type: cty.Map(cty.String),
			},
		},
		DataFunc: fetchHackerOneReports(loader),
	}
}

func fetchHackerOneReports(loader ClientLoadFn) plugin.RetrieveDataFunc {
	return func(ctx context.Context, params *plugin.RetrieveDataParams) (plugin.Data, hcl.Diagnostics) {
		cli, err := makeClient(loader, params.Config)
		if err != nil {
			return nil, hcl.Diagnostics{{
				Severity: hcl.DiagError,
				Summary:  "Failed to create client",
				Detail:   err.Error(),
			}}
		}
		req, err := parseHackerOneReportsArgs(params.Args)
		if err != nil {
			return nil, hcl.Diagnostics{{
				Severity: hcl.DiagError,
				Summary:  "Failed to parse arguments",
				Detail:   err.Error(),
			}}
		}

		data := make([]any, 0)
		if req.PageNumber != nil {
			res, err := cli.GetAllReports(ctx, req)
			if err != nil {
				return nil, hcl.Diagnostics{{
					Severity: hcl.DiagError,
					Summary:  "Failed to fetch reports",
					Detail:   err.Error(),
				}}
			}
			data = append(data, res.Data...)
		} else {
			limit := -1
			if req.PageSize != nil {
				limit = *req.PageSize
			}
			for page := minPage; ; page++ {
				req.PageNumber = client.Int(page)
				res, err := cli.GetAllReports(ctx, req)
				if err != nil {
					return nil, hcl.Diagnostics{{
						Severity: hcl.DiagError,
						Summary:  "Failed to fetch reports",
						Detail:   err.Error(),
					}}
				}
				if res.Data == nil {
					res.Data = make([]any, 0)
				}
				data = append(data, res.Data...)
				if len(res.Data) == 0 || (limit > 0 && len(data) >= limit) {
					break
				}
			}
		}
		dst, err := plugin.ParseDataAny(data)
		if err != nil {
			return nil, hcl.Diagnostics{{
				Severity: hcl.DiagError,
				Summary:  "Failed to parse data",
				Detail:   err.Error(),
			}}
		}
		return dst, nil
	}
}

func parseHackerOneReportsArgs(args cty.Value) (*client.GetAllReportsReq, error) {
	if args.IsNull() {
		return nil, fmt.Errorf("args are required")
	}
	var req client.GetAllReportsReq
	size := args.GetAttr("size")
	if !size.IsNull() {
		n, _ := size.AsBigFloat().Int64()
		if n <= 0 {
			return nil, fmt.Errorf("size must be greater than 0")
		}
		req.PageSize = client.Int(int(n))
	}
	pageNumber := args.GetAttr("page_number")
	if !pageNumber.IsNull() {
		n, _ := pageNumber.AsBigFloat().Int64()
		if n <= 0 {
			return nil, fmt.Errorf("page_number must be greater than 0")
		}
		req.PageNumber = client.Int(int(n))
	}
	sort := args.GetAttr("sort")
	if !sort.IsNull() && sort.AsString() != "" {
		req.Sort = client.String(sort.AsString())
	}
	program := args.GetAttr("program")
	if !program.IsNull() {
		programs := program.AsValueSlice()
		for _, p := range programs {
			req.FilterProgram = append(req.FilterProgram, p.AsString())
		}
	}
	inboxIDs := args.GetAttr("inbox_ids")
	if !inboxIDs.IsNull() {
		ids := inboxIDs.AsValueSlice()
		for _, id := range ids {
			n, _ := id.AsBigFloat().Int64()
			req.FilterInboxIDs = append(req.FilterInboxIDs, int(n))
		}
	}
	if len(req.FilterProgram)+len(req.FilterInboxIDs) == 0 {
		return nil, fmt.Errorf("at least one of program or inbox_ids must be provided")
	}
	reporter := args.GetAttr("reporter")
	if !reporter.IsNull() {
		reporters := reporter.AsValueSlice()
		for _, r := range reporters {
			req.FilterReporter = append(req.FilterReporter, r.AsString())
		}
	}
	assignee := args.GetAttr("assignee")
	if !assignee.IsNull() {
		assignees := assignee.AsValueSlice()
		for _, a := range assignees {
			req.FilterAssignee = append(req.FilterAssignee, a.AsString())
		}
	}
	state := args.GetAttr("state")
	if !state.IsNull() {
		states := state.AsValueSlice()
		for _, s := range states {
			req.FilterState = append(req.FilterState, s.AsString())
		}
	}
	id := args.GetAttr("id")
	if !id.IsNull() {
		ids := id.AsValueSlice()
		for _, i := range ids {
			n, _ := i.AsBigFloat().Int64()
			req.FilterID = append(req.FilterID, int(n))
		}
	}
	weaknessID := args.GetAttr("weakness_id")
	if !weaknessID.IsNull() {
		ids := weaknessID.AsValueSlice()
		for _, i := range ids {
			n, _ := i.AsBigFloat().Int64()
			req.FilterWeaknessID = append(req.FilterWeaknessID, int(n))
		}
	}
	severity := args.GetAttr("severity")
	if !severity.IsNull() {
		severities := severity.AsValueSlice()
		for _, s := range severities {
			req.FilterSeverity = append(req.FilterSeverity, s.AsString())
		}
	}
	hackerPublished := args.GetAttr("hacker_published")
	if !hackerPublished.IsNull() {
		req.FilterHackerPublished = client.Bool(hackerPublished.True())
	}
	createdAtGT := args.GetAttr("created_at__gt")
	if !createdAtGT.IsNull() {
		t, err := time.Parse(time.RFC3339, createdAtGT.AsString())
		if err != nil {
			return nil, fmt.Errorf("failed to parse created_at__gt: %w", err)
		}
		req.FilterCreatedAtGT = &t
	}
	createdAtLT := args.GetAttr("created_at__lt")
	if !createdAtLT.IsNull() {
		t, err := time.Parse(time.RFC3339, createdAtLT.AsString())
		if err != nil {
			return nil, fmt.Errorf("failed to parse created_at__lt: %w", err)
		}
		req.FilterCreatedAtLT = &t
	}
	submittedAtGT := args.GetAttr("submitted_at__gt")
	if !submittedAtGT.IsNull() {
		t, err := time.Parse(time.RFC3339, submittedAtGT.AsString())
		if err != nil {
			return nil, fmt.Errorf("failed to parse submitted_at__gt: %w", err)
		}
		req.FilterSubmittedAtGT = &t
	}
	submittedAtLT := args.GetAttr("submitted_at__lt")
	if !submittedAtLT.IsNull() {
		t, err := time.Parse(time.RFC3339, submittedAtLT.AsString())
		if err != nil {
			return nil, fmt.Errorf("failed to parse submitted_at__lt: %w", err)
		}
		req.FilterSubmittedAtLT = &t
	}
	triagedAtGT := args.GetAttr("triaged_at__gt")
	if !triagedAtGT.IsNull() {
		t, err := time.Parse(time.RFC3339, triagedAtGT.AsString())
		if err != nil {
			return nil, fmt.Errorf("failed to parse triaged_at__gt: %w", err)
		}
		req.FilterTriagedAtGT = &t
	}
	triagedAtLT := args.GetAttr("triaged_at__lt")
	if !triagedAtLT.IsNull() {
		t, err := time.Parse(time.RFC3339, triagedAtLT.AsString())
		if err != nil {
			return nil, fmt.Errorf("failed to parse triaged_at__lt: %w", err)
		}
		req.FilterTriagedAtLT = &t
	}
	triagedAtNull := args.GetAttr("triaged_at__null")
	if !triagedAtNull.IsNull() {
		req.FilterTriagedAtNull = client.Bool(triagedAtNull.True())
	}
	closedAtGT := args.GetAttr("closed_at__gt")
	if !closedAtGT.IsNull() {
		t, err := time.Parse(time.RFC3339, closedAtGT.AsString())
		if err != nil {
			return nil, fmt.Errorf("failed to parse closed_at__gt: %w", err)
		}
		req.FilterClosedAtGT = &t
	}
	closedAtLT := args.GetAttr("closed_at__lt")
	if !closedAtLT.IsNull() {
		t, err := time.Parse(time.RFC3339, closedAtLT.AsString())
		if err != nil {
			return nil, fmt.Errorf("failed to parse closed_at__lt: %w", err)
		}
		req.FilterClosedAtLT = &t
	}
	closedAtNull := args.GetAttr("closed_at__null")
	if !closedAtNull.IsNull() {
		req.FilterClosedAtNull = client.Bool(closedAtNull.True())
	}
	disclosedAtGT := args.GetAttr("disclosed_at__gt")
	if !disclosedAtGT.IsNull() {
		t, err := time.Parse(time.RFC3339, disclosedAtGT.AsString())
		if err != nil {
			return nil, fmt.Errorf("failed to parse disclosed_at__gt: %w", err)
		}
		req.FilterDisclosedAtGT = &t
	}
	disclosedAtLT := args.GetAttr("disclosed_at__lt")
	if !disclosedAtLT.IsNull() {
		t, err := time.Parse(time.RFC3339, disclosedAtLT.AsString())
		if err != nil {
			return nil, fmt.Errorf("failed to parse disclosed_at__lt: %w", err)
		}
		req.FilterDisclosedAtLT = &t
	}
	disclosedAtNull := args.GetAttr("disclosed_at__null")
	if !disclosedAtNull.IsNull() {
		req.FilterDisclosedAtNull = client.Bool(disclosedAtNull.True())
	}
	reporterAgreedOnGoingPublic := args.GetAttr("reporter_agreed_on_going_public")
	if !reporterAgreedOnGoingPublic.IsNull() {
		req.FilterReporterAgreedOnGoingPublic = client.Bool(reporterAgreedOnGoingPublic.True())
	}
	bountyAwardedAtGT := args.GetAttr("bounty_awarded_at__gt")
	if !bountyAwardedAtGT.IsNull() {
		t, err := time.Parse(time.RFC3339, bountyAwardedAtGT.AsString())
		if err != nil {
			return nil, fmt.Errorf("failed to parse bounty_awarded_at__gt: %w", err)
		}
		req.FilterBountyAwardedAtGT = &t
	}
	bountyAwardedAtLT := args.GetAttr("bounty_awarded_at__lt")
	if !bountyAwardedAtLT.IsNull() {
		t, err := time.Parse(time.RFC3339, bountyAwardedAtLT.AsString())
		if err != nil {
			return nil, fmt.Errorf("failed to parse bounty_awarded_at__lt: %w", err)
		}
		req.FilterBountyAwardedAtLT = &t
	}
	bountyAwardedAtNull := args.GetAttr("bounty_awarded_at__null")
	if !bountyAwardedAtNull.IsNull() {
		req.FilterBountyAwardedAtNull = client.Bool(bountyAwardedAtNull.True())
	}
	swagAwardedAtGT := args.GetAttr("swag_awarded_at__gt")
	if !swagAwardedAtGT.IsNull() {
		t, err := time.Parse(time.RFC3339, swagAwardedAtGT.AsString())
		if err != nil {
			return nil, fmt.Errorf("failed to parse swag_awarded_at__gt: %w", err)
		}
		req.FilterSwagAwardedAtGT = &t
	}
	swagAwardedAtLT := args.GetAttr("swag_awarded_at__lt")
	if !swagAwardedAtLT.IsNull() {
		t, err := time.Parse(time.RFC3339, swagAwardedAtLT.AsString())
		if err != nil {
			return nil, fmt.Errorf("failed to parse swag_awarded_at__lt: %w", err)
		}
		req.FilterSwagAwardedAtLT = &t
	}
	swagAwardedAtNull := args.GetAttr("swag_awarded_at__null")
	if !swagAwardedAtNull.IsNull() {
		req.FilterSwagAwardedAtNull = client.Bool(swagAwardedAtNull.True())
	}
	lastReportActivityAtGT := args.GetAttr("last_report_activity_at__gt")
	if !lastReportActivityAtGT.IsNull() {
		t, err := time.Parse(time.RFC3339, lastReportActivityAtGT.AsString())
		if err != nil {
			return nil, fmt.Errorf("failed to parse last_report_activity_at__gt: %w", err)
		}
		req.FilterLastReportActivityAtGT = &t
	}
	lastReportActivityAtLT := args.GetAttr("last_report_activity_at__lt")
	if !lastReportActivityAtLT.IsNull() {
		t, err := time.Parse(time.RFC3339, lastReportActivityAtLT.AsString())
		if err != nil {
			return nil, fmt.Errorf("failed to parse last_report_activity_at__lt: %w", err)
		}
		req.FilterLastReportActivityAtLT = &t
	}
	firstProgramActivityAtGT := args.GetAttr("first_program_activity_at__gt")
	if !firstProgramActivityAtGT.IsNull() {
		t, err := time.Parse(time.RFC3339, firstProgramActivityAtGT.AsString())
		if err != nil {
			return nil, fmt.Errorf("failed to parse first_program_activity_at__gt: %w", err)
		}
		req.FilterFirstProgramActivityAtGT = &t
	}
	firstProgramActivityAtLT := args.GetAttr("first_program_activity_at__lt")
	if !firstProgramActivityAtLT.IsNull() {
		t, err := time.Parse(time.RFC3339, firstProgramActivityAtLT.AsString())
		if err != nil {
			return nil, fmt.Errorf("failed to parse first_program_activity_at__lt: %w", err)
		}
		req.FilterFirstProgramActivityAtLT = &t
	}
	firstProgramActivityAtNull := args.GetAttr("first_program_activity_at__null")
	if !firstProgramActivityAtNull.IsNull() {
		req.FilterFirstProgramActivityAtNull = client.Bool(firstProgramActivityAtNull.True())
	}
	lastProgramActivityAtGT := args.GetAttr("last_program_activity_at__gt")
	if !lastProgramActivityAtGT.IsNull() {
		t, err := time.Parse(time.RFC3339, lastProgramActivityAtGT.AsString())
		if err != nil {
			return nil, fmt.Errorf("failed to parse last_program_activity_at__gt: %w", err)
		}
		req.FilterLastProgramActivityAtGT = &t
	}
	lastProgramActivityAtLT := args.GetAttr("last_program_activity_at__lt")
	if !lastProgramActivityAtLT.IsNull() {
		t, err := time.Parse(time.RFC3339, lastProgramActivityAtLT.AsString())
		if err != nil {
			return nil, fmt.Errorf("failed to parse last_program_activity_at__lt: %w", err)
		}
		req.FilterLastProgramActivityAtLT = &t
	}
	lastProgramActivityAtNull := args.GetAttr("last_program_activity_at__null")
	if !lastProgramActivityAtNull.IsNull() {
		req.FilterLastProgramActivityAtNull = client.Bool(lastProgramActivityAtNull.True())
	}
	lastActivityAtGT := args.GetAttr("last_activity_at__gt")
	if !lastActivityAtGT.IsNull() {
		t, err := time.Parse(time.RFC3339, lastActivityAtGT.AsString())
		if err != nil {
			return nil, fmt.Errorf("failed to parse last_activity_at__gt: %w", err)
		}
		req.FilterLastActivityAtGT = &t
	}
	lastActivityAtLT := args.GetAttr("last_activity_at__lt")
	if !lastActivityAtLT.IsNull() {
		t, err := time.Parse(time.RFC3339, lastActivityAtLT.AsString())
		if err != nil {
			return nil, fmt.Errorf("failed to parse last_activity_at__lt: %w", err)
		}
		req.FilterLastActivityAtLT = &t
	}
	lastPublicActivityAtGT := args.GetAttr("last_public_activity_at__gt")
	if !lastPublicActivityAtGT.IsNull() {
		t, err := time.Parse(time.RFC3339, lastPublicActivityAtGT.AsString())
		if err != nil {
			return nil, fmt.Errorf("failed to parse last_public_activity_at__gt: %w", err)
		}
		req.FilterLastPublicActivityAtGT = &t
	}
	lastPublicActivityAtLT := args.GetAttr("last_public_activity_at__lt")
	if !lastPublicActivityAtLT.IsNull() {
		t, err := time.Parse(time.RFC3339, lastPublicActivityAtLT.AsString())
		if err != nil {
			return nil, fmt.Errorf("failed to parse last_public_activity_at__lt: %w", err)
		}
		req.FilterLastPublicActivityAtLT = &t
	}
	keyword := args.GetAttr("keyword")
	if !keyword.IsNull() {
		req.FilterKeyword = client.String(keyword.AsString())
	}
	customFields := args.GetAttr("custom_fields")
	if !customFields.IsNull() {
		fields := customFields.AsValueMap()
		req.FilterCustomFields = make(map[string]string, len(fields))
		for k, v := range fields {
			req.FilterCustomFields[k] = v.AsString()
		}
	}
	return &req, nil
}
