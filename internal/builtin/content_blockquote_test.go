package builtin

import (
	"context"
	"testing"

	"github.com/hashicorp/hcl/v2"
	"github.com/stretchr/testify/suite"
	"github.com/zclconf/go-cty/cty"

	"github.com/blackstork-io/fabric/plugin"
)

type BlockQuoteTestSuite struct {
	suite.Suite
	schema *plugin.ContentProvider
}

func TestBlockQuoteSuite(t *testing.T) {
	suite.Run(t, &BlockQuoteTestSuite{})
}

func (s *BlockQuoteTestSuite) SetupSuite() {
	s.schema = makeBlockQuoteContentProvider()
}

func (s *BlockQuoteTestSuite) TestSchema() {
	s.Nil(s.schema.Config)
	s.NotNil(s.schema.Args)
	s.NotNil(s.schema.ContentFunc)
}

func (s *BlockQuoteTestSuite) TestMissingText() {
	ctx := context.Background()
	content, diags := s.schema.ContentFunc(ctx, &plugin.ProvideContentParams{
		Args: cty.ObjectVal(map[string]cty.Value{
			"value": cty.NullVal(cty.String),
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

func (s *BlockQuoteTestSuite) TestCallBlockquote() {
	ctx := context.Background()
	content, diags := s.schema.ContentFunc(ctx, &plugin.ProvideContentParams{
		Args: cty.ObjectVal(map[string]cty.Value{
			"value": cty.StringVal(`Hello {{.name}}!`),
		}),
		DataContext: plugin.MapData{
			"name": plugin.StringData("World"),
		},
	})
	s.Empty(diags)
	s.Equal(&plugin.ContentResult{
		Content: &plugin.ContentElement{
			Markdown: "> Hello World!",
		},
	}, content)
}

func (s *BlockQuoteTestSuite) TestCallBlockquoteMultiline() {
	ctx := context.Background()
	content, diags := s.schema.ContentFunc(ctx, &plugin.ProvideContentParams{
		Args: cty.ObjectVal(map[string]cty.Value{
			"value": cty.StringVal("Hello\n{{.name}}\nfor you!"),
		}),
		DataContext: plugin.MapData{
			"name": plugin.StringData("World"),
		},
	})
	s.Empty(diags)
	s.Equal(&plugin.ContentResult{
		Content: &plugin.ContentElement{
			Markdown: "> Hello\n> World\n> for you!",
		},
	}, content)
}

func (s *BlockQuoteTestSuite) TestCallBlockquoteMultilineDoubleNewline() {
	ctx := context.Background()
	content, diags := s.schema.ContentFunc(ctx, &plugin.ProvideContentParams{
		Args: cty.ObjectVal(map[string]cty.Value{
			"value": cty.StringVal("Hello\n{{.name}}\n\nfor you!"),
		}),
		DataContext: plugin.MapData{
			"name": plugin.StringData("World"),
		},
	})
	s.Empty(diags)
	s.Equal(&plugin.ContentResult{
		Content: &plugin.ContentElement{
			Markdown: "> Hello\n> World\n> \n> for you!",
		},
	}, content)
}
