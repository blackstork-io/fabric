package builtin

import (
	"context"
	"testing"

	"github.com/hashicorp/hcl/v2"
	"github.com/stretchr/testify/suite"
	"github.com/zclconf/go-cty/cty"

	"github.com/blackstork-io/fabric/plugin"
)

type ImageGeneratorTestSuite struct {
	suite.Suite
	schema *plugin.ContentProvider
}

func TestImageGeneratorTestSuite(t *testing.T) {
	suite.Run(t, &ImageGeneratorTestSuite{})
}

func (s *ImageGeneratorTestSuite) SetupSuite() {
	s.schema = makeImageContentProvider()
}

func (s *ImageGeneratorTestSuite) TestSchema() {
	provider := makeImageContentProvider()
	s.Nil(provider.Config)
	s.NotNil(provider.Args)
	s.NotNil(provider.ContentFunc)
}

func (s *ImageGeneratorTestSuite) TestMissingImageSource() {
	args := cty.ObjectVal(map[string]cty.Value{
		"src": cty.NullVal(cty.String),
		"alt": cty.NullVal(cty.String),
	})
	ctx := context.Background()
	content, diags := s.schema.ContentFunc(ctx, &plugin.ProvideContentParams{
		Args: args,
	})
	s.Nil(content)
	s.Equal(hcl.Diagnostics{{
		Severity: hcl.DiagError,
		Summary:  "Failed to parse arguments",
		Detail:   "src is required",
	}}, diags)
}

func (s *ImageGeneratorTestSuite) TestCallImageSourceEmpty() {
	args := cty.ObjectVal(map[string]cty.Value{
		"src": cty.StringVal(""),
		"alt": cty.NullVal(cty.String),
	})
	ctx := context.Background()
	content, diags := s.schema.ContentFunc(ctx, &plugin.ProvideContentParams{
		Args: args,
	})
	s.Nil(content)
	s.Equal(hcl.Diagnostics{{
		Severity: hcl.DiagError,
		Summary:  "Failed to parse arguments",
		Detail:   "src is required",
	}}, diags)
}

func (s *ImageGeneratorTestSuite) TestCallImageSourceValid() {
	args := cty.ObjectVal(map[string]cty.Value{
		"src": cty.StringVal("https://example.com/image.png"),
		"alt": cty.NullVal(cty.String),
	})
	ctx := context.Background()
	result, diags := s.schema.ContentFunc(ctx, &plugin.ProvideContentParams{
		Args: args,
	})
	s.Equal("![](https://example.com/image.png)", result.Content.Print())
	s.Empty(diags)
}

func (s *ImageGeneratorTestSuite) TestCallImageSourceValidWithAlt() {
	args := cty.ObjectVal(map[string]cty.Value{
		"src": cty.StringVal("https://example.com/image.png"),
		"alt": cty.StringVal("alt text"),
	})
	ctx := context.Background()
	result, diags := s.schema.ContentFunc(ctx, &plugin.ProvideContentParams{
		Args: args,
	})
	s.Equal("![alt text](https://example.com/image.png)", result.Content.Print())
	s.Empty(diags)
}
