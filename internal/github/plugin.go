package github

import (
	"context"

	gh "github.com/google/go-github/v58/github"

	"github.com/blackstork-io/fabric/plugin"
)

var DefaultClientLoader = func(token string) Client {
	return gh.NewClient(nil).WithAuthToken(token).Issues
}

const (
	minPage  = 1
	pageSize = 30
)

type ClientLoaderFn func(token string) Client

type Client interface {
	ListByRepo(ctx context.Context, owner, repo string, opts *gh.IssueListByRepoOptions) ([]*gh.Issue, *gh.Response, error)
}

func Plugin(version string, clientLoader ClientLoaderFn) *plugin.Schema {
	if clientLoader == nil {
		clientLoader = DefaultClientLoader
	}
	return &plugin.Schema{
		Name:    "blackstork/github",
		Version: version,
		DataSources: plugin.DataSources{
			"github_issues": makeGithubIssuesDataSchema(clientLoader),
		},
	}
}
