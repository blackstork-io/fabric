package image

import (
	"testing"

	"github.com/blackstork-io/fabric/plugininterface/v1"
	"github.com/hashicorp/hcl/v2"
	"github.com/stretchr/testify/suite"
	"github.com/zclconf/go-cty/cty"
)

type PluginTestSuite struct {
	suite.Suite
	plugin plugininterface.PluginRPC
}

func TestPluginSuite(t *testing.T) {
	suite.Run(t, &PluginTestSuite{})
}

func (s *PluginTestSuite) SetupSuite() {
	s.plugin = Plugin{}
}

func (s *PluginTestSuite) TestGetPlugins() {
	plugins := s.plugin.GetPlugins()
	s.Require().Len(plugins, 1, "expected 1 plugin")
	got := plugins[0]
	s.Equal("image", got.Name)
	s.Equal("content", got.Kind)
	s.Equal("blackstork", got.Namespace)
	s.Equal(Version.String(), got.Version.Cast().String())
	s.Nil(got.ConfigSpec)
	s.NotNil(got.InvocationSpec)
}

func (s *PluginTestSuite) TestMissingImageSource() {
	args := plugininterface.Args{
		Kind: "content",
		Name: "image",
		Args: cty.ObjectVal(map[string]cty.Value{
			"src": cty.NullVal(cty.String),
			"alt": cty.NullVal(cty.String),
		}),
	}
	result := s.plugin.Call(args)
	expected := plugininterface.Result{
		Diags: hcl.Diagnostics{{
			Severity: hcl.DiagError,
			Summary:  "Failed to parse arguments",
			Detail:   "src is required",
		}},
	}
	s.Equal(expected, result)
}

func (s *PluginTestSuite) TestCallImageSourceEmpty() {
	args := plugininterface.Args{
		Kind: "content",
		Name: "image",
		Args: cty.ObjectVal(map[string]cty.Value{
			"src": cty.StringVal(""),
			"alt": cty.NullVal(cty.String),
		}),
	}
	result := s.plugin.Call(args)
	expected := plugininterface.Result{
		Diags: hcl.Diagnostics{{
			Severity: hcl.DiagError,
			Summary:  "Failed to parse arguments",
			Detail:   "src is required",
		}},
	}
	s.Equal(expected, result)
}

func (s *PluginTestSuite) TestCallImageSourceValid() {
	args := plugininterface.Args{
		Kind: "content",
		Name: "image",
		Args: cty.ObjectVal(map[string]cty.Value{
			"src": cty.StringVal("https://example.com/image.png"),
			"alt": cty.NullVal(cty.String),
		}),
	}
	result := s.plugin.Call(args)
	expected := plugininterface.Result{
		Result: "![](https://example.com/image.png)",
	}
	s.Equal(expected, result)
}

func (s *PluginTestSuite) TestCallImageSourceValidWithAlt() {
	args := plugininterface.Args{
		Kind: "content",
		Name: "image",
		Args: cty.ObjectVal(map[string]cty.Value{
			"src": cty.StringVal("https://example.com/image.png"),
			"alt": cty.StringVal("alt text"),
		}),
	}
	result := s.plugin.Call(args)
	expected := plugininterface.Result{
		Result: "![alt text](https://example.com/image.png)",
	}
	s.Equal(expected, result)
}
