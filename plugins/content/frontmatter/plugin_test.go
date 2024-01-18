package frontmatter

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
	s.Equal("frontmatter", got.Name)
	s.Equal("content", got.Kind)
	s.Equal("blackstork", got.Namespace)
	s.Equal(Version.String(), got.Version.Cast().String())
	s.Nil(got.ConfigSpec)
	s.NotNil(got.InvocationSpec)
}

func (s *PluginTestSuite) TestCallInvalidFormat() {
	args := plugininterface.Args{
		Kind: "content",
		Name: "frontmatter",
		Args: cty.ObjectVal(map[string]cty.Value{
			"content": cty.ObjectVal(map[string]cty.Value{
				"foo": cty.StringVal("bar"),
			}),
			"format": cty.StringVal("invalid_type"),
		}),
		Context: map[string]any{},
	}
	result := s.plugin.Call(args)
	expected := plugininterface.Result{
		Diags: hcl.Diagnostics{{
			Severity: hcl.DiagError,
			Summary:  "Failed to parse arguments",
			Detail:   "invalid format: invalid_type",
		}},
	}
	s.Equal(expected, result)
}

func (s *PluginTestSuite) TestCallContentAndQueryResultMissing() {
	args := plugininterface.Args{
		Kind: "content",
		Name: "frontmatter",
		Args: cty.ObjectVal(map[string]cty.Value{
			"content": cty.NullVal(cty.DynamicPseudoType),
			"format":  cty.NullVal(cty.String),
		}),
		Context: map[string]any{},
	}
	result := s.plugin.Call(args)
	expected := plugininterface.Result{
		Diags: hcl.Diagnostics{{
			Severity: hcl.DiagError,
			Summary:  "Failed to parse arguments",
			Detail:   "query_result and content are nil",
		}},
	}
	s.Equal(expected, result)
}

func (s *PluginTestSuite) TestCallInvalidQueryResult() {
	args := plugininterface.Args{
		Kind: "content",
		Name: "frontmatter",
		Args: cty.ObjectVal(map[string]cty.Value{
			"content": cty.NullVal(cty.DynamicPseudoType),
			"format":  cty.NullVal(cty.String),
		}),
		Context: map[string]any{
			"query_result": "invalid_type",
		},
	}
	result := s.plugin.Call(args)
	expected := plugininterface.Result{
		Diags: hcl.Diagnostics{{
			Severity: hcl.DiagError,
			Summary:  "Failed to parse arguments",
			Detail:   "invalid query result: string",
		}},
	}
	s.Equal(expected, result)
}

func (s *PluginTestSuite) TestCallContentAndDataContextNil() {
	args := plugininterface.Args{
		Kind: "content",
		Name: "frontmatter",
		Args: cty.ObjectVal(map[string]cty.Value{
			"content": cty.NullVal(cty.DynamicPseudoType),
			"format":  cty.NullVal(cty.String),
		}),
		Context: nil,
	}
	result := s.plugin.Call(args)
	expected := plugininterface.Result{
		Diags: hcl.Diagnostics{{
			Severity: hcl.DiagError,
			Summary:  "Failed to parse arguments",
			Detail:   "query_result and content are nil",
		}},
	}
	s.Equal(expected, result)
}

func (s *PluginTestSuite) TestCallWithContent() {
	args := plugininterface.Args{
		Kind: "content",
		Name: "frontmatter",
		Args: cty.ObjectVal(map[string]cty.Value{
			"content": cty.ObjectVal(map[string]cty.Value{
				"baz": cty.NumberIntVal(1),
				"foo": cty.StringVal("bar"),
				"quux": cty.ObjectVal(map[string]cty.Value{
					"corge":  cty.StringVal("grault"),
					"garply": cty.BoolVal(false),
				}),
				"qux": cty.BoolVal(true),
				"waldo": cty.ListVal([]cty.Value{
					cty.StringVal("fred"),
					cty.StringVal("plugh"),
				}),
			}),
			"format": cty.NullVal(cty.String),
		}),
		Context: nil,
	}
	result := s.plugin.Call(args)
	expected := plugininterface.Result{
		Result: "---\n" +
			"baz: 1\n" +
			"foo: bar\n" +
			"quux:\n" +
			"    corge: grault\n" +
			"    garply: false\n" +
			"qux: true\n" +
			"waldo:\n" +
			"    - fred\n" +
			"    - plugh\n" +
			"---\n",
	}
	s.Equal(expected, result)
}

func (s *PluginTestSuite) TestCallWithQueryResult() {
	args := plugininterface.Args{
		Kind: "content",
		Name: "frontmatter",
		Args: cty.ObjectVal(map[string]cty.Value{
			"content": cty.NullVal(cty.DynamicPseudoType),
			"format":  cty.NullVal(cty.String),
		}),
		Context: map[string]any{
			"query_result": map[string]any{
				"baz": int64(1),
				"foo": "bar",
				"quux": map[string]any{
					"corge":  "grault",
					"garply": false,
				},
				"qux": true,
				"waldo": []any{
					"fred",
					"plugh",
				},
			},
		},
	}
	result := s.plugin.Call(args)
	expected := plugininterface.Result{
		Result: "---\n" +
			"baz: 1\n" +
			"foo: bar\n" +
			"quux:\n" +
			"    corge: grault\n" +
			"    garply: false\n" +
			"qux: true\n" +
			"waldo:\n" +
			"    - fred\n" +
			"    - plugh\n" +
			"---\n",
	}
	s.Equal(expected, result)
}

func (s *PluginTestSuite) TestCallFormatYaml() {
	args := plugininterface.Args{
		Kind: "content",
		Name: "frontmatter",
		Args: cty.ObjectVal(map[string]cty.Value{
			"content": cty.NullVal(cty.DynamicPseudoType),
			"format":  cty.StringVal("yaml"),
		}),
		Context: map[string]any{
			"query_result": map[string]any{
				"baz": int64(1),
				"foo": "bar",
				"quux": map[string]any{
					"corge":  "grault",
					"garply": false,
				},
				"qux": true,
				"waldo": []any{
					"fred",
					"plugh",
				},
			},
		},
	}
	result := s.plugin.Call(args)
	expected := plugininterface.Result{
		Result: "---\n" +
			"baz: 1\n" +
			"foo: bar\n" +
			"quux:\n" +
			"    corge: grault\n" +
			"    garply: false\n" +
			"qux: true\n" +
			"waldo:\n" +
			"    - fred\n" +
			"    - plugh\n" +
			"---\n",
	}
	s.Equal(expected, result)
}

func (s *PluginTestSuite) TestCallFormatTOML() {
	args := plugininterface.Args{
		Kind: "content",
		Name: "frontmatter",
		Args: cty.ObjectVal(map[string]cty.Value{
			"content": cty.NullVal(cty.DynamicPseudoType),
			"format":  cty.StringVal("toml"),
		}),
		Context: map[string]any{
			"query_result": map[string]any{
				"baz": int64(1),
				"foo": "bar",
				"quux": map[string]any{
					"corge":  "grault",
					"garply": false,
				},
				"qux": true,
				"waldo": []any{
					"fred",
					"plugh",
				},
			},
		},
	}
	result := s.plugin.Call(args)
	expected := plugininterface.Result{
		Result: "+++\n" +
			"baz = 1\n" +
			"foo = 'bar'\n" +
			"qux = true\n" +
			"waldo = ['fred', 'plugh']\n\n" +
			"[quux]\n" +
			"corge = 'grault'\n" +
			"garply = false\n" +
			"+++\n",
	}
	s.Equal(expected, result)
}

func (s *PluginTestSuite) TestCallFormatJSON() {
	args := plugininterface.Args{
		Kind: "content",
		Name: "frontmatter",
		Args: cty.ObjectVal(map[string]cty.Value{
			"content": cty.NullVal(cty.DynamicPseudoType),
			"format":  cty.StringVal("json"),
		}),
		Context: map[string]any{
			"query_result": map[string]any{
				"baz": int64(1),
				"foo": "bar",
				"quux": map[string]any{
					"corge":  "grault",
					"garply": false,
				},
				"qux": true,
				"waldo": []any{
					"fred",
					"plugh",
				},
			},
		},
	}
	result := s.plugin.Call(args)
	expected := plugininterface.Result{
		Result: "{\n" +
			"  \"baz\": 1,\n" +
			"  \"foo\": \"bar\",\n" +
			"  \"quux\": {\n" +
			"    \"corge\": \"grault\",\n" +
			"    \"garply\": false\n" +
			"  },\n" +
			"  \"qux\": true,\n" +
			"  \"waldo\": [\n" +
			"    \"fred\",\n" +
			"    \"plugh\"\n" +
			"  ]\n" +
			"}\n",
	}
	s.Equal(expected, result)
}

func (s *PluginTestSuite) TestCallFailedRender() {
	args := plugininterface.Args{
		Kind: "content",
		Name: "frontmatter",
		Args: cty.ObjectVal(map[string]cty.Value{
			"content": cty.NullVal(cty.DynamicPseudoType),
			"format":  cty.StringVal("json"),
		}),
		Context: map[string]any{
			"query_result": map[string]any{
				"foo": func() {},
			},
		},
	}
	result := s.plugin.Call(args)
	expected := plugininterface.Result{
		Diags: hcl.Diagnostics{{
			Severity: hcl.DiagError,
			Summary:  "Failed to render frontmatter",
			Detail:   "json: unsupported type: func()",
		}},
	}
	s.Equal(expected, result)
}
