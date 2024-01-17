package openai

import (
	"testing"

	"github.com/blackstork-io/fabric/plugininterface/v1"
	"github.com/blackstork-io/fabric/plugins/content/openai/client"
	"github.com/blackstork-io/fabric/plugins/content/openai/internal/mocks"
	"github.com/hashicorp/hcl/v2"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
	"github.com/zclconf/go-cty/cty"
)

type PluginTestSuite struct {
	suite.Suite
	plugin plugininterface.PluginRPC
	cli    *mocks.Client
}

func TestPluginSuite(t *testing.T) {
	suite.Run(t, &PluginTestSuite{})
}

func (s *PluginTestSuite) SetupSuite() {
	s.plugin = Plugin{
		ClientLoader: func(opts ...client.Option) client.Client {
			return s.cli
		},
	}
}

func (s *PluginTestSuite) SetupTest() {
	s.cli = &mocks.Client{}
}

func (s *PluginTestSuite) TearDownTest() {
	s.cli.AssertExpectations(s.T())
}

func (s *PluginTestSuite) TestGetPlugins() {
	plugins := s.plugin.GetPlugins()
	s.Require().Len(plugins, 1, "expected 1 plugin")
	got := plugins[0]
	s.Equal("openai_text", got.Name)
	s.Equal("content", got.Kind)
	s.Equal("blackstork", got.Namespace)
	s.Equal(Version.String(), got.Version.Cast().String())
	s.NotNil(got.ConfigSpec)
	s.NotNil(got.InvocationSpec)
}

func (s *PluginTestSuite) TestCallBasic() {

	s.cli.On("GenerateChatCompletion", mock.Anything, &client.ChatCompletionParams{
		Model: defaultModel,
		Messages: []client.ChatCompletionMessage{
			{
				Role:    "user",
				Content: "Tell me a story",
			},
		},
	}).Return(&client.ChatCompletionResult{
		Choices: []client.ChatCompletionChoice{
			{
				FinishedReason: "stop",
				Index:          0,
				Message: client.ChatCompletionMessage{
					Role:    "assistant",
					Content: "Once upon a time.",
				},
			},
		},
	}, nil)

	args := plugininterface.Args{
		Kind: "content",
		Name: "openai_text",
		Args: cty.ObjectVal(map[string]cty.Value{
			"prompt": cty.StringVal("Tell me a story"),
		}),
		Config: cty.ObjectVal(map[string]cty.Value{
			"api_key":         cty.StringVal("api_key_123"),
			"organization_id": cty.NullVal(cty.String),
			"system_prompt":   cty.NullVal(cty.String),
		}),
		Context: map[string]any{},
	}
	result := s.plugin.Call(args)
	expected := plugininterface.Result{
		Result: "Once upon a time.",
	}
	s.Equal(expected, result)
}
func (s *PluginTestSuite) TestCallAdvanced() {

	s.cli.On("GenerateChatCompletion", mock.Anything, &client.ChatCompletionParams{
		Model: defaultModel,
		Messages: []client.ChatCompletionMessage{
			{
				Role:    "system",
				Content: "Some system message.",
			},
			{
				Role:    "user",
				Content: "Tell me a story\n```\n{\n  \"foo\": \"bar\"\n}\n```",
			},
		},
	}).Return(&client.ChatCompletionResult{
		Choices: []client.ChatCompletionChoice{
			{
				FinishedReason: "stop",
				Index:          0,
				Message: client.ChatCompletionMessage{
					Role:    "assistant",
					Content: "Once upon a time.",
				},
			},
		},
	}, nil)

	args := plugininterface.Args{
		Kind: "content",
		Name: "openai_text",
		Args: cty.ObjectVal(map[string]cty.Value{
			"prompt": cty.StringVal("Tell me a story"),
		}),
		Config: cty.ObjectVal(map[string]cty.Value{
			"api_key":         cty.StringVal("api_key_123"),
			"organization_id": cty.StringVal("org_id_123"),
			"system_prompt":   cty.StringVal("Some system message."),
		}),
		Context: map[string]any{
			"query_result": map[string]any{
				"foo": "bar",
			},
		},
	}
	result := s.plugin.Call(args)
	expected := plugininterface.Result{
		Result: "Once upon a time.",
	}
	s.Equal(expected, result)
}

func (s *PluginTestSuite) TestCallMissingPrompt() {
	args := plugininterface.Args{
		Kind: "content",
		Name: "openai_text",
		Args: cty.ObjectVal(map[string]cty.Value{
			"prompt": cty.NullVal(cty.String),
		}),
		Config: cty.ObjectVal(map[string]cty.Value{
			"api_key":         cty.StringVal("api_key_123"),
			"organization_id": cty.NullVal(cty.String),
			"system_prompt":   cty.NullVal(cty.String),
		}),
		Context: map[string]any{},
	}
	result := s.plugin.Call(args)
	expected := plugininterface.Result{
		Diags: hcl.Diagnostics{{
			Severity: hcl.DiagError,
			Summary:  "Failed to generate text",
			Detail:   "prompt is required in invocation",
		}},
	}
	s.Equal(expected, result)
}

func (s *PluginTestSuite) TestCallMissingAPIKey() {
	args := plugininterface.Args{
		Kind: "content",
		Name: "openai_text",
		Args: cty.ObjectVal(map[string]cty.Value{
			"prompt": cty.StringVal("Tell me a story"),
		}),
		Config: cty.ObjectVal(map[string]cty.Value{
			"api_key":         cty.NullVal(cty.String),
			"organization_id": cty.NullVal(cty.String),
			"system_prompt":   cty.NullVal(cty.String),
		}),
		Context: map[string]any{},
	}
	result := s.plugin.Call(args)
	expected := plugininterface.Result{
		Diags: hcl.Diagnostics{{
			Severity: hcl.DiagError,
			Summary:  "Failed to create client",
			Detail:   "api_key is required in configuration",
		}},
	}
	s.Equal(expected, result)
}

func (s *PluginTestSuite) TestCallFailingClient() {
	args := plugininterface.Args{
		Kind: "content",
		Name: "openai_text",
		Args: cty.ObjectVal(map[string]cty.Value{
			"prompt": cty.StringVal("Tell me a story"),
		}),
		Config: cty.ObjectVal(map[string]cty.Value{
			"api_key":         cty.StringVal("api_key_123"),
			"organization_id": cty.NullVal(cty.String),
			"system_prompt":   cty.NullVal(cty.String),
		}),
		Context: map[string]any{},
	}
	s.cli.On("GenerateChatCompletion", mock.Anything, mock.Anything).Return(nil, client.Error{
		Type:    "error_type",
		Message: "error_message",
	})
	result := s.plugin.Call(args)
	expected := plugininterface.Result{
		Diags: hcl.Diagnostics{{
			Severity: hcl.DiagError,
			Summary:  "Failed to generate text",
			Detail:   "openai[error_type]: error_message",
		}},
	}
	s.Equal(expected, result)
}
