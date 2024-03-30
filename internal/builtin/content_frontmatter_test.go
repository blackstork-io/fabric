package builtin

import (
	"context"
	"testing"

	"github.com/hashicorp/hcl/v2"
	"github.com/stretchr/testify/suite"
	"github.com/zclconf/go-cty/cty"

	"github.com/blackstork-io/fabric/plugin"
)

type FrontMatterGeneratorTestSuite struct {
	suite.Suite
	schema *plugin.Schema
}

func TestFrontMatterGeneratorSuite(t *testing.T) {
	suite.Run(t, &FrontMatterGeneratorTestSuite{})
}

func (s *FrontMatterGeneratorTestSuite) SetupSuite() {
	s.schema = Plugin("")
}

func (s *FrontMatterGeneratorTestSuite) TestSchema() {
	provider := s.schema.ContentProviders["frontmatter"]
	s.NotNil(provider)
	s.Nil(provider.Config)
	s.NotNil(provider.Args)
	s.NotNil(provider.ContentFunc)
}

func (s *FrontMatterGeneratorTestSuite) TestInvalidFormat() {
	args := cty.ObjectVal(map[string]cty.Value{
		"format": cty.StringVal("invalid_type"),
		"content": cty.ObjectVal(map[string]cty.Value{
			"foo": cty.StringVal("bar"),
		}),
	})
	ctx := context.Background()
	document := plugin.ContentSection{}
	content, diags := s.schema.ProvideContent(ctx, "frontmatter", &plugin.ProvideContentParams{
		Args: args,
		DataContext: plugin.MapData{
			"document": plugin.MapData{
				"content": document.AsData(),
			},
		},
	})
	s.Nil(content)
	s.Equal(hcl.Diagnostics{{
		Severity: hcl.DiagError,
		Summary:  "Failed to parse arguments",
		Detail:   "invalid format: invalid_type",
	}}, diags)
}

func (s *FrontMatterGeneratorTestSuite) TestContentAndQueryResultMissing() {
	args := cty.ObjectVal(map[string]cty.Value{
		"content": cty.NullVal(cty.DynamicPseudoType),
		"format":  cty.NullVal(cty.String),
	})
	ctx := context.Background()
	document := plugin.ContentSection{}
	content, diags := s.schema.ProvideContent(ctx, "frontmatter", &plugin.ProvideContentParams{
		Args: args,
		DataContext: plugin.MapData{
			"document": plugin.MapData{
				"content": document.AsData(),
			},
		},
	})
	s.Nil(content)
	s.Equal(hcl.Diagnostics{{
		Severity: hcl.DiagError,
		Summary:  "Failed to parse arguments",
		Detail:   "query_result and content are nil",
	}}, diags)
}

func (s *FrontMatterGeneratorTestSuite) TestInvalidQueryResult() {
	args := cty.ObjectVal(map[string]cty.Value{
		"content": cty.NullVal(cty.DynamicPseudoType),
		"format":  cty.NullVal(cty.String),
	})
	ctx := context.Background()
	document := plugin.ContentSection{}
	content, diags := s.schema.ProvideContent(ctx, "frontmatter", &plugin.ProvideContentParams{
		Args: args,
		DataContext: plugin.MapData{
			"query_result": plugin.StringData("invalid_type"),
			"document": plugin.MapData{
				"content": document.AsData(),
			},
		},
	})
	s.Nil(content)
	s.Equal(hcl.Diagnostics{{
		Severity: hcl.DiagError,
		Summary:  "Failed to parse arguments",
		Detail:   "invalid query result: plugin.StringData",
	}}, diags)
}

func (s *FrontMatterGeneratorTestSuite) TestContentAndDataContextNil() {
	args := cty.ObjectVal(map[string]cty.Value{
		"content": cty.NullVal(cty.DynamicPseudoType),
		"format":  cty.NullVal(cty.String),
	})
	ctx := context.Background()
	document := plugin.ContentSection{}
	content, diags := s.schema.ProvideContent(ctx, "frontmatter", &plugin.ProvideContentParams{
		Args: args,
		DataContext: plugin.MapData{
			"document": plugin.MapData{
				"content": document.AsData(),
			},
		},
	})
	s.Nil(content)
	s.Equal(hcl.Diagnostics{{
		Severity: hcl.DiagError,
		Summary:  "Failed to parse arguments",
		Detail:   "query_result and content are nil",
	}}, diags)
}

func (s *FrontMatterGeneratorTestSuite) TestWithContent() {
	args := cty.ObjectVal(map[string]cty.Value{
		"content": cty.ObjectVal(map[string]cty.Value{
			"baz": cty.NumberIntVal(1),
			"foo": cty.StringVal("bar"),
			"quux": cty.ObjectVal(map[string]cty.Value{
				"corge":  cty.StringVal("grault"),
				"garply": cty.BoolVal(false),
			}),
			"qux": cty.BoolVal(true),
			"waldo": cty.ListVal([]cty.Value{
				cty.StringVal("fred"),
				cty.StringVal("plugh"),
			}),
		}),
		"format": cty.NullVal(cty.String),
	})
	ctx := context.Background()
	document := plugin.ContentSection{}
	result, diags := s.schema.ProvideContent(ctx, "frontmatter", &plugin.ProvideContentParams{
		Args: args,
		DataContext: plugin.MapData{
			"document": plugin.MapData{
				"content": document.AsData(),
			},
		},
	})
	s.Require().Nil(diags)
	s.Equal("---\n"+
		"baz: 1\n"+
		"foo: bar\n"+
		"quux:\n"+
		"    corge: grault\n"+
		"    garply: false\n"+
		"qux: true\n"+
		"waldo:\n"+
		"    - fred\n"+
		"    - plugh\n"+
		"---\n", result.Content.Print())
}

func (s *FrontMatterGeneratorTestSuite) TestWithQueryResult() {
	args := cty.ObjectVal(map[string]cty.Value{
		"content": cty.NullVal(cty.DynamicPseudoType),
		"format":  cty.NullVal(cty.String),
	})
	ctx := context.Background()
	document := plugin.ContentSection{}
	result, diags := s.schema.ProvideContent(ctx, "frontmatter", &plugin.ProvideContentParams{
		Args: args,
		DataContext: plugin.MapData{
			"query_result": plugin.MapData{
				"baz": plugin.NumberData(1),
				"foo": plugin.StringData("bar"),
				"quux": plugin.MapData{
					"corge":  plugin.StringData("grault"),
					"garply": plugin.BoolData(false),
				},
				"qux": plugin.BoolData(true),
				"waldo": plugin.ListData{
					plugin.StringData("fred"),
					plugin.StringData("plugh"),
				},
			},
			"document": plugin.MapData{
				"content": document.AsData(),
			},
		},
	})
	s.Equal("---\n"+
		"baz: 1\n"+
		"foo: bar\n"+
		"quux:\n"+
		"    corge: grault\n"+
		"    garply: false\n"+
		"qux: true\n"+
		"waldo:\n"+
		"    - fred\n"+
		"    - plugh\n"+
		"---\n", result.Content.Print())
	s.Nil(diags)
}

func (s *FrontMatterGeneratorTestSuite) TestFormatYaml() {
	args := cty.ObjectVal(map[string]cty.Value{
		"content": cty.NullVal(cty.DynamicPseudoType),
		"format":  cty.StringVal("yaml"),
	})
	ctx := context.Background()
	document := plugin.ContentSection{}
	result, diags := s.schema.ProvideContent(ctx, "frontmatter", &plugin.ProvideContentParams{
		Args: args,
		DataContext: plugin.MapData{
			"query_result": plugin.MapData{
				"baz": plugin.NumberData(1),
				"foo": plugin.StringData("bar"),
				"quux": plugin.MapData{
					"corge":  plugin.StringData("grault"),
					"garply": plugin.BoolData(false),
				},
				"qux": plugin.BoolData(true),
				"waldo": plugin.ListData{
					plugin.StringData("fred"),
					plugin.StringData("plugh"),
				},
			},
			"document": plugin.MapData{
				"content": document.AsData(),
			},
		},
	})
	s.Equal("---\n"+
		"baz: 1\n"+
		"foo: bar\n"+
		"quux:\n"+
		"    corge: grault\n"+
		"    garply: false\n"+
		"qux: true\n"+
		"waldo:\n"+
		"    - fred\n"+
		"    - plugh\n"+
		"---\n", result.Content.Print())
	s.Nil(diags)
}

func (s *FrontMatterGeneratorTestSuite) TestFormatTOML() {
	args := cty.ObjectVal(map[string]cty.Value{
		"content": cty.NullVal(cty.DynamicPseudoType),
		"format":  cty.StringVal("toml"),
	})
	ctx := context.Background()
	document := plugin.ContentSection{}
	result, diags := s.schema.ProvideContent(ctx, "frontmatter", &plugin.ProvideContentParams{
		Args: args,
		DataContext: plugin.MapData{
			"query_result": plugin.MapData{
				"baz": plugin.NumberData(1),
				"foo": plugin.StringData("bar"),
				"quux": plugin.MapData{
					"corge":  plugin.StringData("grault"),
					"garply": plugin.BoolData(false),
				},
				"qux": plugin.BoolData(true),
				"waldo": plugin.ListData{
					plugin.StringData("fred"),
					plugin.StringData("plugh"),
				},
			},
			"document": plugin.MapData{
				"content": document.AsData(),
			},
		},
	})
	s.Equal("+++\n"+
		"baz = 1.0\n"+
		"foo = 'bar'\n"+
		"qux = true\n"+
		"waldo = ['fred', 'plugh']\n\n"+
		"[quux]\n"+
		"corge = 'grault'\n"+
		"garply = false\n"+
		"+++\n", result.Content.Print())
	s.Nil(diags)
}

func (s *FrontMatterGeneratorTestSuite) TestFormatJSON() {
	args := cty.ObjectVal(map[string]cty.Value{
		"content": cty.NullVal(cty.DynamicPseudoType),
		"format":  cty.StringVal("json"),
	})
	ctx := context.Background()
	document := plugin.ContentSection{}
	result, diags := s.schema.ProvideContent(ctx, "frontmatter", &plugin.ProvideContentParams{
		Args: args,
		DataContext: plugin.MapData{
			"query_result": plugin.MapData{
				"baz": plugin.NumberData(1),
				"foo": plugin.StringData("bar"),
				"quux": plugin.MapData{
					"corge":  plugin.StringData("grault"),
					"garply": plugin.BoolData(false),
				},
				"qux": plugin.BoolData(true),
				"waldo": plugin.ListData{
					plugin.StringData("fred"),
					plugin.StringData("plugh"),
				},
			},
			"document": plugin.MapData{
				"content": document.AsData(),
			},
		},
	})
	s.Equal("{\n"+
		"  \"baz\": 1,\n"+
		"  \"foo\": \"bar\",\n"+
		"  \"quux\": {\n"+
		"    \"corge\": \"grault\",\n"+
		"    \"garply\": false\n"+
		"  },\n"+
		"  \"qux\": true,\n"+
		"  \"waldo\": [\n"+
		"    \"fred\",\n"+
		"    \"plugh\"\n"+
		"  ]\n"+
		"}\n", result.Content.Print())
	s.Nil(diags)
}
