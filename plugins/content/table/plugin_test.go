package table

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
	s.Equal("table", got.Name)
	s.Equal("content", got.Kind)
	s.Equal("blackstork", got.Namespace)
	s.Equal(Version.String(), got.Version.Cast().String())
	s.Nil(got.ConfigSpec)
	s.NotNil(got.InvocationSpec)
}

func (s *PluginTestSuite) TestCallNilQueryResult() {
	args := plugininterface.Args{
		Kind: "content",
		Name: "table",
		Args: cty.ObjectVal(map[string]cty.Value{
			"columns": cty.ListVal([]cty.Value{
				cty.ObjectVal(map[string]cty.Value{
					"header": cty.StringVal("{{.col_prefix}} Name"),
					"value":  cty.StringVal("{{.name}}"),
				}),
				cty.ObjectVal(map[string]cty.Value{
					"header": cty.StringVal("{{.col_prefix}} Age"),
					"value":  cty.StringVal("{{.age}}"),
				}),
			}),
		}),
		Context: map[string]any{
			"col_prefix":   "User",
			"query_result": nil,
		},
	}
	result := s.plugin.Call(args)
	expected := plugininterface.Result{
		Result: "|User Name|User Age|\n|-|-|\n",
	}
	s.Equal(expected, result)
}

func (s *PluginTestSuite) TestCallEmptyQueryResult() {
	args := plugininterface.Args{
		Kind: "content",
		Name: "table",
		Args: cty.ObjectVal(map[string]cty.Value{
			"columns": cty.ListVal([]cty.Value{
				cty.ObjectVal(map[string]cty.Value{
					"header": cty.StringVal("{{.col_prefix}} Name"),
					"value":  cty.StringVal("{{.name}}"),
				}),
				cty.ObjectVal(map[string]cty.Value{
					"header": cty.StringVal("{{.col_prefix}} Age"),
					"value":  cty.StringVal("{{.age}}"),
				}),
			}),
		}),
		Context: map[string]any{
			"col_prefix":   "User",
			"query_result": []any{},
		},
	}
	result := s.plugin.Call(args)
	expected := plugininterface.Result{
		Result: "|User Name|User Age|\n|-|-|\n",
	}
	s.Equal(expected, result)
}
func (s *PluginTestSuite) TestCallBasic() {
	args := plugininterface.Args{
		Kind: "content",
		Name: "table",
		Args: cty.ObjectVal(map[string]cty.Value{
			"columns": cty.ListVal([]cty.Value{
				cty.ObjectVal(map[string]cty.Value{
					"header": cty.StringVal("{{.col_prefix}} Name"),
					"value":  cty.StringVal("{{.name}}"),
				}),
				cty.ObjectVal(map[string]cty.Value{
					"header": cty.StringVal("{{.col_prefix}} Age"),
					"value":  cty.StringVal("{{.age}}"),
				}),
			}),
		}),
		Context: map[string]any{
			"col_prefix": "User",
			"query_result": []any{
				map[string]any{
					"name": "John",
					"age":  42,
				},
				map[string]any{
					"name": "Jane",
					"age":  43,
				},
			},
		},
	}
	result := s.plugin.Call(args)
	expected := plugininterface.Result{
		Result: "|User Name|User Age|\n|-|-|\n|John|42|\n|Jane|43|\n",
	}
	s.Equal(expected, result)
}

func (s *PluginTestSuite) TestCallMissingHeader() {
	args := plugininterface.Args{
		Kind: "content",
		Name: "table",
		Args: cty.ObjectVal(map[string]cty.Value{
			"columns": cty.ListVal([]cty.Value{
				cty.ObjectVal(map[string]cty.Value{
					"value": cty.StringVal("{{.name}}"),
				}),
				cty.ObjectVal(map[string]cty.Value{
					"value": cty.StringVal("{{.age}}"),
				}),
			}),
		}),
		Context: nil,
	}
	result := s.plugin.Call(args)
	expected := plugininterface.Result{
		Diags: hcl.Diagnostics{{
			Severity: hcl.DiagError,
			Summary:  "Failed to parse arguments",
			Detail:   "missing header in table cell",
		}},
	}
	s.Equal(expected, result)
}
func (s *PluginTestSuite) TestCallMissingValue() {
	args := plugininterface.Args{
		Kind: "content",
		Name: "table",
		Args: cty.ObjectVal(map[string]cty.Value{
			"columns": cty.ListVal([]cty.Value{
				cty.ObjectVal(map[string]cty.Value{
					"header": cty.StringVal("{{.col_prefix}} Name"),
				}),
				cty.ObjectVal(map[string]cty.Value{
					"header": cty.StringVal("{{.col_prefix}} Age"),
				}),
			}),
		}),
		Context: nil,
	}
	result := s.plugin.Call(args)
	expected := plugininterface.Result{
		Diags: hcl.Diagnostics{{
			Severity: hcl.DiagError,
			Summary:  "Failed to parse arguments",
			Detail:   "missing value in table cell",
		}},
	}
	s.Equal(expected, result)
}

func (s *PluginTestSuite) TestCallNilHeader() {
	args := plugininterface.Args{
		Kind: "content",
		Name: "table",
		Args: cty.ObjectVal(map[string]cty.Value{
			"columns": cty.ListVal([]cty.Value{
				cty.ObjectVal(map[string]cty.Value{
					"header": cty.NullVal(cty.String),
					"value":  cty.StringVal("{{.name}}"),
				}),
				cty.ObjectVal(map[string]cty.Value{
					"header": cty.NullVal(cty.String),
					"value":  cty.StringVal("{{.age}}"),
				}),
			}),
		}),
		Context: nil,
	}
	result := s.plugin.Call(args)
	expected := plugininterface.Result{
		Diags: hcl.Diagnostics{{
			Severity: hcl.DiagError,
			Summary:  "Failed to parse arguments",
			Detail:   "missing header in table cell",
		}},
	}
	s.Equal(expected, result)
}
func (s *PluginTestSuite) TestCallNilValue() {
	args := plugininterface.Args{
		Kind: "content",
		Name: "table",
		Args: cty.ObjectVal(map[string]cty.Value{
			"columns": cty.ListVal([]cty.Value{
				cty.ObjectVal(map[string]cty.Value{
					"header": cty.StringVal("{{.col_prefix}} Name"),
					"value":  cty.NullVal(cty.String),
				}),
				cty.ObjectVal(map[string]cty.Value{
					"header": cty.StringVal("{{.col_prefix}} Age"),
					"value":  cty.NullVal(cty.String),
				}),
			}),
		}),
		Context: nil,
	}
	result := s.plugin.Call(args)
	expected := plugininterface.Result{
		Diags: hcl.Diagnostics{{
			Severity: hcl.DiagError,
			Summary:  "Failed to parse arguments",
			Detail:   "missing value in table cell",
		}},
	}
	s.Equal(expected, result)
}

func (s *PluginTestSuite) TestCallNilColumns() {
	args := plugininterface.Args{
		Kind: "content",
		Name: "table",
		Args: cty.ObjectVal(map[string]cty.Value{
			"columns": cty.NullVal(cty.List(cty.Object(map[string]cty.Type{}))),
		}),
	}
	result := s.plugin.Call(args)
	expected := plugininterface.Result{
		Diags: hcl.Diagnostics{{
			Severity: hcl.DiagError,
			Summary:  "Failed to parse arguments",
			Detail:   "columns is required",
		}},
	}
	s.Equal(expected, result)
}

func (s *PluginTestSuite) TestCallEmptyColumns() {
	args := plugininterface.Args{
		Kind: "content",
		Name: "table",
		Args: cty.ObjectVal(map[string]cty.Value{
			"columns": cty.ListValEmpty(cty.Object(map[string]cty.Type{})),
		}),
	}
	result := s.plugin.Call(args)
	expected := plugininterface.Result{
		Diags: hcl.Diagnostics{{
			Severity: hcl.DiagError,
			Summary:  "Failed to parse arguments",
			Detail:   "columns must not be empty",
		}},
	}
	s.Equal(expected, result)
}

func (s *PluginTestSuite) TestCallInvalidHeaderTemplate() {
	args := plugininterface.Args{
		Kind: "content",
		Name: "table",
		Args: cty.ObjectVal(map[string]cty.Value{
			"columns": cty.ListVal([]cty.Value{
				cty.ObjectVal(map[string]cty.Value{
					"header": cty.StringVal("{{.col_prefix} Name"),
					"value":  cty.StringVal("{{.name}}"),
				}),
				cty.ObjectVal(map[string]cty.Value{
					"header": cty.StringVal("{{.col_prefix}} Age"),
					"value":  cty.StringVal("{{.age}}"),
				}),
			}),
		}),
	}
	result := s.plugin.Call(args)
	expected := plugininterface.Result{
		Diags: hcl.Diagnostics{{
			Severity: hcl.DiagError,
			Summary:  "Failed to parse arguments",
			Detail:   "failed to parse header template: template: header:1: bad character U+007D '}'",
		}},
	}
	s.Equal(expected, result)
}

func (s *PluginTestSuite) TestCallInvalidValueTemplate() {
	args := plugininterface.Args{
		Kind: "content",
		Name: "table",
		Args: cty.ObjectVal(map[string]cty.Value{
			"columns": cty.ListVal([]cty.Value{
				cty.ObjectVal(map[string]cty.Value{
					"header": cty.StringVal("{{.col_prefix}} Name"),
					"value":  cty.StringVal("{{.name}"),
				}),
				cty.ObjectVal(map[string]cty.Value{
					"header": cty.StringVal("{{.col_prefix}} Age"),
					"value":  cty.StringVal("{{.age}}"),
				}),
			}),
		}),
	}
	result := s.plugin.Call(args)
	expected := plugininterface.Result{
		Diags: hcl.Diagnostics{{
			Severity: hcl.DiagError,
			Summary:  "Failed to parse arguments",
			Detail:   "failed to parse value template: template: value:1: bad character U+007D '}'",
		}},
	}
	s.Equal(expected, result)
}
