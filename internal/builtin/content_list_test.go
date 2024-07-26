package builtin

import (
	"context"
	"testing"

	"github.com/hashicorp/hcl/v2"
	"github.com/stretchr/testify/suite"

	"github.com/blackstork-io/fabric/pkg/diagnostics"
	"github.com/blackstork-io/fabric/pkg/diagnostics/diagtest"
	"github.com/blackstork-io/fabric/plugin"
	"github.com/blackstork-io/fabric/plugin/plugindata"
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
	dataCtx := plugindata.Map{}
	args := plugintest.DecodeAndAssert(s.T(), s.schema.Args, `
		items = null
	`, dataCtx, diagtest.Asserts{})

	ctx := context.Background()
	result, diags := s.schema.ContentFunc(ctx, &plugin.ProvideContentParams{
		Args:        args,
		DataContext: dataCtx,
	})
	s.Nil(result)
	s.Equal(diagnostics.Diag{{
		Severity: hcl.DiagError,
		Summary:  "Failed to parse arguments",
		Detail:   "Data is nil",
	}}, diags)
}

func (s *ListGeneratorTestSuite) TestNonArrayQueryResult() {
	dataCtx := plugindata.Map{}
	args := plugintest.DecodeAndAssert(s.T(), s.schema.Args, `
		items = "not_an_array"
	`, dataCtx, diagtest.Asserts{})

	ctx := context.Background()
	result, diags := s.schema.ContentFunc(ctx, &plugin.ProvideContentParams{
		Args:        args,
		DataContext: dataCtx,
	})
	s.Nil(result)
	s.Equal(diagnostics.Diag{{
		Severity: hcl.DiagError,
		Summary:  "Failed to parse arguments",
		Detail:   "Data must be a list",
	}}, diags)
}

func (s *ListGeneratorTestSuite) TestUnordered() {
	dataCtx := plugindata.Map{}
	args := plugintest.DecodeAndAssert(s.T(), s.schema.Args, `
		items = ["bar", "baz"]
		item_template = "foo {{.}}"
		format = "unordered"
	`, dataCtx, diagtest.Asserts{})

	ctx := context.Background()
	result, diags := s.schema.ContentFunc(ctx, &plugin.ProvideContentParams{
		Args:        args,
		DataContext: dataCtx,
	})
	s.Equal("* foo bar\n* foo baz\n", mdprint.PrintString(result.Content))
	s.Empty(diags)
}

func (s *ListGeneratorTestSuite) TestOrdered() {
	dataCtx := plugindata.Map{}
	args := plugintest.DecodeAndAssert(s.T(), s.schema.Args, `
		items = ["bar", "baz"]
		item_template = "foo {{.}}"
		format = "ordered"
	`, dataCtx, diagtest.Asserts{})

	ctx := context.Background()
	result, diags := s.schema.ContentFunc(ctx, &plugin.ProvideContentParams{
		Args:        args,
		DataContext: dataCtx,
	})
	s.Equal("1. foo bar\n2. foo baz\n", mdprint.PrintString(result.Content))
	s.Empty(diags)
}

func (s *ListGeneratorTestSuite) TestTaskList() {
	dataCtx := plugindata.Map{}
	args := plugintest.DecodeAndAssert(s.T(), s.schema.Args, `
		items = ["bar", "baz"]
		item_template = "foo {{.}}"
		format = "tasklist"
	`, dataCtx, diagtest.Asserts{})

	ctx := context.Background()
	result, diags := s.schema.ContentFunc(ctx, &plugin.ProvideContentParams{
		Args:        args,
		DataContext: dataCtx,
	})
	s.Equal("* [ ] foo bar\n* [ ] foo baz\n", mdprint.PrintString(result.Content))
	s.Empty(diags)
}

func (s *ListGeneratorTestSuite) TestBasic() {
	dataCtx := plugindata.Map{}
	args := plugintest.DecodeAndAssert(s.T(), s.schema.Args, `
		items = ["bar", "baz"]
		item_template = "foo {{.}}"
	`, dataCtx, diagtest.Asserts{})

	ctx := context.Background()
	result, diags := s.schema.ContentFunc(ctx, &plugin.ProvideContentParams{
		Args:        args,
		DataContext: dataCtx,
	})
	s.Equal("* foo bar\n* foo baz\n", mdprint.PrintString(result.Content))
	s.Empty(diags)
}

func (s *ListGeneratorTestSuite) TestAdvanced() {
	dataCtx := plugindata.Map{}
	args := plugintest.DecodeAndAssert(s.T(), s.schema.Args, `
		items = [
			{
				bar = "bar1",
				baz = "baz1",
			},
			{
				bar = "bar2",
				baz = "baz2",
			}
		]
		item_template = "foo {{.bar}} {{.baz | upper}}"
	`, dataCtx, diagtest.Asserts{})

	ctx := context.Background()
	result, diags := s.schema.ContentFunc(ctx, &plugin.ProvideContentParams{
		Args:        args,
		DataContext: dataCtx,
	})
	s.Equal("* foo bar1 BAZ1\n* foo bar2 BAZ2\n", mdprint.PrintString(result.Content))
	s.Empty(diags)
}

func (s *ListGeneratorTestSuite) TestEmptyQueryResult() {
	dataCtx := plugindata.Map{}
	args := plugintest.DecodeAndAssert(s.T(), s.schema.Args, `
		items = []
		item_template = "foo {{.}}"
	`, dataCtx, diagtest.Asserts{})

	ctx := context.Background()
	result, diags := s.schema.ContentFunc(ctx, &plugin.ProvideContentParams{
		Args:        args,
		DataContext: dataCtx,
	})
	s.Equal("", mdprint.PrintString(result.Content))
	s.Empty(diags)
}

func (s *ListGeneratorTestSuite) TestMissingItemTemplate() {
	dataCtx := plugindata.Map{}
	args := plugintest.DecodeAndAssert(s.T(), s.schema.Args, `
		items = ["bar", "baz"]
	`, dataCtx, diagtest.Asserts{})

	ctx := context.Background()
	result, diags := s.schema.ContentFunc(ctx, &plugin.ProvideContentParams{
		Args:        args,
		DataContext: dataCtx,
	})
	s.Equal("* bar\n* baz\n", mdprint.PrintString(result.Content))
	s.Empty(diags)
}
