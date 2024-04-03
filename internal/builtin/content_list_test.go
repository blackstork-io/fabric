package builtin

import (
	"context"
	"testing"

	"github.com/hashicorp/hcl/v2"
	"github.com/stretchr/testify/suite"
	"github.com/zclconf/go-cty/cty"

	"github.com/blackstork-io/fabric/plugin"
)

type ListGeneratorTestSuite struct {
	suite.Suite
	schema *plugin.ContentProvider
}

func TestListGeneratorTestSuite(t *testing.T) {
	suite.Run(t, &ListGeneratorTestSuite{})
}

func (s *ListGeneratorTestSuite) SetupSuite() {
	s.schema = makeListContentProvider()
}

func (s *ListGeneratorTestSuite) TestSchema() {
	s.NotNil(s.schema)
	s.Nil(s.schema.Config)
	s.NotNil(s.schema.Args)
	s.NotNil(s.schema.ContentFunc)
}

func (s *ListGeneratorTestSuite) TestNilQueryResult() {
	args := cty.ObjectVal(map[string]cty.Value{
		"item_template": cty.StringVal("{{.}}"),
		"format":        cty.NullVal(cty.String),
	})
	ctx := context.Background()
	result, diags := s.schema.ContentFunc(ctx, &plugin.ProvideContentParams{
		Args:        args,
		DataContext: plugin.MapData{},
	})
	s.Nil(result)
	s.Equal(hcl.Diagnostics{{
		Severity: hcl.DiagError,
		Summary:  "Failed to render template",
		Detail:   "query_result is required in data context",
	}}, diags)
}

func (s *ListGeneratorTestSuite) TestNonArrayQueryResult() {
	args := cty.ObjectVal(map[string]cty.Value{
		"item_template": cty.StringVal("{{.}}"),
		"format":        cty.NullVal(cty.String),
	})
	ctx := context.Background()
	result, diags := s.schema.ContentFunc(ctx, &plugin.ProvideContentParams{
		Args: args,
		DataContext: plugin.MapData{
			"query_result": plugin.StringData("not_an_array"),
		},
	})
	s.Nil(result)
	s.Equal(hcl.Diagnostics{{
		Severity: hcl.DiagError,
		Summary:  "Failed to render template",
		Detail:   "query_result must be an array",
	}}, diags)
}

func (s *ListGeneratorTestSuite) TestUnordered() {
	args := cty.ObjectVal(map[string]cty.Value{
		"item_template": cty.StringVal("foo {{.}}"),
		"format":        cty.StringVal("unordered"),
	})
	ctx := context.Background()
	result, diags := s.schema.ContentFunc(ctx, &plugin.ProvideContentParams{
		Args: args,
		DataContext: plugin.MapData{
			"query_result": plugin.ListData{
				plugin.StringData("bar"),
				plugin.StringData("baz"),
			},
		},
	})
	s.Equal("* foo bar\n* foo baz\n", result.Content.Print())
	s.Empty(diags)
}

func (s *ListGeneratorTestSuite) TestOrdered() {
	args := cty.ObjectVal(map[string]cty.Value{
		"item_template": cty.StringVal("foo {{.}}"),
		"format":        cty.StringVal("ordered"),
	})
	ctx := context.Background()
	result, diags := s.schema.ContentFunc(ctx, &plugin.ProvideContentParams{
		Args: args,
		DataContext: plugin.MapData{
			"query_result": plugin.ListData{
				plugin.StringData("bar"),
				plugin.StringData("baz"),
			},
		},
	})
	s.Equal("1. foo bar\n2. foo baz\n", result.Content.Print())
	s.Empty(diags)
}

func (s *ListGeneratorTestSuite) TestTaskList() {
	args := cty.ObjectVal(map[string]cty.Value{
		"item_template": cty.StringVal("foo {{.}}"),
		"format":        cty.StringVal("tasklist"),
	})
	ctx := context.Background()
	result, diags := s.schema.ContentFunc(ctx, &plugin.ProvideContentParams{
		Args: args,
		DataContext: plugin.MapData{
			"query_result": plugin.ListData{
				plugin.StringData("bar"),
				plugin.StringData("baz"),
			},
		},
	})
	s.Equal("* [ ] foo bar\n* [ ] foo baz\n", result.Content.Print())
	s.Empty(diags)
}

func (s *ListGeneratorTestSuite) TestBasic() {
	args := cty.ObjectVal(map[string]cty.Value{
		"item_template": cty.StringVal("foo {{.}}"),
		"format":        cty.NullVal(cty.String),
	})
	ctx := context.Background()
	result, diags := s.schema.ContentFunc(ctx, &plugin.ProvideContentParams{
		Args: args,
		DataContext: plugin.MapData{
			"query_result": plugin.ListData{
				plugin.StringData("bar"),
				plugin.StringData("baz"),
			},
		},
	})
	s.Equal("* foo bar\n* foo baz\n", result.Content.Print())
	s.Empty(diags)
}

func (s *ListGeneratorTestSuite) TestAdvanced() {
	args := cty.ObjectVal(map[string]cty.Value{
		"item_template": cty.StringVal("foo {{.bar}} {{.baz | upper}}"),
		"format":        cty.NullVal(cty.String),
	})
	ctx := context.Background()
	result, diags := s.schema.ContentFunc(ctx, &plugin.ProvideContentParams{
		Args: args,
		DataContext: plugin.MapData{
			"query_result": plugin.ListData{
				plugin.MapData{
					"bar": plugin.StringData("bar1"),
					"baz": plugin.StringData("baz1"),
				},
				plugin.MapData{
					"bar": plugin.StringData("bar2"),
					"baz": plugin.StringData("baz2"),
				},
			},
		},
	})
	s.Equal("* foo bar1 BAZ1\n* foo bar2 BAZ2\n", result.Content.Print())
	s.Empty(diags)
}

func (s *ListGeneratorTestSuite) TestEmptyQueryResult() {
	args := cty.ObjectVal(map[string]cty.Value{
		"item_template": cty.StringVal("foo {{.}}"),
		"format":        cty.NullVal(cty.String),
	})
	ctx := context.Background()
	result, diags := s.schema.ContentFunc(ctx, &plugin.ProvideContentParams{
		Args: args,
		DataContext: plugin.MapData{
			"query_result": plugin.ListData{},
		},
	})
	s.Equal("", result.Content.Print())
	s.Empty(diags)
}

func (s *ListGeneratorTestSuite) TestMissingItemTemplate() {
	args := cty.ObjectVal(map[string]cty.Value{
		"item_template": cty.NullVal(cty.String),
		"format":        cty.NullVal(cty.String),
	})
	ctx := context.Background()
	result, diags := s.schema.ContentFunc(ctx, &plugin.ProvideContentParams{
		Args: args,
		DataContext: plugin.MapData{
			"query_result": plugin.ListData{},
		},
	})
	s.Nil(result)
	s.Equal(hcl.Diagnostics{{
		Severity: hcl.DiagError,
		Summary:  "Failed to parse template",
		Detail:   "item_template is required",
	}}, diags)
}

func (s *ListGeneratorTestSuite) TestInvalidFormat() {
	args := cty.ObjectVal(map[string]cty.Value{
		"item_template": cty.StringVal("foo {{.}}"),
		"format":        cty.StringVal("invalid"),
	})
	ctx := context.Background()
	result, diags := s.schema.ContentFunc(ctx, &plugin.ProvideContentParams{
		Args: args,
		DataContext: plugin.MapData{
			"query_result": plugin.ListData{},
		},
	})
	s.Nil(result)
	s.Equal(hcl.Diagnostics{{
		Severity: hcl.DiagError,
		Summary:  "Failed to parse template",
		Detail:   "invalid format: invalid",
	}}, diags)
}

func (s *ListGeneratorTestSuite) TestMissingDataContext() {
	args := cty.ObjectVal(map[string]cty.Value{
		"item_template": cty.StringVal("foo {{.}}"),
		"format":        cty.NullVal(cty.String),
	})
	ctx := context.Background()
	result, diags := s.schema.ContentFunc(ctx, &plugin.ProvideContentParams{
		Args: args,
	})
	s.Nil(result)
	s.Equal(hcl.Diagnostics{{
		Severity: hcl.DiagError,
		Summary:  "Failed to render template",
		Detail:   "data context is required",
	}}, diags)
}
