package github

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	gh "github.com/google/go-github/v58/github"
	"github.com/hashicorp/hcl/v2"
	"github.com/zclconf/go-cty/cty"

	"github.com/blackstork-io/fabric/pkg/diagnostics"
	"github.com/blackstork-io/fabric/plugin"
	"github.com/blackstork-io/fabric/plugin/dataspec"
	"github.com/blackstork-io/fabric/plugin/dataspec/constraint"
)

func makeGithubIssuesDataSchema(loader ClientLoaderFn) *plugin.DataSource {
	return &plugin.DataSource{
		DataFunc: fetchGithubIssuesData(loader),
		Config: dataspec.ObjectSpec{
			&dataspec.AttrSpec{
				Name:        "github_token",
				Type:        cty.String,
				Constraints: constraint.RequiredNonNull,
				Secret:      true,
			},
		},
		Args: dataspec.ObjectSpec{
			&dataspec.AttrSpec{
				Name:        "repository",
				Type:        cty.String,
				Constraints: constraint.RequiredNonNull,
			},
			&dataspec.AttrSpec{
				Name: "milestone",
				Type: cty.String,
			},
			&dataspec.AttrSpec{
				Name: "state",
				Type: cty.String,
			},
			&dataspec.AttrSpec{
				Name: "assignee",
				Type: cty.String,
			},
			&dataspec.AttrSpec{
				Name: "creator",
				Type: cty.String,
			},
			&dataspec.AttrSpec{
				Name: "mentioned",
				Type: cty.String,
			},
			&dataspec.AttrSpec{
				Name: "labels",
				Type: cty.List(cty.String),
			},
			&dataspec.AttrSpec{
				Name: "sort",
				Type: cty.String,
			},
			&dataspec.AttrSpec{
				Name: "direction",
				Type: cty.String,
			},
			&dataspec.AttrSpec{
				Name: "since",
				Type: cty.String,
			},
			&dataspec.AttrSpec{
				Name: "limit",
				Type: cty.Number,
			},
		},
	}
}

func fetchGithubIssuesData(loader ClientLoaderFn) plugin.RetrieveDataFunc {
	return func(ctx context.Context, params *plugin.RetrieveDataParams) (plugin.Data, diagnostics.Diag) {
		tkn, err := parseIssuesConfig(params.Config)
		if err != nil {
			return nil, diagnostics.Diag{{
				Severity: hcl.DiagError,
				Summary:  "Failed to parse configuration",
				Detail:   err.Error(),
			}}
		}
		opts, err := parseIssuesArgs(params.Args)
		if err != nil {
			return nil, diagnostics.Diag{{
				Severity: hcl.DiagError,
				Summary:  "Failed to parse arguments",
				Detail:   err.Error(),
			}}
		}
		client := loader(tkn)
		// iterate over pages until we get all issues or reach the limit if specified
		var issues plugin.ListData
		for page := minPage; ; page++ {
			opts.opts.Page = page
			opts.opts.PerPage = pageSize
			issuesPage, resp, err := client.Issues().ListByRepo(context.Background(), opts.owner, opts.name, opts.opts)
			if err != nil {
				return nil, diagnostics.Diag{{
					Severity: hcl.DiagError,
					Summary:  "Failed to list issues",
					Detail:   err.Error(),
				}}
			}
			for _, issue := range issuesPage {
				encoded, err := encodeGithubIssue(issue)
				if err != nil {
					return nil, diagnostics.Diag{{
						Severity: hcl.DiagError,
						Summary:  "Failed to encode issue",
						Detail:   err.Error(),
					}}
				}
				issues = append(issues, encoded)
			}
			if resp.NextPage == 0 || (opts.limit > 0 && int64(len(issues)) >= opts.limit) {
				break
			}
		}
		// if limit is specified, truncate the issues slice
		if opts.limit > 0 && int64(len(issues)) > opts.limit {
			issues = issues[:opts.limit]
		}
		return issues, nil
	}
}

func parseIssuesConfig(cfg cty.Value) (string, error) {
	githubToken := cfg.GetAttr("github_token")
	if githubToken.IsNull() || githubToken.AsString() == "" {
		return "", fmt.Errorf("github_token is required")
	}
	return githubToken.AsString(), nil
}

type parsedIssuesArgs struct {
	owner string
	name  string
	limit int64
	opts  *gh.IssueListByRepoOptions
}

func parseIssuesArgs(args cty.Value) (*parsedIssuesArgs, error) {
	repository := args.GetAttr("repository")
	if repository.IsNull() || repository.AsString() == "" {
		return nil, fmt.Errorf("repository is required")
	}
	repositoryParts := strings.Split(repository.AsString(), "/")
	if len(repositoryParts) != 2 {
		return nil, fmt.Errorf("repository must be in the format of owner/name")
	}
	owner := repositoryParts[0]
	name := repositoryParts[1]
	opts := &gh.IssueListByRepoOptions{}
	if milestone := args.GetAttr("milestone"); !milestone.IsNull() && milestone.AsString() != "" {
		opts.Milestone = milestone.AsString()
	}
	if state := args.GetAttr("state"); !state.IsNull() && state.AsString() != "" {
		opts.State = state.AsString()
	}
	if assignee := args.GetAttr("assignee"); !assignee.IsNull() && assignee.AsString() != "" {
		opts.Assignee = assignee.AsString()
	}
	if creator := args.GetAttr("creator"); !creator.IsNull() && creator.AsString() != "" {
		opts.Creator = creator.AsString()
	}
	if mentioned := args.GetAttr("mentioned"); !mentioned.IsNull() && mentioned.AsString() != "" {
		opts.Mentioned = mentioned.AsString()
	}
	if labels := args.GetAttr("labels"); !labels.IsNull() && labels.LengthInt() > 0 {
		arr := make([]string, labels.LengthInt())
		for i, label := range labels.AsValueSlice() {
			arr[i] = label.AsString()
		}
		opts.Labels = arr
	}
	if sort := args.GetAttr("sort"); !sort.IsNull() && sort.AsString() != "" {
		opts.Sort = sort.AsString()
	}
	if direction := args.GetAttr("direction"); !direction.IsNull() && direction.AsString() != "" {
		opts.Direction = direction.AsString()
	}
	if since := args.GetAttr("since"); !since.IsNull() && since.AsString() != "" {
		ts, err := time.Parse(time.RFC3339, since.AsString())
		if err != nil {
			return nil, fmt.Errorf("since must be in RFC3339 format")
		}
		opts.Since = ts
	}
	parsed := &parsedIssuesArgs{
		owner: owner,
		name:  name,
		opts:  opts,
		limit: -1,
	}
	if limit := args.GetAttr("limit"); !limit.IsNull() && limit.AsBigFloat().IsInt() {
		parsed.limit, _ = limit.AsBigFloat().Int64()
		if parsed.limit <= 0 {
			return nil, fmt.Errorf("limit must be greater than 0")
		}
	}
	return parsed, nil
}

func encodeGithubIssue(issue *gh.Issue) (plugin.Data, error) {
	raw, err := json.Marshal(issue)
	if err != nil {
		return nil, fmt.Errorf("failed to encode issue: %w", err)
	}
	return plugin.UnmarshalJSONData(raw)
}
