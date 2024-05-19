package builtin

import (
	"context"
	"testing"

	"github.com/stretchr/testify/suite"

	"github.com/blackstork-io/fabric/pkg/diagnostics/diagtest"
	"github.com/blackstork-io/fabric/plugin"
	"github.com/blackstork-io/fabric/plugin/plugintest"
	"github.com/blackstork-io/fabric/print/mdprint"
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
	plugintest.DecodeAndAssert(s.T(), s.schema.Args, `
		src = null
		alt = null
		`,
		diagtest.Asserts{{
			diagtest.IsError,
			diagtest.SummaryContains("Attribute must be non-null"),
		}})
}

func (s *ImageGeneratorTestSuite) TestCallImageSourceEmpty() {
	plugintest.DecodeAndAssert(s.T(), s.schema.Args, `
		src = ""
		alt = null
		`,
		diagtest.Asserts{{
			diagtest.IsError,
			diagtest.DetailContains(`The length`, `"src"`, `>= 1`),
		}})
}

func (s *ImageGeneratorTestSuite) TestCallImageSourceValid() {
	args := plugintest.DecodeAndAssert(s.T(), s.schema.Args, `
		src = "https://example.com/image.png"
		`,
		nil)

	ctx := context.Background()
	result, diags := s.schema.ContentFunc(ctx, &plugin.ProvideContentParams{
		Args: args,
	})
	s.Equal("![](https://example.com/image.png)", mdprint.PrintString(result.Content))
	s.Empty(diags)
}

func (s *ImageGeneratorTestSuite) TestCallImageSourceValidWithAlt() {
	args := plugintest.DecodeAndAssert(s.T(), s.schema.Args, `
		src = "https://example.com/image.png"
		alt = "alt text"
		`,
		nil)

	ctx := context.Background()
	result, diags := s.schema.ContentFunc(ctx, &plugin.ProvideContentParams{
		Args: args,
	})
	s.Equal("![alt text](https://example.com/image.png)", mdprint.PrintString(result.Content))
	s.Empty(diags)
}

func (s *ImageGeneratorTestSuite) TestCallImageSourceTemplateRender() {
	args := plugintest.DecodeAndAssert(s.T(), s.schema.Args, `
		src = "./{{ add 1 2 }}.png"
		`,
		nil)

	ctx := context.Background()
	result, diags := s.schema.ContentFunc(ctx, &plugin.ProvideContentParams{
		Args: args,
	})
	s.Equal("![](./3.png)", mdprint.PrintString(result.Content))
	s.Empty(diags)
}

func (s *ImageGeneratorTestSuite) TestCallImageAltTemplateRender() {
	args := plugintest.DecodeAndAssert(s.T(), s.schema.Args, `
		src = "./{{ add 1 2 }}.png"
		alt = "{{ add 2 3 }} alt text"
		`,
		nil)

	ctx := context.Background()
	result, diags := s.schema.ContentFunc(ctx, &plugin.ProvideContentParams{
		Args: args,
	})
	s.Equal("![5 alt text](./3.png)", mdprint.PrintString(result.Content))
	s.Empty(diags)
}
