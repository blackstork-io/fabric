package github

import (
	"testing"
	"time"

	"github.com/blackstork-io/fabric/plugininterface/v1"
	"github.com/blackstork-io/fabric/plugins/data/github/internal/mocks"
	gh "github.com/google/go-github/v58/github"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
	"github.com/zclconf/go-cty/cty"
)

type PluginTestSuite struct {
	suite.Suite
	plugin plugininterface.PluginRPC
	cli    *mocks.Client
}

func TestPluginSuite(t *testing.T) {
	suite.Run(t, &PluginTestSuite{})
}

func (s *PluginTestSuite) SetupSuite() {
	s.plugin = Plugin{
		Loader: func(token string) Client {
			return s.cli
		},
	}
}

func (s *PluginTestSuite) SetupTest() {
	s.cli = &mocks.Client{}
}

func (s *PluginTestSuite) TearDownTest() {
	s.cli.AssertExpectations(s.T())
}

func (s *PluginTestSuite) TestGetPlugins() {
	plugins := s.plugin.GetPlugins()
	s.Require().Len(plugins, 1, "expected 1 plugin")
	got := plugins[0]
	s.Equal("github_issues", got.Name)
	s.Equal("data", got.Kind)
	s.Equal("blackstork", got.Namespace)
	s.Equal(Version.String(), got.Version.Cast().String())
	s.NotNil(got.ConfigSpec)
	s.NotNil(got.InvocationSpec)
}

func int64ptr(i int64) *int64 { return &i }

func (s *PluginTestSuite) TestCallBasic() {
	s.cli.On("ListByRepo", mock.Anything, "testorg", "testrepo", &gh.IssueListByRepoOptions{
		ListOptions: gh.ListOptions{
			PerPage: 30,
			Page:    1,
		},
	}).
		Return([]*gh.Issue{
			{
				ID: int64ptr(123),
			},
			{
				ID: int64ptr(124),
			},
		}, &gh.Response{}, nil)

	args := plugininterface.Args{
		Kind: "data",
		Name: "github_issues",
		Config: cty.ObjectVal(map[string]cty.Value{
			"github_token": cty.StringVal("testtoken"),
		}),
		Args: cty.ObjectVal(map[string]cty.Value{
			"repository": cty.StringVal("testorg/testrepo"),
			"limit":      cty.NullVal(cty.Number),
			"milestone":  cty.NullVal(cty.String),
			"state":      cty.NullVal(cty.String),
			"assignee":   cty.NullVal(cty.String),
			"creator":    cty.NullVal(cty.String),
			"mentioned":  cty.NullVal(cty.String),
			"labels":     cty.ListValEmpty(cty.String),
			"sort":       cty.NullVal(cty.String),
			"direction":  cty.NullVal(cty.String),
			"since":      cty.NullVal(cty.String),
		}),
	}
	result := s.plugin.Call(args)
	expected := plugininterface.Result{
		Result: []any{
			&gh.Issue{
				ID: int64ptr(123),
			},
			&gh.Issue{
				ID: int64ptr(124),
			},
		},
	}
	s.Equal(expected, result)
}

func (s *PluginTestSuite) TestCallAdvanced() {
	since, err := time.Parse(time.RFC3339, "2021-01-01T00:00:00Z")
	s.Require().NoError(err)
	s.cli.On("ListByRepo", mock.Anything, "testorg", "testrepo", &gh.IssueListByRepoOptions{
		Milestone: "testmilestone",
		State:     "open",
		Assignee:  "testassignee",
		Creator:   "testcreator",
		Labels: []string{
			"testlabel1",
			"testlabel2",
		},
		Sort:      "created",
		Direction: "asc",
		Mentioned: "testmentioned",
		Since:     since,
		ListOptions: gh.ListOptions{
			PerPage: 30,
			Page:    1,
		},
	}).
		Return([]*gh.Issue{
			{
				ID: int64ptr(123),
			},
			{
				ID: int64ptr(124),
			},
			{
				ID: int64ptr(125),
			},
		}, &gh.Response{}, nil)

	args := plugininterface.Args{
		Kind: "data",
		Name: "github_issues",
		Config: cty.ObjectVal(map[string]cty.Value{
			"github_token": cty.StringVal("testtoken"),
		}),
		Args: cty.ObjectVal(map[string]cty.Value{
			"repository": cty.StringVal("testorg/testrepo"),
			"limit":      cty.NumberIntVal(2),
			"milestone":  cty.StringVal("testmilestone"),
			"state":      cty.StringVal("open"),
			"assignee":   cty.StringVal("testassignee"),
			"creator":    cty.StringVal("testcreator"),
			"mentioned":  cty.StringVal("testmentioned"),
			"labels": cty.ListVal([]cty.Value{
				cty.StringVal("testlabel1"),
				cty.StringVal("testlabel2"),
			}),
			"sort":      cty.StringVal("created"),
			"direction": cty.StringVal("asc"),
			"since":     cty.StringVal("2021-01-01T00:00:00Z"),
		}),
	}
	result := s.plugin.Call(args)
	expected := plugininterface.Result{
		Result: []any{
			&gh.Issue{
				ID: int64ptr(123),
			},
			&gh.Issue{
				ID: int64ptr(124),
			},
		},
	}
	s.Equal(expected, result)
}
