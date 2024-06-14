package builtin

import (
	"bytes"
	"context"
	"encoding/json"
	"testing"

	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/hclwrite"
	"github.com/pelletier/go-toml/v2"
	"github.com/stretchr/testify/suite"
	"github.com/zclconf/go-cty/cty"
	"github.com/zclconf/go-cty/cty/gocty"
	"gopkg.in/yaml.v3"

	"github.com/blackstork-io/fabric/pkg/diagnostics"
	"github.com/blackstork-io/fabric/pkg/diagnostics/diagtest"
	"github.com/blackstork-io/fabric/pkg/utils"
	"github.com/blackstork-io/fabric/plugin"
	"github.com/blackstork-io/fabric/plugin/plugintest"
	"github.com/blackstork-io/fabric/print/mdprint"
)

type testStructInner struct {
	Corge  string `cty:"Corge" yaml:"Corge" toml:"Corge" json:"Corge"`
	Garply bool   `cty:"Garply" yaml:"Garply" toml:"Garply" json:"Garply"`
}

type testStruct struct {
	Baz   float32         `cty:"Baz" yaml:"Baz" toml:"Baz" json:"Baz"`
	Foo   string          `cty:"Foo" yaml:"Foo" toml:"Foo" json:"Foo"`
	Quux  testStructInner `cty:"Quux" yaml:"Quux" toml:"Quux" json:"Quux"`
	Qux   bool            `cty:"Qux" yaml:"Qux" toml:"Qux" json:"Qux"`
	Waldo []string        `cty:"Waldo" yaml:"Waldo" toml:"Waldo" json:"Waldo"`
}

var testVal = testStruct{
	Baz: 1,
	Foo: "bar",
	Quux: testStructInner{
		Corge:  "grault",
		Garply: false,
	},
	Qux:   true,
	Waldo: []string{"fred", "plugh"},
}

var testValCty = utils.Must(gocty.ToCtyValue(
	&testVal,
	utils.Must(gocty.ImpliedType(&testVal)),
))

type FrontMatterGeneratorTestSuite struct {
	suite.Suite
	schema *plugin.Schema
}

func TestFrontMatterGeneratorSuite(t *testing.T) {
	suite.Run(t, &FrontMatterGeneratorTestSuite{})
}

func (s *FrontMatterGeneratorTestSuite) parseFrontmatter(contentStr string) (format string) {
	require := s.Require()

	content := bytes.TrimSpace([]byte(contentStr))
	var result testStruct
	var err error
	switch {
	case bytes.HasPrefix(content, []byte("---\n")) && bytes.HasSuffix(content, []byte("\n---")):
		content = bytes.Trim(content, "-\n")
		err = yaml.Unmarshal(content, &result)
		format = "yaml"
	case bytes.HasPrefix(content, []byte("+++\n")) && bytes.HasSuffix(content, []byte("\n+++")):
		content = bytes.Trim(content, "+\n")
		err = toml.Unmarshal(content, &result)
		format = "toml"
	default:
		err = json.Unmarshal(content, &result)
		format = "json"
	}
	require.NoError(err)
	require.Equal(testVal, result)
	return
}

func (s *FrontMatterGeneratorTestSuite) SetupSuite() {
	s.schema = Plugin("", nil, nil)
}

func (s *FrontMatterGeneratorTestSuite) TestSchema() {
	provider := s.schema.ContentProviders["frontmatter"]
	s.NotNil(provider)
	s.Nil(provider.Config)
	s.NotNil(provider.Args)
	s.NotNil(provider.ContentFunc)
}

func (s *FrontMatterGeneratorTestSuite) TestInvalidFormat() {
	val := cty.ObjectVal(map[string]cty.Value{
		"format": cty.StringVal("invalid_type"),
		"content": cty.ObjectVal(map[string]cty.Value{
			"foo": cty.StringVal("bar"),
		}),
	})
	plugintest.ReencodeCTY(s.T(), s.schema.ContentProviders["frontmatter"].Args, val, diagtest.Asserts{{
		diagtest.IsError,
		diagtest.SummaryContains("Attribute", "not one of"),
	}})
}

func (s *FrontMatterGeneratorTestSuite) TestContentAndQueryResultMissing() {
	val := cty.ObjectVal(map[string]cty.Value{
		"content": cty.NullVal(cty.DynamicPseudoType),
		"format":  cty.NullVal(cty.String),
	})
	args := plugintest.ReencodeCTY(s.T(), s.schema.ContentProviders["frontmatter"].Args, val, nil)

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
	s.Equal(diagnostics.Diag{{
		Severity: hcl.DiagError,
		Summary:  "Failed to parse arguments",
		Detail:   "Content is nil",
	}}, diags)
}

func (s *FrontMatterGeneratorTestSuite) TestInvalidQueryResult() {
	val := `
		format = null
		content = "invalid_type"
	`
	document := plugin.ContentSection{}
	dataCtx := plugin.MapData{
		"document": plugin.MapData{
			"content": document.AsData(),
		},
	}

	args := plugintest.DecodeAndAssert(s.T(), s.schema.ContentProviders["frontmatter"].Args, val, dataCtx, diagtest.Asserts{})

	ctx := context.Background()
	content, diags := s.schema.ProvideContent(ctx, "frontmatter", &plugin.ProvideContentParams{
		Args:        args,
		DataContext: dataCtx,
	})
	s.Nil(content)

	diagtest.Asserts{{
		diagtest.IsError,
		diagtest.DetailContains("invalid", "plugin.StringData"),
	}}.AssertMatch(s.T(), diags, nil)
}

func (s *FrontMatterGeneratorTestSuite) TestContentAndDataContextNil() {
	val := cty.ObjectVal(map[string]cty.Value{
		"content": cty.NullVal(cty.DynamicPseudoType),
		"format":  cty.NullVal(cty.String),
	})
	args := plugintest.ReencodeCTY(s.T(), s.schema.ContentProviders["frontmatter"].Args, val, nil)

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
	s.Equal(diagnostics.Diag{{
		Severity: hcl.DiagError,
		Summary:  "Failed to parse arguments",
		Detail:   "Content is nil",
	}}, diags)
}

func (s *FrontMatterGeneratorTestSuite) TestWithContent() {
	f := hclwrite.NewEmptyFile()
	f.Body().SetAttributeValue("content", testValCty)
	body := string(f.Bytes())

	document := plugin.ContentSection{}
	dataCtx := plugin.MapData{
		"document": plugin.MapData{
			"content": document.AsData(),
		},
	}

	args := plugintest.DecodeAndAssert(
		s.T(), s.schema.ContentProviders["frontmatter"].Args,
		body, dataCtx, diagtest.Asserts{},
	)

	ctx := context.Background()
	result, diags := s.schema.ProvideContent(ctx, "frontmatter", &plugin.ProvideContentParams{
		Args: args,
		DataContext: plugin.MapData{
			"document": plugin.MapData{
				"content": document.AsData(),
			},
		},
	})
	s.Require().Nil(diags)
	format := s.parseFrontmatter(mdprint.PrintString(result.Content))
	s.Equal("yaml", format)
}

func (s *FrontMatterGeneratorTestSuite) TestFormatYaml() {
	f := hclwrite.NewEmptyFile()
	hclBody := f.Body()
	hclBody.SetAttributeValue("content", testValCty)
	hclBody.SetAttributeValue("format", cty.StringVal("yaml"))
	body := string(f.Bytes())

	document := plugin.ContentSection{}
	dataCtx := plugin.MapData{
		"document": plugin.MapData{
			"content": document.AsData(),
		},
	}

	args := plugintest.DecodeAndAssert(
		s.T(), s.schema.ContentProviders["frontmatter"].Args,
		body, dataCtx, diagtest.Asserts{},
	)

	ctx := context.Background()
	result, diags := s.schema.ProvideContent(ctx, "frontmatter", &plugin.ProvideContentParams{
		Args: args,
		DataContext: plugin.MapData{
			"document": plugin.MapData{
				"content": document.AsData(),
			},
		},
	})
	s.Require().Nil(diags)
	format := s.parseFrontmatter(mdprint.PrintString(result.Content))
	s.Equal("yaml", format)
}

func (s *FrontMatterGeneratorTestSuite) TestFormatTOML() {
	f := hclwrite.NewEmptyFile()
	hclBody := f.Body()
	hclBody.SetAttributeValue("content", testValCty)
	hclBody.SetAttributeValue("format", cty.StringVal("toml"))
	body := string(f.Bytes())

	document := plugin.ContentSection{}
	dataCtx := plugin.MapData{
		"document": plugin.MapData{
			"content": document.AsData(),
		},
	}

	args := plugintest.DecodeAndAssert(
		s.T(), s.schema.ContentProviders["frontmatter"].Args,
		body, dataCtx, diagtest.Asserts{},
	)

	ctx := context.Background()
	result, diags := s.schema.ProvideContent(ctx, "frontmatter", &plugin.ProvideContentParams{
		Args: args,
		DataContext: plugin.MapData{
			"document": plugin.MapData{
				"content": document.AsData(),
			},
		},
	})
	s.Require().Nil(diags)
	format := s.parseFrontmatter(mdprint.PrintString(result.Content))
	s.Equal("toml", format)
}

func (s *FrontMatterGeneratorTestSuite) TestFormatJSON() {
	f := hclwrite.NewEmptyFile()
	hclBody := f.Body()
	hclBody.SetAttributeValue("content", testValCty)
	hclBody.SetAttributeValue("format", cty.StringVal("json"))
	body := string(f.Bytes())

	document := plugin.ContentSection{}
	dataCtx := plugin.MapData{
		"document": plugin.MapData{
			"content": document.AsData(),
		},
	}

	args := plugintest.DecodeAndAssert(
		s.T(), s.schema.ContentProviders["frontmatter"].Args,
		body, dataCtx, diagtest.Asserts{},
	)

	ctx := context.Background()
	result, diags := s.schema.ProvideContent(ctx, "frontmatter", &plugin.ProvideContentParams{
		Args: args,
		DataContext: plugin.MapData{
			"document": plugin.MapData{
				"content": document.AsData(),
			},
		},
	})
	s.Require().Nil(diags)
	format := s.parseFrontmatter(mdprint.PrintString(result.Content))
	s.Equal("json", format)
}
