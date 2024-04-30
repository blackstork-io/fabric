package builtin

import (
	"context"
	"testing"

	"github.com/hashicorp/hcl/v2"
	"github.com/stretchr/testify/suite"

	"github.com/blackstork-io/fabric/internal/testtools"
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
	testtools.DecodeAndAssert(s.T(), s.schema.Args, `
		src = null
		alt = null
		`,
		[][]testtools.Assert{{
			testtools.IsError,
			testtools.SummaryContains("Argument must be non-null"),
		}})
}

func (s *ImageGeneratorTestSuite) TestCallImageSourceEmpty() {
	args := testtools.DecodeAndAssert(s.T(), s.schema.Args, `
		src = ""
		alt = null
		`,
		nil)

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
	args := testtools.DecodeAndAssert(s.T(), s.schema.Args, `
		src = "https://example.com/image.png"
		`,
		nil)

	ctx := context.Background()
	result, diags := s.schema.ContentFunc(ctx, &plugin.ProvideContentParams{
		Args: args,
	})
	s.Equal("![](https://example.com/image.png)", result.Content.Print())
	s.Empty(diags)
}

func (s *ImageGeneratorTestSuite) TestCallImageSourceValidWithAlt() {
	args := testtools.DecodeAndAssert(s.T(), s.schema.Args, `
		src = "https://example.com/image.png"
		alt = "alt text"
		`,
		nil)

	ctx := context.Background()
	result, diags := s.schema.ContentFunc(ctx, &plugin.ProvideContentParams{
		Args: args,
	})
	s.Equal("![alt text](https://example.com/image.png)", result.Content.Print())
	s.Empty(diags)
}

func (s *ImageGeneratorTestSuite) TestCallImageSourceTemplateRender() {
	args := testtools.DecodeAndAssert(s.T(), s.schema.Args, `
		src = "./{{ add 1 2 }}.png"
		`,
		nil)

	ctx := context.Background()
	result, diags := s.schema.ContentFunc(ctx, &plugin.ProvideContentParams{
		Args: args,
	})
	s.Equal("![](./3.png)", result.Content.Print())
	s.Empty(diags)
}

func (s *ImageGeneratorTestSuite) TestCallImageAltTemplateRender() {
	args := testtools.DecodeAndAssert(s.T(), s.schema.Args, `
		src = "./{{ add 1 2 }}.png"
		alt = "{{ add 2 3 }} alt text"
		`,
		nil)

	ctx := context.Background()
	result, diags := s.schema.ContentFunc(ctx, &plugin.ProvideContentParams{
		Args: args,
	})
	s.Equal("![5 alt text](./3.png)", result.Content.Print())
	s.Empty(diags)
}
