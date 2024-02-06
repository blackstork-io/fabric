package openai

import (
	"context"
	"testing"

	"github.com/hashicorp/hcl/v2"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
	"github.com/zclconf/go-cty/cty"

	"github.com/blackstork-io/fabric/internal/openai/client"
	client_mocks "github.com/blackstork-io/fabric/mocks/internalpkg/openai/client"
	"github.com/blackstork-io/fabric/plugin"
)

type OpenAITextContentTestSuite struct {
	suite.Suite
	schema *plugin.Schema
	cli    *client_mocks.Client
}

func TestOpenAITextContentSuite(t *testing.T) {
	suite.Run(t, &OpenAITextContentTestSuite{})
}

func (s *OpenAITextContentTestSuite) SetupSuite() {
	s.schema = Plugin("", func(opts ...client.Option) client.Client {
		return s.cli
	})
}

func (s *OpenAITextContentTestSuite) SetupTest() {
	s.cli = &client_mocks.Client{}
}

func (s *OpenAITextContentTestSuite) TearDownTest() {
	s.cli.AssertExpectations(s.T())
}

func (s *OpenAITextContentTestSuite) TestSchema() {
	provider := s.schema.ContentProviders["openai_text"]
	s.Require().NotNil(provider)
	s.NotNil(provider.Args)
	s.NotNil(provider.ContentFunc)
	s.NotNil(provider.Config)
}

func (s *OpenAITextContentTestSuite) TestBasic() {
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
	ctx := context.Background()
	content, diags := s.schema.ProvideContent(ctx, "openai_text", &plugin.ProvideContentParams{
		Args: cty.ObjectVal(map[string]cty.Value{
			"prompt": cty.StringVal("Tell me a story"),
			"model":  cty.NullVal(cty.String),
		}),
		Config: cty.ObjectVal(map[string]cty.Value{
			"api_key":         cty.StringVal("api_key_123"),
			"organization_id": cty.NullVal(cty.String),
			"system_prompt":   cty.NullVal(cty.String),
		}),
		DataContext: plugin.MapData{},
	})
	s.Nil(diags)
	s.Equal(&plugin.Content{
		Markdown: "Once upon a time.",
	}, content)
}

func (s *OpenAITextContentTestSuite) TestAdvanced() {
	s.cli.On("GenerateChatCompletion", mock.Anything, &client.ChatCompletionParams{
		Model: "model_123",
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
	ctx := context.Background()
	content, diags := s.schema.ProvideContent(ctx, "openai_text", &plugin.ProvideContentParams{
		Args: cty.ObjectVal(map[string]cty.Value{
			"prompt": cty.StringVal("Tell me a story"),
			"model":  cty.StringVal("model_123"),
		}),
		Config: cty.ObjectVal(map[string]cty.Value{
			"api_key":         cty.StringVal("api_key_123"),
			"organization_id": cty.StringVal("org_id_123"),
			"system_prompt":   cty.StringVal("Some system message."),
		}),
		DataContext: plugin.MapData{
			"query_result": plugin.MapData{
				"foo": plugin.StringData("bar"),
			},
		},
	})
	s.Nil(diags)
	s.Equal(&plugin.Content{
		Markdown: "Once upon a time.",
	}, content)
}

func (s *OpenAITextContentTestSuite) TestMissingPrompt() {
	ctx := context.Background()
	content, diags := s.schema.ProvideContent(ctx, "openai_text", &plugin.ProvideContentParams{
		Args: cty.ObjectVal(map[string]cty.Value{
			"prompt": cty.NullVal(cty.String),
			"model":  cty.NullVal(cty.String),
		}),
		Config: cty.ObjectVal(map[string]cty.Value{
			"api_key":         cty.StringVal("api_key_123"),
			"organization_id": cty.NullVal(cty.String),
			"system_prompt":   cty.NullVal(cty.String),
		}),
		DataContext: plugin.MapData{},
	})
	s.Nil(content)
	s.Equal(hcl.Diagnostics{{
		Severity: hcl.DiagError,
		Summary:  "Failed to generate text",
		Detail:   "prompt is required in invocation",
	}}, diags)
}

func (s *OpenAITextContentTestSuite) TestMissingAPIKey() {
	ctx := context.Background()
	content, diags := s.schema.ProvideContent(ctx, "openai_text", &plugin.ProvideContentParams{
		Args: cty.ObjectVal(map[string]cty.Value{
			"prompt": cty.StringVal("Tell me a story"),
			"model":  cty.NullVal(cty.String),
		}),
		Config: cty.ObjectVal(map[string]cty.Value{
			"api_key":         cty.NullVal(cty.String),
			"organization_id": cty.NullVal(cty.String),
			"system_prompt":   cty.NullVal(cty.String),
		}),
		DataContext: plugin.MapData{},
	})
	s.Nil(content)
	s.Equal(hcl.Diagnostics{{
		Severity: hcl.DiagError,
		Summary:  "Failed to create client",
		Detail:   "api_key is required in configuration",
	}}, diags)
}

func (s *OpenAITextContentTestSuite) TestFailingClient() {
	s.cli.On("GenerateChatCompletion", mock.Anything, mock.Anything).Return(nil, client.Error{
		Type:    "error_type",
		Message: "error_message",
	})
	ctx := context.Background()
	content, diags := s.schema.ProvideContent(ctx, "openai_text", &plugin.ProvideContentParams{
		Args: cty.ObjectVal(map[string]cty.Value{
			"prompt": cty.StringVal("Tell me a story"),
			"model":  cty.NullVal(cty.String),
		}),
		Config: cty.ObjectVal(map[string]cty.Value{
			"api_key":         cty.StringVal("api_key_123"),
			"organization_id": cty.NullVal(cty.String),
			"system_prompt":   cty.NullVal(cty.String),
		}),
		DataContext: plugin.MapData{},
	})
	s.Nil(content)
	s.Equal(hcl.Diagnostics{{
		Severity: hcl.DiagError,
		Summary:  "Failed to generate text",
		Detail:   "openai[error_type]: error_message",
	}}, diags)
}

func (s *OpenAITextContentTestSuite) TestCancellation() {
	s.cli.On("GenerateChatCompletion", mock.Anything, mock.Anything).Return(nil, context.Canceled)
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	content, diags := s.schema.ProvideContent(ctx, "openai_text", &plugin.ProvideContentParams{
		Args: cty.ObjectVal(map[string]cty.Value{
			"prompt": cty.StringVal("Tell me a story"),
			"model":  cty.NullVal(cty.String),
		}),
		Config: cty.ObjectVal(map[string]cty.Value{
			"api_key":         cty.StringVal("api_key_123"),
			"organization_id": cty.NullVal(cty.String),
			"system_prompt":   cty.NullVal(cty.String),
		}),
		DataContext: plugin.MapData{},
	})
	s.Nil(content)
	s.Equal(hcl.Diagnostics{{
		Severity: hcl.DiagError,
		Summary:  "Failed to generate text",
		Detail:   "context canceled",
	}}, diags)
}

func (s *OpenAITextContentTestSuite) TestErrorEncoding() {
	want := client.Error{
		Type:    "invalid_request_error",
		Message: "message of error",
	}
	s.cli.On("GenerateChatCompletion", mock.Anything, mock.Anything).Return(nil, want)
	ctx := context.Background()
	content, diags := s.schema.ProvideContent(ctx, "openai_text", &plugin.ProvideContentParams{
		Args: cty.ObjectVal(map[string]cty.Value{
			"prompt": cty.StringVal("Tell me a story"),
			"model":  cty.NullVal(cty.String),
		}),
		Config: cty.ObjectVal(map[string]cty.Value{
			"api_key":         cty.StringVal("api_key_123"),
			"organization_id": cty.NullVal(cty.String),
			"system_prompt":   cty.NullVal(cty.String),
		}),
		DataContext: plugin.MapData{},
	})
	s.Nil(content)
	s.Equal(hcl.Diagnostics{{
		Severity: hcl.DiagError,
		Summary:  "Failed to generate text",
		Detail:   "openai[invalid_request_error]: message of error",
	}}, diags)
}
