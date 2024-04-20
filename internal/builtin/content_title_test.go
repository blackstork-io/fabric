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

type TitleTestSuite struct {
	suite.Suite
	schema *plugin.ContentProvider
}

func TestTitleSuite(t *testing.T) {
	suite.Run(t, &TitleTestSuite{})
}

func (s *TitleTestSuite) SetupSuite() {
	s.schema = makeTitleContentProvider()
}

func (s *TitleTestSuite) TestSchema() {
	s.Nil(s.schema.Config)
	s.NotNil(s.schema.Args)
	s.NotNil(s.schema.ContentFunc)
}

func (s *TitleTestSuite) TestMissingValue() {
	val := cty.ObjectVal(map[string]cty.Value{
		"value":         cty.NullVal(cty.String),
		"absolute_size": cty.NullVal(cty.Number),
		"relative_size": cty.NullVal(cty.Number),
	})
	testtools.ReencodeCTY(s.T(), s.schema.Args, val, [][]testtools.Assert{{
		testtools.IsError,
		testtools.SummaryContains("Non-null value is required"),
	}})
}

func (s *TitleTestSuite) TestTDefault() {
	val := cty.ObjectVal(map[string]cty.Value{
		"value":         cty.StringVal("Hello {{.name}}!"),
		"absolute_size": cty.NullVal(cty.Number),
		"relative_size": cty.NullVal(cty.Number),
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
	s.Equal("## Hello World!", result.Content.Print())
}

func (s *TitleTestSuite) TestWithTextMultiline() {
	val := cty.ObjectVal(map[string]cty.Value{
		"value":         cty.StringVal("Hello\n{{.name}}\nfor you!"),
		"absolute_size": cty.NullVal(cty.Number),
		"relative_size": cty.NullVal(cty.Number),
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
	s.Equal("## Hello World for you!", result.Content.Print())
}

func (s *TitleTestSuite) TestWithSize() {
	val := cty.ObjectVal(map[string]cty.Value{
		"value":         cty.StringVal("Hello {{.name}}!"),
		"absolute_size": cty.NumberIntVal(2),
		"relative_size": cty.NullVal(cty.Number),
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
	s.Equal("### Hello World!", result.Content.Print())
}

func (s *TitleTestSuite) TestWithSizeTooBig() {
	val := cty.ObjectVal(map[string]cty.Value{
		"value":         cty.StringVal("Hello {{.name}}!"),
		"absolute_size": cty.NumberIntVal(7),
		"relative_size": cty.NullVal(cty.Number),
	})
	args := testtools.ReencodeCTY(s.T(), s.schema.Args, val, nil)
	ctx := context.Background()
	result, diags := s.schema.ContentFunc(ctx, &plugin.ProvideContentParams{
		Args: args,
		DataContext: plugin.MapData{
			"name": plugin.StringData("World"),
		},
	})
	s.Equal(hcl.Diagnostics{{
		Severity: hcl.DiagError,
		Summary:  "Failed to parse arguments",
		Detail:   "absolute_size must be between 0 and 5",
	}}, diags)
	s.Nil(result)
}
