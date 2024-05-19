package builtin

import (
	"context"
	"testing"

	"github.com/hashicorp/hcl/v2"
	"github.com/stretchr/testify/suite"
	"github.com/zclconf/go-cty/cty"

	"github.com/blackstork-io/fabric/pkg/diagnostics"
	"github.com/blackstork-io/fabric/pkg/diagnostics/diagtest"
	"github.com/blackstork-io/fabric/plugin"
	"github.com/blackstork-io/fabric/plugin/plugintest"
	"github.com/blackstork-io/fabric/print/mdprint"
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
	val := cty.ObjectVal(map[string]cty.Value{
		"item_template": cty.StringVal("{{.}}"),
		"format":        cty.NullVal(cty.String),
	})
	args := plugintest.ReencodeCTY(s.T(), s.schema.Args, val, nil)

	ctx := context.Background()
	result, diags := s.schema.ContentFunc(ctx, &plugin.ProvideContentParams{
		Args:        args,
		DataContext: plugin.MapData{},
	})
	s.Nil(result)
	s.Equal(diagnostics.Diag{{
		Severity: hcl.DiagError,
		Summary:  "Failed to render template",
		Detail:   "query_result is required in data context",
	}}, diags)
}

func (s *ListGeneratorTestSuite) TestNonArrayQueryResult() {
	val := cty.ObjectVal(map[string]cty.Value{
		"item_template": cty.StringVal("{{.}}"),
		"format":        cty.NullVal(cty.String),
	})
	args := plugintest.ReencodeCTY(s.T(), s.schema.Args, val, nil)
	ctx := context.Background()
	result, diags := s.schema.ContentFunc(ctx, &plugin.ProvideContentParams{
		Args: args,
		DataContext: plugin.MapData{
			"query_result": plugin.StringData("not_an_array"),
		},
	})
	s.Nil(result)
	s.Equal(diagnostics.Diag{{
		Severity: hcl.DiagError,
		Summary:  "Failed to render template",
		Detail:   "query_result must be an array",
	}}, diags)
}

func (s *ListGeneratorTestSuite) TestUnordered() {
	val := cty.ObjectVal(map[string]cty.Value{
		"item_template": cty.StringVal("foo {{.}}"),
		"format":        cty.StringVal("unordered"),
	})
	args := plugintest.ReencodeCTY(s.T(), s.schema.Args, val, nil)
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
	s.Equal("* foo bar\n* foo baz\n", mdprint.PrintString(result.Content))
	s.Empty(diags)
}

func (s *ListGeneratorTestSuite) TestOrdered() {
	val := cty.ObjectVal(map[string]cty.Value{
		"item_template": cty.StringVal("foo {{.}}"),
		"format":        cty.StringVal("ordered"),
	})
	args := plugintest.ReencodeCTY(s.T(), s.schema.Args, val, nil)
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
	s.Equal("1. foo bar\n2. foo baz\n", mdprint.PrintString(result.Content))
	s.Empty(diags)
}

func (s *ListGeneratorTestSuite) TestTaskList() {
	val := cty.ObjectVal(map[string]cty.Value{
		"item_template": cty.StringVal("foo {{.}}"),
		"format":        cty.StringVal("tasklist"),
	})
	args := plugintest.ReencodeCTY(s.T(), s.schema.Args, val, nil)
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
	s.Equal("* [ ] foo bar\n* [ ] foo baz\n", mdprint.PrintString(result.Content))
	s.Empty(diags)
}

func (s *ListGeneratorTestSuite) TestBasic() {
	val := cty.ObjectVal(map[string]cty.Value{
		"item_template": cty.StringVal("foo {{.}}"),
		"format":        cty.NullVal(cty.String),
	})
	args := plugintest.ReencodeCTY(s.T(), s.schema.Args, val, nil)
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
	s.Equal("* foo bar\n* foo baz\n", mdprint.PrintString(result.Content))
	s.Empty(diags)
}

func (s *ListGeneratorTestSuite) TestAdvanced() {
	val := cty.ObjectVal(map[string]cty.Value{
		"item_template": cty.StringVal("foo {{.bar}} {{.baz | upper}}"),
		"format":        cty.NullVal(cty.String),
	})
	args := plugintest.ReencodeCTY(s.T(), s.schema.Args, val, nil)
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
	s.Equal("* foo bar1 BAZ1\n* foo bar2 BAZ2\n", mdprint.PrintString(result.Content))
	s.Empty(diags)
}

func (s *ListGeneratorTestSuite) TestEmptyQueryResult() {
	val := cty.ObjectVal(map[string]cty.Value{
		"item_template": cty.StringVal("foo {{.}}"),
		"format":        cty.NullVal(cty.String),
	})
	args := plugintest.ReencodeCTY(s.T(), s.schema.Args, val, nil)
	ctx := context.Background()
	result, diags := s.schema.ContentFunc(ctx, &plugin.ProvideContentParams{
		Args: args,
		DataContext: plugin.MapData{
			"query_result": plugin.ListData{},
		},
	})
	s.Equal("", mdprint.PrintString(result.Content))
	s.Empty(diags)
}

func (s *ListGeneratorTestSuite) TestMissingItemTemplate() {
	val := cty.ObjectVal(map[string]cty.Value{
		"item_template": cty.NullVal(cty.String),
		"format":        cty.NullVal(cty.String),
	})
	plugintest.ReencodeCTY(s.T(), s.schema.Args, val, diagtest.Asserts{{
		diagtest.IsError,
		diagtest.SummaryContains("Attribute must be non-null"),
	}})
}

func (s *ListGeneratorTestSuite) TestInvalidFormat() {
	val := cty.ObjectVal(map[string]cty.Value{
		"item_template": cty.StringVal("foo {{.}}"),
		"format":        cty.StringVal("invalid"),
	})
	plugintest.ReencodeCTY(s.T(), s.schema.Args, val, diagtest.Asserts{{
		diagtest.IsError,
		diagtest.SummaryContains("Attribute is not one of the allowed values"),
	}})
}

func (s *ListGeneratorTestSuite) TestMissingDataContext() {
	val := cty.ObjectVal(map[string]cty.Value{
		"item_template": cty.StringVal("foo {{.}}"),
		"format":        cty.NullVal(cty.String),
	})
	args := plugintest.ReencodeCTY(s.T(), s.schema.Args, val, nil)
	ctx := context.Background()
	result, diags := s.schema.ContentFunc(ctx, &plugin.ProvideContentParams{
		Args: args,
	})
	s.Nil(result)
	s.Equal(diagnostics.Diag{{
		Severity: hcl.DiagError,
		Summary:  "Failed to render template",
		Detail:   "data context is required",
	}}, diags)
}
