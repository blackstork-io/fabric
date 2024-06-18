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

type TableGeneratorTestSuite struct {
	suite.Suite
	schema *plugin.ContentProvider
}

func TestTableGeneratorTestSuite(t *testing.T) {
	suite.Run(t, &TableGeneratorTestSuite{})
}

func (s *TableGeneratorTestSuite) SetupSuite() {
	s.schema = makeTableContentProvider()
}

func (s *TableGeneratorTestSuite) TestSchema() {
	s.NotNil(s.schema)
	s.Nil(s.schema.Config)
	s.NotNil(s.schema.Args)
	s.NotNil(s.schema.ContentFunc)
}

func (s *TableGeneratorTestSuite) TestNilRowVars() {
	val := `
	columns = [
		{
			header = "{{.col_prefix}} Name"
			value  = "{{.name}}"
		},
		{
			header = "{{.col_prefix}} Age"
			value  = "{{.age}}"
		}
	]
	rows = null
	`
	dataCtx := plugin.MapData{
		"col_prefix": plugin.StringData("User"),
	}
	args := plugintest.DecodeAndAssert(s.T(), s.schema.Args, val, dataCtx, diagtest.Asserts{})
	ctx := context.Background()
	result, diags := s.schema.ContentFunc(ctx, &plugin.ProvideContentParams{
		Args:        args,
		DataContext: dataCtx,
	})
	diagtest.AssertNoErrors(s.T(), diags, nil)
	s.Equal("|User Name|User Age|\n|---|---|\n", mdprint.PrintString(result.Content))
	s.Nil(diags)
}

func (s *TableGeneratorTestSuite) TestEmptyQueryResult() {
	val := `
	columns = [
		{
			header = "{{.col_prefix}} Name"
			value  = "{{.name}}"
		},
		{
			header = "{{.col_prefix}} Age"
			value  = "{{.age}}"
		}
	]
	rows = null
	`

	dataCtx := plugin.MapData{
		"col_prefix": plugin.StringData("User"),
	}
	args := plugintest.DecodeAndAssert(s.T(), s.schema.Args, val, dataCtx, diagtest.Asserts{})
	ctx := context.Background()
	result, diags := s.schema.ContentFunc(ctx, &plugin.ProvideContentParams{
		Args:        args,
		DataContext: dataCtx,
	})
	s.Equal("|User Name|User Age|\n|---|---|\n", mdprint.PrintString(result.Content))
	s.Nil(diags)
}

func (s *TableGeneratorTestSuite) TestBasic() {
	val := `
	columns = [
		{
			header = "{{.col_prefix}} Name"
			value  = "{{.row.value.name}}"
		},
		{
			header = "{{.col_prefix}} Age"
			value  = "{{.row.value.age}}"
		}
	]
	rows = [
		{name = "John", age = 42},
		{name = "Jane", age = 43}
	]
	`
	dataCtx := plugin.MapData{
		"col_prefix": plugin.StringData("User"),
	}
	args := plugintest.DecodeAndAssert(s.T(), s.schema.Args, val, dataCtx, diagtest.Asserts{})
	ctx := context.Background()
	result, diags := s.schema.ContentFunc(ctx, &plugin.ProvideContentParams{
		Args:        args,
		DataContext: dataCtx,
	})
	s.Equal("|User Name|User Age|\n|---|---|\n|John|42|\n|Jane|43|\n", mdprint.PrintString(result.Content))
	s.Nil(diags)
}

func (s *TableGeneratorTestSuite) TestSprigTemplate() {
	val := `
	columns = [
		{
			header = "{{.col_prefix | upper}} Name"
			value  = "{{.row.value.name | upper}}"
		},
		{
			header = "{{.col_prefix}} Age"
			value  = "{{.row.value.age}}"
		}
	]
	rows = [
		{name = "John", age = 42},
		{name = "Jane", age = 43}
	]
	`
	dataCtx := plugin.MapData{
		"col_prefix": plugin.StringData("User"),
	}
	args := plugintest.DecodeAndAssert(s.T(), s.schema.Args, val, dataCtx, diagtest.Asserts{})
	ctx := context.Background()
	result, diags := s.schema.ContentFunc(ctx, &plugin.ProvideContentParams{
		Args:        args,
		DataContext: dataCtx,
	})
	s.Equal("|USER Name|User Age|\n|---|---|\n|JOHN|42|\n|JANE|43|\n", mdprint.PrintString(result.Content))
	s.Nil(diags)
}

func (s *TableGeneratorTestSuite) TestMissingHeader() {
	val := cty.ObjectVal(map[string]cty.Value{
		"columns": cty.ListVal([]cty.Value{
			cty.ObjectVal(map[string]cty.Value{
				"value": cty.StringVal("{{.name}}"),
			}),
			cty.ObjectVal(map[string]cty.Value{
				"value": cty.StringVal("{{.age}}"),
			}),
		}),
	})
	plugintest.ReencodeCTY(s.T(), s.schema.Args, val, diagtest.Asserts{{
		diagtest.IsError,
		diagtest.DetailContains("argument", "header", "required"),
	}})
}

func (s *TableGeneratorTestSuite) TestNilHeader() {
	val := cty.ObjectVal(map[string]cty.Value{
		"columns": cty.ListVal([]cty.Value{
			cty.ObjectVal(map[string]cty.Value{
				"header": cty.NullVal(cty.String),
				"value":  cty.StringVal("{{.name}}"),
			}),
			cty.ObjectVal(map[string]cty.Value{
				"header": cty.NullVal(cty.String),
				"value":  cty.StringVal("{{.age}}"),
			}),
		}),
	})
	args := plugintest.ReencodeCTY(s.T(), s.schema.Args, val, nil)
	ctx := context.Background()
	result, diags := s.schema.ContentFunc(ctx, &plugin.ProvideContentParams{
		Args: args,
		DataContext: plugin.MapData{
			"col_prefix": plugin.StringData("User"),
		},
	})
	s.Nil(result)
	s.Equal(diagnostics.Diag{{
		Severity: hcl.DiagError,
		Summary:  "Failed to parse arguments",
		Detail:   "missing header in table cell",
	}}, diags)
}

func (s *TableGeneratorTestSuite) TestNilValue() {
	val := cty.ObjectVal(map[string]cty.Value{
		"columns": cty.ListVal([]cty.Value{
			cty.ObjectVal(map[string]cty.Value{
				"header": cty.StringVal("{{.col_prefix}} Name"),
				"value":  cty.NullVal(cty.String),
			}),
			cty.ObjectVal(map[string]cty.Value{
				"header": cty.StringVal("{{.col_prefix}} Age"),
				"value":  cty.NullVal(cty.String),
			}),
		}),
	})
	args := plugintest.ReencodeCTY(s.T(), s.schema.Args, val, nil)
	ctx := context.Background()
	result, diags := s.schema.ContentFunc(ctx, &plugin.ProvideContentParams{
		Args: args,
		DataContext: plugin.MapData{
			"col_prefix": plugin.StringData("User"),
		},
	})
	s.Nil(result)
	s.Equal(diagnostics.Diag{{
		Severity: hcl.DiagError,
		Summary:  "Failed to parse arguments",
		Detail:   "missing value in table cell",
	}}, diags)
}

func (s *TableGeneratorTestSuite) TestNilColumns() {
	val := cty.ObjectVal(map[string]cty.Value{
		"columns": cty.NullVal(cty.List(cty.Object(map[string]cty.Type{}))),
	})
	plugintest.ReencodeCTY(s.T(), s.schema.Args, val, diagtest.Asserts{{
		diagtest.IsError,
		diagtest.SummaryContains("Attribute must be non-null"),
	}})
}

func (s *TableGeneratorTestSuite) TestEmptyColumns() {
	val := `
	columns = []
	rows = [
		{name = "John", age = 42},
		{name = "Jane", age = 43}
	]
	`
	dataCtx := plugin.MapData{
		"col_prefix": plugin.StringData("User"),
	}
	plugintest.DecodeAndAssert(s.T(), s.schema.Args, val, dataCtx, diagtest.Asserts{{
		diagtest.IsError,
		diagtest.DetailContains("length", `"columns"`, ">= 1"),
	}})
}

func (s *TableGeneratorTestSuite) TestInvalidHeaderTemplate() {
	val := cty.ObjectVal(map[string]cty.Value{
		"columns": cty.ListVal([]cty.Value{
			cty.ObjectVal(map[string]cty.Value{
				"header": cty.StringVal("{{.col_prefix} Name"),
				"value":  cty.StringVal("{{.name}}"),
			}),
			cty.ObjectVal(map[string]cty.Value{
				"header": cty.StringVal("{{.col_prefix}} Age"),
				"value":  cty.StringVal("{{.age}}"),
			}),
		}),
	})
	args := plugintest.ReencodeCTY(s.T(), s.schema.Args, val, nil)
	ctx := context.Background()
	result, diags := s.schema.ContentFunc(ctx, &plugin.ProvideContentParams{
		Args: args,
		DataContext: plugin.MapData{
			"col_prefix": plugin.StringData("User"),
		},
	})
	s.Nil(result)
	s.Equal(diagnostics.Diag{{
		Severity: hcl.DiagError,
		Summary:  "Failed to parse arguments",
		Detail:   "failed to parse header template: template: header:1: bad character U+007D '}'",
	}}, diags)
}

func (s *TableGeneratorTestSuite) TestInvalidValueTemplate() {
	val := cty.ObjectVal(map[string]cty.Value{
		"columns": cty.ListVal([]cty.Value{
			cty.ObjectVal(map[string]cty.Value{
				"header": cty.StringVal("{{.col_prefix}} Name"),
				"value":  cty.StringVal("{{.name}"),
			}),
			cty.ObjectVal(map[string]cty.Value{
				"header": cty.StringVal("{{.col_prefix}} Age"),
				"value":  cty.StringVal("{{.age}}"),
			}),
		}),
	})
	args := plugintest.ReencodeCTY(s.T(), s.schema.Args, val, nil)
	ctx := context.Background()
	result, diags := s.schema.ContentFunc(ctx, &plugin.ProvideContentParams{
		Args: args,
		DataContext: plugin.MapData{
			"col_prefix": plugin.StringData("User"),
		},
	})
	s.Nil(result)
	s.Equal(diagnostics.Diag{{
		Severity: hcl.DiagError,
		Summary:  "Failed to parse arguments",
		Detail:   "failed to parse value template: template: value:1: bad character U+007D '}'",
	}}, diags)
}
