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
	val := cty.ObjectVal(map[string]cty.Value{
		"src": cty.NullVal(cty.String),
		"alt": cty.NullVal(cty.String),
	})
	testtools.ReencodeCTY(s.T(), s.schema.Args, val, [][]testtools.Assert{{
		testtools.IsError,
		testtools.SummaryContains("Argument must be non-null"),
	}})
}

func (s *ImageGeneratorTestSuite) TestCallImageSourceEmpty() {
	val := cty.ObjectVal(map[string]cty.Value{
		"src": cty.StringVal(""),
		"alt": cty.NullVal(cty.String),
	})
	args := testtools.ReencodeCTY(s.T(), s.schema.Args, val, nil)
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
	val := cty.ObjectVal(map[string]cty.Value{
		"src": cty.StringVal("https://example.com/image.png"),
		"alt": cty.NullVal(cty.String),
	})
	args := testtools.ReencodeCTY(s.T(), s.schema.Args, val, nil)

	ctx := context.Background()
	result, diags := s.schema.ContentFunc(ctx, &plugin.ProvideContentParams{
		Args: args,
	})
	s.Equal("![](https://example.com/image.png)", result.Content.Print())
	s.Empty(diags)
}

func (s *ImageGeneratorTestSuite) TestCallImageSourceValidWithAlt() {
	val := cty.ObjectVal(map[string]cty.Value{
		"src": cty.StringVal("https://example.com/image.png"),
		"alt": cty.StringVal("alt text"),
	})
	args := testtools.ReencodeCTY(s.T(), s.schema.Args, val, nil)

	ctx := context.Background()
	result, diags := s.schema.ContentFunc(ctx, &plugin.ProvideContentParams{
		Args: args,
	})
	s.Equal("![alt text](https://example.com/image.png)", result.Content.Print())
	s.Empty(diags)
}
