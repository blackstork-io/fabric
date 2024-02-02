package list

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
	s.Equal("list", got.Name)
	s.Equal("content", got.Kind)
	s.Equal("blackstork", got.Namespace)
	s.Equal(Version.String(), got.Version.Cast().String())
	s.Nil(got.ConfigSpec)
	s.NotNil(got.InvocationSpec)
}

func (s *PluginTestSuite) TestCallNilQueryResult() {
	args := plugininterface.Args{
		Kind: "content",
		Name: "list",
		Args: cty.ObjectVal(map[string]cty.Value{
			"item_template": cty.StringVal("{{.}}"),
			"format":        cty.NullVal(cty.String),
		}),
		Context: map[string]any{
			"foo":          "foo_value",
			"query_result": nil,
		},
	}
	result := s.plugin.Call(args)
	expected := plugininterface.Result{
		Diags: hcl.Diagnostics{{
			Severity: hcl.DiagError,
			Summary:  "Failed to render template",
			Detail:   "query_result is required in data context",
		}},
	}
	s.Equal(expected, result)
}

func (s *PluginTestSuite) TestCallNonArrayQueryResult() {
	args := plugininterface.Args{
		Kind: "content",
		Name: "list",
		Args: cty.ObjectVal(map[string]cty.Value{
			"item_template": cty.StringVal("{{.}}"),
			"format":        cty.NullVal(cty.String),
		}),
		Context: map[string]any{
			"query_result": "not_an_array",
		},
	}
	result := s.plugin.Call(args)
	expected := plugininterface.Result{
		Diags: hcl.Diagnostics{{
			Severity: hcl.DiagError,
			Summary:  "Failed to render template",
			Detail:   "query_result must be an array",
		}},
	}
	s.Equal(expected, result)
}

func (s *PluginTestSuite) TestCallOrdered() {
	args := plugininterface.Args{
		Kind: "content",
		Name: "list",
		Args: cty.ObjectVal(map[string]cty.Value{
			"item_template": cty.StringVal("foo {{.}}"),
			"format":        cty.StringVal("ordered"),
		}),
		Context: map[string]any{
			"query_result": []any{
				"bar",
				"baz",
			},
		},
	}
	result := s.plugin.Call(args)
	expected := plugininterface.Result{
		Result: "1. foo bar\n2. foo baz\n",
	}
	s.Equal(expected, result)
}

func (s *PluginTestSuite) TestCallTaskList() {
	args := plugininterface.Args{
		Kind: "content",
		Name: "list",
		Args: cty.ObjectVal(map[string]cty.Value{
			"item_template": cty.StringVal("foo {{.}}"),
			"format":        cty.StringVal("tasklist"),
		}),
		Context: map[string]any{
			"query_result": []any{
				"bar",
				"baz",
			},
		},
	}
	result := s.plugin.Call(args)
	expected := plugininterface.Result{
		Result: "* [ ] foo bar\n* [ ] foo baz\n",
	}
	s.Equal(expected, result)
}

func (s *PluginTestSuite) TestCallBasic() {
	args := plugininterface.Args{
		Kind: "content",
		Name: "list",
		Args: cty.ObjectVal(map[string]cty.Value{
			"item_template": cty.StringVal("foo {{.}}"),
			"format":        cty.NullVal(cty.String),
		}),
		Context: map[string]any{
			"query_result": []any{
				"bar",
				"baz",
			},
		},
	}
	result := s.plugin.Call(args)
	expected := plugininterface.Result{
		Result: "* foo bar\n* foo baz\n",
	}
	s.Equal(expected, result)
}
func (s *PluginTestSuite) TestCallAdvanced() {
	args := plugininterface.Args{
		Kind: "content",
		Name: "list",
		Args: cty.ObjectVal(map[string]cty.Value{
			"item_template": cty.StringVal("foo {{.bar}} {{.baz}}"),
			"format":        cty.NullVal(cty.String),
		}),
		Context: map[string]any{
			"query_result": []any{
				map[string]any{
					"bar": "bar1",
					"baz": "baz1",
				},
				map[string]any{
					"bar": "bar2",
					"baz": "baz2",
				},
			},
		},
	}
	result := s.plugin.Call(args)
	expected := plugininterface.Result{
		Result: "* foo bar1 baz1\n* foo bar2 baz2\n",
	}
	s.Equal(expected, result)
}

func (s *PluginTestSuite) TestCallEmptyQueryResult() {
	args := plugininterface.Args{
		Kind: "content",
		Name: "list",
		Args: cty.ObjectVal(map[string]cty.Value{
			"item_template": cty.StringVal("foo {{.}}"),
			"format":        cty.NullVal(cty.String),
		}),
		Context: map[string]any{
			"query_result": []any{},
		},
	}
	result := s.plugin.Call(args)
	expected := plugininterface.Result{
		Result: "",
	}
	s.Equal(expected, result)
}

func (s *PluginTestSuite) TestCallMissingItemTemplate() {
	args := plugininterface.Args{
		Kind: "content",
		Name: "list",
		Args: cty.ObjectVal(map[string]cty.Value{
			"item_template": cty.NullVal(cty.String),
			"format":        cty.NullVal(cty.String),
		}),
	}
	result := s.plugin.Call(args)
	expected := plugininterface.Result{
		Diags: hcl.Diagnostics{{
			Severity: hcl.DiagError,
			Summary:  "Failed to parse template",
			Detail:   "item_template is required",
		}},
	}
	s.Equal(expected, result)
}

func (s *PluginTestSuite) TestCallInvalidItemTemplate() {
	args := plugininterface.Args{
		Kind: "content",
		Name: "list",
		Args: cty.ObjectVal(map[string]cty.Value{
			"item_template": cty.StringVal("{{"),
			"format":        cty.NullVal(cty.String),
		}),
	}
	result := s.plugin.Call(args)
	expected := plugininterface.Result{
		Diags: hcl.Diagnostics{{
			Severity: hcl.DiagError,
			Summary:  "Failed to parse template",
			Detail:   "template: item:1: unclosed action",
		}},
	}
	s.Equal(expected, result)
}

func (s *PluginTestSuite) TestCallMissingDataContext() {
	args := plugininterface.Args{
		Kind: "content",
		Name: "list",
		Args: cty.ObjectVal(map[string]cty.Value{
			"item_template": cty.StringVal("foo {{.}}"),
			"format":        cty.NullVal(cty.String),
		}),
	}
	result := s.plugin.Call(args)
	expected := plugininterface.Result{
		Diags: hcl.Diagnostics{{
			Severity: hcl.DiagError,
			Summary:  "Failed to render template",
			Detail:   "data context is required",
		}},
	}
	s.Equal(expected, result)
}
