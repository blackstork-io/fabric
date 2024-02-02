package text

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
	s.Equal("text", got.Name)
	s.Equal("content", got.Kind)
	s.Equal("blackstork", got.Namespace)
	s.Equal(Version.String(), got.Version.Cast().String())
	s.Nil(got.ConfigSpec)
	s.NotNil(got.InvocationSpec)
}

func (s *PluginTestSuite) TestCallMissingText() {
	args := plugininterface.Args{
		Args: cty.ObjectVal(map[string]cty.Value{
			"text":                cty.NullVal(cty.String),
			"format_as":           cty.NullVal(cty.String),
			"absolute_title_size": cty.NullVal(cty.Number),
			"code_language":       cty.NullVal(cty.String),
		}),
		Context: map[string]any{
			"name": "World",
		},
	}
	expected := plugininterface.Result{
		Diags: hcl.Diagnostics{{
			Severity: hcl.DiagError,
			Summary:  "Failed to render text",
			Detail:   "text is required",
		}},
	}
	s.Equal(expected, s.plugin.Call(args))
}
func (s *PluginTestSuite) TestCallText() {
	args := plugininterface.Args{
		Args: cty.ObjectVal(map[string]cty.Value{
			"text":                cty.StringVal("Hello {{.name}}!"),
			"format_as":           cty.NullVal(cty.String),
			"absolute_title_size": cty.NullVal(cty.Number),
			"code_language":       cty.NullVal(cty.String),
		}),
		Context: map[string]any{
			"name": "World",
		},
	}
	expected := plugininterface.Result{
		Result: "Hello World!",
	}
	s.Equal(expected, s.plugin.Call(args))
}

func (s *PluginTestSuite) TestCallTextNoTemplate() {
	args := plugininterface.Args{
		Args: cty.ObjectVal(map[string]cty.Value{
			"text":                cty.StringVal("Hello World!"),
			"format_as":           cty.NullVal(cty.String),
			"absolute_title_size": cty.NullVal(cty.Number),
			"code_language":       cty.NullVal(cty.String),
		}),
		Context: nil,
	}
	expected := plugininterface.Result{
		Result: "Hello World!",
	}
	s.Equal(expected, s.plugin.Call(args))
}

func (s *PluginTestSuite) TestCallTitleDefault() {
	args := plugininterface.Args{
		Args: cty.ObjectVal(map[string]cty.Value{
			"text":                cty.StringVal("Hello {{.name}}!"),
			"format_as":           cty.StringVal("title"),
			"absolute_title_size": cty.NullVal(cty.Number),
			"code_language":       cty.NullVal(cty.String),
		}),
		Context: map[string]any{
			"name": "World",
		},
	}
	expected := plugininterface.Result{
		Result: "# Hello World!",
	}
	s.Equal(expected, s.plugin.Call(args))
}

func (s *PluginTestSuite) TestCallTitleWithTextMultiline() {
	args := plugininterface.Args{
		Args: cty.ObjectVal(map[string]cty.Value{
			"text":                cty.StringVal("Hello\n{{.name}}\nfor you!"),
			"format_as":           cty.StringVal("title"),
			"absolute_title_size": cty.NullVal(cty.Number),
			"code_language":       cty.NullVal(cty.String),
		}),
		Context: map[string]any{
			"name": "World",
		},
	}
	expected := plugininterface.Result{
		Result: "# Hello World for you!",
	}
	s.Equal(expected, s.plugin.Call(args))
}

func (s *PluginTestSuite) TestCallTitleWithSize() {
	args := plugininterface.Args{
		Args: cty.ObjectVal(map[string]cty.Value{
			"text":                cty.StringVal("Hello {{.name}}!"),
			"format_as":           cty.StringVal("title"),
			"absolute_title_size": cty.NumberIntVal(3),
			"code_language":       cty.NullVal(cty.String),
		}),
		Context: map[string]any{
			"name": "World",
		},
	}
	expected := plugininterface.Result{
		Result: "### Hello World!",
	}
	s.Equal(expected, s.plugin.Call(args))
}

func (s *PluginTestSuite) TestCallTitleWithSizeTooSmall() {
	args := plugininterface.Args{
		Args: cty.ObjectVal(map[string]cty.Value{
			"text":                cty.StringVal("Hello {{.name}}!"),
			"format_as":           cty.StringVal("title"),
			"absolute_title_size": cty.NumberIntVal(0),
			"code_language":       cty.NullVal(cty.String),
		}),
		Context: map[string]any{
			"name": "World",
		},
	}
	expected := plugininterface.Result{
		Diags: hcl.Diagnostics{{
			Severity: hcl.DiagError,
			Summary:  "Failed to render text",
			Detail:   "absolute_title_size must be between 1 and 6",
		}},
	}
	s.Equal(expected, s.plugin.Call(args))
}

func (s *PluginTestSuite) TestCallTitleWithSizeTooBig() {
	args := plugininterface.Args{
		Args: cty.ObjectVal(map[string]cty.Value{
			"text":                cty.StringVal("Hello {{.name}}!"),
			"format_as":           cty.StringVal("title"),
			"absolute_title_size": cty.NumberIntVal(7),
			"code_language":       cty.NullVal(cty.String),
		}),
		Context: map[string]any{
			"name": "World",
		},
	}
	expected := plugininterface.Result{
		Diags: hcl.Diagnostics{{
			Severity: hcl.DiagError,
			Summary:  "Failed to render text",
			Detail:   "absolute_title_size must be between 1 and 6",
		}},
	}
	s.Equal(expected, s.plugin.Call(args))
}

func (s *PluginTestSuite) TestCallInvalidFormat() {
	args := plugininterface.Args{
		Args: cty.ObjectVal(map[string]cty.Value{
			"text":                cty.StringVal("Hello World!"),
			"format_as":           cty.StringVal("unknown"),
			"absolute_title_size": cty.NullVal(cty.Number),
			"code_language":       cty.NullVal(cty.String),
		}),
		Context: nil,
	}
	expected := plugininterface.Result{
		Diags: hcl.Diagnostics{{
			Severity: hcl.DiagError,
			Summary:  "Failed to render text",
			Detail:   "format_as must be one of text, title, code, blockquote",
		}},
	}
	s.Equal(expected, s.plugin.Call(args))
}

func (s *PluginTestSuite) TestCallInvalidTemplate() {
	args := plugininterface.Args{
		Args: cty.ObjectVal(map[string]cty.Value{
			"text":                cty.StringVal("Hello {{.name}!"),
			"format_as":           cty.NullVal(cty.String),
			"absolute_title_size": cty.NullVal(cty.Number),
			"code_language":       cty.NullVal(cty.String),
		}),
		Context: map[string]any{
			"name": "World",
		},
	}
	expected := plugininterface.Result{
		Diags: hcl.Diagnostics{{
			Severity: hcl.DiagError,
			Summary:  "Failed to render text",
			Detail:   "failed to parse text template: template: text:1: bad character U+007D '}'",
		}},
	}
	s.Equal(expected, s.plugin.Call(args))
}

func (s *PluginTestSuite) TestCallCodeDefault() {
	args := plugininterface.Args{
		Args: cty.ObjectVal(map[string]cty.Value{
			"text":                cty.StringVal(`Hello {{.name}}!`),
			"format_as":           cty.StringVal("code"),
			"absolute_title_size": cty.NullVal(cty.Number),
			"code_language":       cty.NullVal(cty.String),
		}),
		Context: map[string]any{
			"name": "World",
		},
	}
	expected := plugininterface.Result{
		Result: "```\nHello World!\n```",
	}
	s.Equal(expected, s.plugin.Call(args))
}

func (s *PluginTestSuite) TestCallCodeNoLanguage() {
	args := plugininterface.Args{
		Args: cty.ObjectVal(map[string]cty.Value{
			"text":                cty.StringVal(`{"hello": "{{.name}}"}`),
			"format_as":           cty.StringVal("code"),
			"absolute_title_size": cty.NullVal(cty.Number),
			"code_language":       cty.StringVal("json"),
		}),
		Context: map[string]any{
			"name": "world",
		},
	}
	expected := plugininterface.Result{
		Result: "```json\n{\"hello\": \"world\"}\n```",
	}
	s.Equal(expected, s.plugin.Call(args))
}

func (s *PluginTestSuite) TestCallBlockquote() {
	args := plugininterface.Args{
		Args: cty.ObjectVal(map[string]cty.Value{
			"text":                cty.StringVal(`Hello {{.name}}!`),
			"format_as":           cty.StringVal("blockquote"),
			"absolute_title_size": cty.NullVal(cty.Number),
			"code_language":       cty.NullVal(cty.String),
		}),
		Context: map[string]any{
			"name": "World",
		},
	}
	expected := plugininterface.Result{
		Result: "> Hello World!",
	}
	s.Equal(expected, s.plugin.Call(args))
}

func (s *PluginTestSuite) TestCallBlockquoteMultiline() {
	args := plugininterface.Args{
		Args: cty.ObjectVal(map[string]cty.Value{
			"text":                cty.StringVal("Hello\n{{.name}}\nfor you!"),
			"format_as":           cty.StringVal("blockquote"),
			"absolute_title_size": cty.NullVal(cty.Number),
			"code_language":       cty.NullVal(cty.String),
		}),
		Context: map[string]any{
			"name": "World",
		},
	}
	expected := plugininterface.Result{
		Result: "> Hello\n> World\n> for you!",
	}
	s.Equal(expected, s.plugin.Call(args))
}

func (s *PluginTestSuite) TestCallBlockquoteMultilineDoubleNewline() {
	args := plugininterface.Args{
		Args: cty.ObjectVal(map[string]cty.Value{
			"text":                cty.StringVal("Hello\n{{.name}}\n\nfor you!"),
			"format_as":           cty.StringVal("blockquote"),
			"absolute_title_size": cty.NullVal(cty.Number),
			"code_language":       cty.NullVal(cty.String),
		}),
		Context: map[string]any{
			"name": "World",
		},
	}
	expected := plugininterface.Result{
		Result: "> Hello\n> World\n> \n> for you!",
	}
	s.Equal(expected, s.plugin.Call(args))
}
