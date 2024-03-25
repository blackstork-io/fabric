package builtin

import (
	"context"
	"testing"

	"github.com/hashicorp/hcl/v2"
	"github.com/stretchr/testify/suite"
	"github.com/zclconf/go-cty/cty"

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
	ctx := context.Background()
	content, diags := s.schema.ContentFunc(ctx, &plugin.ProvideContentParams{
		Args: cty.ObjectVal(map[string]cty.Value{
			"value":    cty.NullVal(cty.String),
			"language": cty.NullVal(cty.String),
		}),
		DataContext: plugin.MapData{
			"name": plugin.StringData("World"),
		},
	})
	s.Nil(content)
	s.Equal(hcl.Diagnostics{{
		Severity: hcl.DiagError,
		Summary:  "Failed to parse arguments",
		Detail:   "value is required",
	}}, diags)
}

func (s *CodeTestSuite) TestCallCodeDefault() {
	ctx := context.Background()
	content, diags := s.schema.ContentFunc(ctx, &plugin.ProvideContentParams{
		Args: cty.ObjectVal(map[string]cty.Value{
			"value":    cty.StringVal(`Hello {{.name}}!`),
			"language": cty.NullVal(cty.String),
		}),
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
	content, diags := s.schema.ContentFunc(ctx, &plugin.ProvideContentParams{
		Args: cty.ObjectVal(map[string]cty.Value{
			"value":    cty.StringVal(`{"hello": "{{.name}}"}`),
			"language": cty.StringVal("json"),
		}),
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
