package builtin

import (
	"context"
	"testing"

	"github.com/hashicorp/hcl/v2"
	"github.com/stretchr/testify/suite"
	"github.com/zclconf/go-cty/cty"

	"github.com/blackstork-io/fabric/internal/testtools"
	"github.com/blackstork-io/fabric/plugin"
)

type TextTestSuite struct {
	suite.Suite
	schema *plugin.ContentProvider
}

func TestTextSuite(t *testing.T) {
	suite.Run(t, &TextTestSuite{})
}

func (s *TextTestSuite) SetupSuite() {
	s.schema = makeTextContentProvider()
}

func (s *TextTestSuite) TestSchema() {
	s.Nil(s.schema.Config)
	s.NotNil(s.schema.Args)
	s.NotNil(s.schema.ContentFunc)
}

func (s *TextTestSuite) TestMissingText() {
	val := cty.ObjectVal(map[string]cty.Value{
		"value": cty.NullVal(cty.String),
	})
	testtools.ReencodeCTY(s.T(), s.schema.Args, val, [][]testtools.Assert{{
		testtools.IsError,
		testtools.SummaryContains("Non-null value is required"),
	}})
}

func (s *TextTestSuite) TestBasic() {
	val := cty.ObjectVal(map[string]cty.Value{
		"value": cty.StringVal("Hello {{.name}}!"),
	})
	args := testtools.ReencodeCTY(s.T(), s.schema.Args, val, nil)

	ctx := context.Background()
	result, diags := s.schema.ContentFunc(ctx, &plugin.ProvideContentParams{
		Args: args,
		DataContext: plugin.MapData{
			"name": plugin.StringData("World"),
		},
	})
	s.Empty(diags)
	s.Equal("Hello World!", result.Content.Print())
}

func (s *TextTestSuite) TestNoTemplate() {
	val := cty.ObjectVal(map[string]cty.Value{
		"value": cty.StringVal("Hello World!"),
	})
	args := testtools.ReencodeCTY(s.T(), s.schema.Args, val, nil)
	ctx := context.Background()
	result, diags := s.schema.ContentFunc(ctx, &plugin.ProvideContentParams{
		Args:        args,
		DataContext: nil,
	})
	s.Empty(diags)
	s.Equal("Hello World!", result.Content.Print())
}

func (s *TextTestSuite) TestCallInvalidTemplate() {
	val := cty.ObjectVal(map[string]cty.Value{
		"value": cty.StringVal("Hello {{.name}!"),
	})
	args := testtools.ReencodeCTY(s.T(), s.schema.Args, val, nil)

	ctx := context.Background()
	result, diags := s.schema.ContentFunc(ctx, &plugin.ProvideContentParams{
		Args: args,
		DataContext: plugin.MapData{
			"name": plugin.StringData("World"),
		},
	})
	s.Nil(result)
	s.Equal(hcl.Diagnostics{{
		Severity: hcl.DiagError,
		Summary:  "Failed to render text",
		Detail:   "failed to parse text template: template: text:1: bad character U+007D '}'",
	}}, diags)
}

func (s *TextTestSuite) TestSprigTemplate() {
	ctx := context.Background()
	result, diags := s.schema.ContentFunc(ctx, &plugin.ProvideContentParams{
		Args: cty.ObjectVal(map[string]cty.Value{
			"value": cty.StringVal("Hello {{.name | upper}}!"),
		}),
		DataContext: plugin.MapData{
			"name": plugin.StringData("World"),
		},
	})
	s.Empty(diags)
	s.Equal("Hello WORLD!", result.Content.Print())
}
