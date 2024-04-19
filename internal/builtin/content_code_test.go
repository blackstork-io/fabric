package builtin

import (
	"context"
	"testing"

	"github.com/stretchr/testify/suite"
	"github.com/zclconf/go-cty/cty"

	"github.com/blackstork-io/fabric/internal/testtools"
	"github.com/blackstork-io/fabric/plugin"
)

type CodeTestSuite struct {
	suite.Suite
	schema *plugin.ContentProvider
}

func TestCodeSuite(t *testing.T) {
	suite.Run(t, &CodeTestSuite{})
}

func (s *CodeTestSuite) SetupSuite() {
	s.schema = makeCodeContentProvider()
}

func (s *CodeTestSuite) TestSchema() {
	s.Nil(s.schema.Config)
	s.NotNil(s.schema.Args)
	s.NotNil(s.schema.ContentFunc)
}

func (s *CodeTestSuite) TestMissingValue() {
	val := cty.ObjectVal(map[string]cty.Value{
		"value":    cty.NullVal(cty.String),
		"language": cty.NullVal(cty.String),
	})
	testtools.ReencodeCTY(s.T(), s.schema.Args, val, [][]testtools.Assert{{
		testtools.SummaryContains("Non-null value is required"),
	}})
}

func (s *CodeTestSuite) TestCallCodeDefault() {
	ctx := context.Background()
	val := cty.ObjectVal(map[string]cty.Value{
		"value":    cty.StringVal(`Hello {{.name}}!`),
		"language": cty.NullVal(cty.String),
	})
	args := testtools.ReencodeCTY(s.T(), s.schema.Args, val, nil)

	content, diags := s.schema.ContentFunc(ctx, &plugin.ProvideContentParams{
		Args: args,
		DataContext: plugin.MapData{
			"name": plugin.StringData("World"),
		},
	})
	s.Empty(diags)
	s.Equal(&plugin.ContentResult{
		Content: &plugin.ContentElement{
			Markdown: "```\nHello World!\n```",
		},
	}, content)
}

func (s *CodeTestSuite) TestCallCodeWithLanguage() {
	ctx := context.Background()
	val := cty.ObjectVal(map[string]cty.Value{
		"value":    cty.StringVal(`{"hello": "{{.name}}"}`),
		"language": cty.StringVal("json"),
	})
	args := testtools.ReencodeCTY(s.T(), s.schema.Args, val, nil)

	content, diags := s.schema.ContentFunc(ctx, &plugin.ProvideContentParams{
		Args: args,
		DataContext: plugin.MapData{
			"name": plugin.StringData("world"),
		},
	})
	s.Empty(diags)
	s.Equal(&plugin.ContentResult{
		Content: &plugin.ContentElement{
			Markdown: "```json\n{\"hello\": \"world\"}\n```",
		},
	}, content)
}
