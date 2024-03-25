package builtin

import (
	"context"
	"testing"

	"github.com/hashicorp/hcl/v2"
	"github.com/stretchr/testify/suite"
	"github.com/zclconf/go-cty/cty"

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
	ctx := context.Background()
	result, diags := s.schema.ContentFunc(ctx, &plugin.ProvideContentParams{
		Args: cty.ObjectVal(map[string]cty.Value{
			"value":         cty.NullVal(cty.String),
			"absolute_size": cty.NullVal(cty.Number),
			"relative_size": cty.NullVal(cty.Number),
		}),
		DataContext: plugin.MapData{
			"name": plugin.StringData("World"),
		},
	})
	s.Nil(result)
	s.Equal(hcl.Diagnostics{{
		Severity: hcl.DiagError,
		Summary:  "Failed to parse arguments",
		Detail:   "value is required",
	}}, diags)
}

func (s *TitleTestSuite) TestTDefault() {
	ctx := context.Background()
	result, diags := s.schema.ContentFunc(ctx, &plugin.ProvideContentParams{
		Args: cty.ObjectVal(map[string]cty.Value{
			"value":         cty.StringVal("Hello {{.name}}!"),
			"absolute_size": cty.NullVal(cty.Number),
			"relative_size": cty.NullVal(cty.Number),
		}),
		DataContext: plugin.MapData{
			"name": plugin.StringData("World"),
		},
	})
	s.Empty(diags)
	s.Equal("## Hello World!", result.Content.Print())
}

func (s *TitleTestSuite) TestWithTextMultiline() {
	ctx := context.Background()
	result, diags := s.schema.ContentFunc(ctx, &plugin.ProvideContentParams{
		Args: cty.ObjectVal(map[string]cty.Value{
			"value":         cty.StringVal("Hello\n{{.name}}\nfor you!"),
			"absolute_size": cty.NullVal(cty.Number),
			"relative_size": cty.NullVal(cty.Number),
		}),
		DataContext: plugin.MapData{
			"name": plugin.StringData("World"),
		},
	})
	s.Empty(diags)
	s.Equal("## Hello World for you!", result.Content.Print())
}

func (s *TitleTestSuite) TestWithSize() {
	ctx := context.Background()
	result, diags := s.schema.ContentFunc(ctx, &plugin.ProvideContentParams{
		Args: cty.ObjectVal(map[string]cty.Value{
			"value":         cty.StringVal("Hello {{.name}}!"),
			"absolute_size": cty.NumberIntVal(2),
			"relative_size": cty.NullVal(cty.Number),
		}),
		DataContext: plugin.MapData{
			"name": plugin.StringData("World"),
		},
	})
	s.Empty(diags)
	s.Equal("### Hello World!", result.Content.Print())
}

func (s *TitleTestSuite) TestWithSizeTooBig() {
	ctx := context.Background()
	result, diags := s.schema.ContentFunc(ctx, &plugin.ProvideContentParams{
		Args: cty.ObjectVal(map[string]cty.Value{
			"value":         cty.StringVal("Hello {{.name}}!"),
			"absolute_size": cty.NumberIntVal(7),
			"relative_size": cty.NullVal(cty.Number),
		}),
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
