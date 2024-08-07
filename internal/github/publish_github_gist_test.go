package github_test

import (
	"context"
	"testing"

	gh "github.com/google/go-github/v58/github"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
	"github.com/zclconf/go-cty/cty"

	"github.com/blackstork-io/fabric/internal/github"
	github_mocks "github.com/blackstork-io/fabric/mocks/internalpkg/github"
	"github.com/blackstork-io/fabric/plugin"
)

type GithubPublishGistTestSuite struct {
	suite.Suite
	plugin  *plugin.Schema
	cli     *github_mocks.Client
	gistCli *github_mocks.GistClient
}

func TestGithubublishGistSuite(t *testing.T) {
	suite.Run(t, &GithubPublishGistTestSuite{})
}

func (s *GithubPublishGistTestSuite) SetupSuite() {
	s.plugin = github.Plugin("1.2.3", func(token string) github.Client {
		return s.cli
	})
}

func (s *GithubPublishGistTestSuite) SetupTest() {
	s.cli = &github_mocks.Client{}
	s.gistCli = &github_mocks.GistClient{}
}

func (s *GithubPublishGistTestSuite) TearDownTest() {
	s.cli.AssertExpectations(s.T())
}

func (s *GithubPublishGistTestSuite) TestSchema() {
	schema := s.plugin.Publishers["github_gist"]
	s.Require().NotNil(schema)
	s.NotNil(schema.Config)
	s.NotNil(schema.Args)
	s.NotNil(schema.PublishFunc)
}

func (s *GithubPublishGistTestSuite) TestBasic() {
	s.cli.On("Gists").Return(s.gistCli)
	s.gistCli.On("Create", mock.Anything, &gh.Gist{
		Description: gh.String("test description"),
		Public:      gh.Bool(false),
		Files: map[gh.GistFilename]gh.GistFile{
			"filename.md": {
				Content:  gh.String("# Header 1\n\nLorem ipsum dolor sit amet, consectetur adipiscing elit."),
				Filename: gh.String("filename.md"),
			},
		},
	}).Return(&gh.Gist{HTMLURL: gh.String("http://gist.github.com/mock")}, &gh.Response{}, nil)
	ctx := context.Background()
	titleMeta := plugin.MapData{
		"provider": plugin.StringData("title"),
		"plugin":   plugin.StringData("blackstork/builtin"),
	}
	diags := s.plugin.Publish(ctx, "github_gist", &plugin.PublishParams{
		Config: cty.ObjectVal(map[string]cty.Value{
			"github_token": cty.StringVal("testtoken"),
		}),
		Args: cty.ObjectVal(map[string]cty.Value{
			"description": cty.StringVal("test description"),
			"filename":    cty.StringVal("filename.md"),
			"make_public": cty.NullVal(cty.Bool),
			"gist_id":     cty.NullVal(cty.String),
		}),
		Format: plugin.OutputFormatMD,
		DataContext: plugin.MapData{
			"document": plugin.MapData{
				"meta": plugin.MapData{
					"name": plugin.StringData("test_document"),
				},
				"content": plugin.MapData{
					"type": plugin.StringData("section"),
					"children": plugin.ListData{
						plugin.MapData{
							"type":     plugin.StringData("element"),
							"markdown": plugin.StringData("# Header 1"),
							"meta":     titleMeta,
						},
						plugin.MapData{
							"type":     plugin.StringData("element"),
							"markdown": plugin.StringData("Lorem ipsum dolor sit amet, consectetur adipiscing elit."),
						},
					},
				},
			},
		},
		DocumentName: "test_doc",
	})
	s.Require().Nil(diags)
}

func (s *GithubPublishGistTestSuite) TestFilenameOptional() {
	s.cli.On("Gists").Return(s.gistCli)
	s.gistCli.On("Create", mock.Anything, &gh.Gist{
		Description: gh.String("test description"),
		Public:      gh.Bool(false),
		Files: map[gh.GistFilename]gh.GistFile{
			"test_doc.md": {
				Content:  gh.String("# Header 1\n\nLorem ipsum dolor sit amet, consectetur adipiscing elit."),
				Filename: gh.String("test_doc.md"),
			},
		},
	}).Return(&gh.Gist{HTMLURL: gh.String("http://gist.github.com/mock")}, &gh.Response{}, nil)
	ctx := context.Background()
	titleMeta := plugin.MapData{
		"provider": plugin.StringData("title"),
		"plugin":   plugin.StringData("blackstork/builtin"),
	}
	diags := s.plugin.Publish(ctx, "github_gist", &plugin.PublishParams{
		Config: cty.ObjectVal(map[string]cty.Value{
			"github_token": cty.StringVal("testtoken"),
		}),
		Args: cty.ObjectVal(map[string]cty.Value{
			"description": cty.StringVal("test description"),
			"filename":    cty.StringVal(""),
			"make_public": cty.NullVal(cty.Bool),
			"gist_id":     cty.NullVal(cty.String),
		}),
		Format: plugin.OutputFormatMD,
		DataContext: plugin.MapData{
			"document": plugin.MapData{
				"meta": plugin.MapData{
					"name": plugin.StringData("test_document"),
				},
				"content": plugin.MapData{
					"type": plugin.StringData("section"),
					"children": plugin.ListData{
						plugin.MapData{
							"type":     plugin.StringData("element"),
							"markdown": plugin.StringData("# Header 1"),
							"meta":     titleMeta,
						},
						plugin.MapData{
							"type":     plugin.StringData("element"),
							"markdown": plugin.StringData("Lorem ipsum dolor sit amet, consectetur adipiscing elit."),
						},
					},
				},
			},
		},
		DocumentName: "test_doc",
	})
	s.Require().Nil(diags)
}

func (s *GithubPublishGistTestSuite) TestResetOldFiles() {
	s.cli.On("Gists").Return(s.gistCli)
	s.gistCli.On("Get", mock.Anything, "gistid").Return(&gh.Gist{
		Files: map[gh.GistFilename]gh.GistFile{
			"oldfile.md": {
				Filename: gh.String("oldfile.md"),
				Content:  gh.String("old content"),
			},
		},
	}, &gh.Response{}, nil)
	s.gistCli.On("Edit", mock.Anything, "gistid", &gh.Gist{
		Description: gh.String("test description"),
		Public:      gh.Bool(false),
		Files: map[gh.GistFilename]gh.GistFile{
			"oldfile.md": {},
			"test_doc.md": {
				Content:  gh.String("# Header 1\n\nLorem ipsum dolor sit amet, consectetur adipiscing elit."),
				Filename: gh.String("test_doc.md"),
			},
		},
	}).Return(&gh.Gist{HTMLURL: gh.String("http://gist.github.com/mock")}, &gh.Response{}, nil)
	ctx := context.Background()
	titleMeta := plugin.MapData{
		"provider": plugin.StringData("title"),
		"plugin":   plugin.StringData("blackstork/builtin"),
	}
	diags := s.plugin.Publish(ctx, "github_gist", &plugin.PublishParams{
		Config: cty.ObjectVal(map[string]cty.Value{
			"github_token": cty.StringVal("testtoken"),
		}),
		Args: cty.ObjectVal(map[string]cty.Value{
			"description": cty.StringVal("test description"),
			"filename":    cty.StringVal(""),
			"make_public": cty.NullVal(cty.Bool),
			"gist_id":     cty.StringVal("gistid"),
		}),
		Format: plugin.OutputFormatMD,
		DataContext: plugin.MapData{
			"document": plugin.MapData{
				"meta": plugin.MapData{
					"name": plugin.StringData("test_document"),
				},
				"content": plugin.MapData{
					"type": plugin.StringData("section"),
					"children": plugin.ListData{
						plugin.MapData{
							"type":     plugin.StringData("element"),
							"markdown": plugin.StringData("# Header 1"),
							"meta":     titleMeta,
						},
						plugin.MapData{
							"type":     plugin.StringData("element"),
							"markdown": plugin.StringData("Lorem ipsum dolor sit amet, consectetur adipiscing elit."),
						},
					},
				},
			},
		},
		DocumentName: "test_doc",
	})
	s.Require().Nil(diags)
}
