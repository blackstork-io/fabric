package builtin

import (
	"context"
	"testing"

	"github.com/stretchr/testify/suite"

	"github.com/blackstork-io/fabric/internal/testtools"
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
	testtools.Decode(s.T(), s.schema.Args, ``, [][]testtools.Assert{{testtools.DetailContains(`The argument "value" is required`)}})
	return
}

func (s *BlockQuoteTestSuite) TestCallBlockquote() {
	ctx := context.Background()
	args := testtools.Decode(s.T(), s.schema.Args, `
		value = "Hello {{.name}}!"
	`, nil)
	content, diags := s.schema.ContentFunc(ctx, &plugin.ProvideContentParams{
		Args: args,
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
	args := testtools.Decode(s.T(), s.schema.Args, `
		value = "Hello\n{{.name}}\nfor you!"
	`, nil)
	content, diags := s.schema.ContentFunc(ctx, &plugin.ProvideContentParams{
		Args: args,
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
	args := testtools.Decode(s.T(), s.schema.Args, `
		value = "Hello\n{{.name}}\n\nfor you!"
	`, nil)
	content, diags := s.schema.ContentFunc(ctx, &plugin.ProvideContentParams{
		Args: args,
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
