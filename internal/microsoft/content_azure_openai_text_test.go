package microsoft_test

import (
	"context"
	"errors"
	"fmt"
	"testing"

	"github.com/Azure/azure-sdk-for-go/sdk/ai/azopenai"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	"github.com/hashicorp/hcl/v2"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
	"github.com/zclconf/go-cty/cty"

	"github.com/blackstork-io/fabric/internal/microsoft"
	client_mocks "github.com/blackstork-io/fabric/mocks/internalpkg/microsoft"
	"github.com/blackstork-io/fabric/pkg/diagnostics"
	"github.com/blackstork-io/fabric/pkg/diagnostics/diagtest"
	"github.com/blackstork-io/fabric/plugin"
	"github.com/blackstork-io/fabric/plugin/plugintest"
	"github.com/blackstork-io/fabric/print/mdprint"
)

type AzureOpenAITextContentTestSuite struct {
	suite.Suite
	plugin *plugin.Schema
	schema *plugin.ContentProvider
	cli    *client_mocks.AzureOpenaiClient
}

func TestAzureOpenAITextContentSuite(t *testing.T) {
	suite.Run(t, &AzureOpenAITextContentTestSuite{})
}

func (s *AzureOpenAITextContentTestSuite) SetupSuite() {
	s.plugin = microsoft.Plugin("1.0.0", nil, (func(apiKey string, endPoint string) (cli microsoft.AzureOpenaiClient, err error) {
		return s.cli, nil
	}))
	s.schema = s.plugin.ContentProviders["azure_openai_text"]
}

func (s *AzureOpenAITextContentTestSuite) SetupTest() {
	s.cli = &client_mocks.AzureOpenaiClient{}
}

func (s *AzureOpenAITextContentTestSuite) TearDownTest() {
	s.cli.AssertExpectations(s.T())
}

func (s *AzureOpenAITextContentTestSuite) TestSchema() {
	s.Require().NotNil(s.plugin)
	s.Require().NotNil(s.schema)
	s.NotNil(s.schema.Args)
	s.NotNil(s.schema.ContentFunc)
	s.NotNil(s.schema.Config)
}

func (s *AzureOpenAITextContentTestSuite) TestBasic() {
	s.cli.On("GetCompletions", mock.Anything, azopenai.CompletionsOptions{
		DeploymentName: to.Ptr("test"),
		MaxTokens:      to.Ptr(int32(1000)),
		Temperature:    to.Ptr(float32(0)),
		Prompt:         []string{"Tell me a story"},
		TopP:           nil,
		N:              to.Ptr(int32(1)),
	}, mock.Anything).Return(azopenai.GetCompletionsResponse{
		Completions: azopenai.Completions{
			Choices: []azopenai.Choice{
				{Text: to.Ptr("Once upon a time.")},
			},
		},
	}, nil)
	ctx := context.Background()
	dataCtx := plugin.MapData{}
	result, diags := s.schema.ContentFunc(ctx, &plugin.ProvideContentParams{
		Args: cty.ObjectVal(map[string]cty.Value{
			"prompt":            cty.StringVal("Tell me a story"),
			"max_tokens":        cty.NumberIntVal(1000),
			"temperature":       cty.NumberFloatVal(0),
			"top_p":             cty.NilVal,
			"completions_count": cty.NumberIntVal(1),
		}),
		Config: cty.ObjectVal(map[string]cty.Value{
			"api_key":           cty.StringVal("testtoken"),
			"resource_endpoint": cty.StringVal("http://test"),
			"deployment_name":   cty.StringVal("test"),
			"api_version":       cty.StringVal("2024-02-01"),
		}),
		DataContext: dataCtx,
	})
	fmt.Println(diags)
	s.Nil(diags)
	s.Equal("Once upon a time.", mdprint.PrintString(result.Content))
}

func (s *AzureOpenAITextContentTestSuite) TestAdvanced() {
	s.cli.On("GetCompletions", mock.Anything, azopenai.CompletionsOptions{
		DeploymentName: to.Ptr("test"),
		MaxTokens:      to.Ptr(int32(1000)),
		Temperature:    to.Ptr(float32(0)),
		Prompt:         []string{"Tell me a story about BAR. {\"foo\":\"bar\"}"},
		TopP:           nil,
		N:              to.Ptr(int32(1)),
	}, mock.Anything).Return(azopenai.GetCompletionsResponse{
		Completions: azopenai.Completions{
			Choices: []azopenai.Choice{
				{Text: to.Ptr("Once upon a time.")},
			},
		},
	}, nil)

	ctx := context.Background()
	dataCtx := plugin.MapData{
		"local": plugin.MapData{
			"foo": plugin.StringData("bar"),
		},
	}

	result, diags := s.schema.ContentFunc(ctx, &plugin.ProvideContentParams{
		Args: cty.ObjectVal(map[string]cty.Value{
			"prompt":            cty.StringVal("Tell me a story about {{.local.foo | upper}}. {{ .local | toRawJson }}"),
			"max_tokens":        cty.NumberIntVal(1000),
			"temperature":       cty.NumberFloatVal(0),
			"top_p":             cty.NilVal,
			"completions_count": cty.NumberIntVal(1),
		}),
		Config: cty.ObjectVal(map[string]cty.Value{
			"api_key":           cty.StringVal("testtoken"),
			"resource_endpoint": cty.StringVal("http://test"),
			"deployment_name":   cty.StringVal("test"),
			"api_version":       cty.StringVal("2024-02-01"),
		}),
		DataContext: dataCtx,
	})
	s.Empty(diags)
	s.Equal("Once upon a time.", mdprint.PrintString(result.Content))
}

func (s *AzureOpenAITextContentTestSuite) TestMissingPrompt() {
	plugintest.DecodeAndAssert(s.T(), s.schema.Args, "", plugin.MapData{}, diagtest.Asserts{
		{
			diagtest.IsError,
			diagtest.SummaryEquals("Missing required argument"),
			diagtest.DetailContains("prompt"),
		},
		{
			diagtest.IsError,
			diagtest.SummaryEquals("Argument value must be non-null"),
			diagtest.DetailContains("prompt"),
		},
	})
}

func (s *AzureOpenAITextContentTestSuite) TestMissingAPIKey() {
	plugintest.DecodeAndAssert(s.T(), s.schema.Args, `
		prompt = "Tell me a story"
	`, plugin.MapData{}, diagtest.Asserts{})
	plugintest.DecodeAndAssert(s.T(), s.schema.Config, `
	`, plugin.MapData{}, diagtest.Asserts{
		{
			diagtest.IsError,
			diagtest.SummaryEquals("Missing required argument"),
			diagtest.DetailContains("api_key"),
		},
		{
			diagtest.IsError,
			diagtest.SummaryEquals("Argument value must be non-null"),
			diagtest.DetailContains("api_key"),
		},
		{
			diagtest.IsError,
			diagtest.SummaryEquals("Missing required argument"),
			diagtest.DetailContains("resource_endpoint"),
		},
		{
			diagtest.IsError,
			diagtest.SummaryEquals("Argument value must be non-null"),
			diagtest.DetailContains("resource_endpoint"),
		},
		{
			diagtest.IsError,
			diagtest.SummaryEquals("Missing required argument"),
			diagtest.DetailContains("deployment_name"),
		},
		{
			diagtest.IsError,
			diagtest.SummaryEquals("Argument value must be non-null"),
			diagtest.DetailContains("deployment_name"),
		},
	})
}

func (s *AzureOpenAITextContentTestSuite) TestFailingClient() {
	s.cli.On("GetCompletions", mock.Anything, azopenai.CompletionsOptions{
		DeploymentName: to.Ptr("test"),
		MaxTokens:      to.Ptr(int32(1000)),
		Temperature:    to.Ptr(float32(0)),
		Prompt:         []string{"Tell me a story"},
		TopP:           nil,
		N:              to.Ptr(int32(1)),
	}, mock.Anything).Return(azopenai.GetCompletionsResponse{}, errors.New("failed to generate text from model"))
	ctx := context.Background()
	dataCtx := plugin.MapData{}
	result, diags := s.schema.ContentFunc(ctx, &plugin.ProvideContentParams{
		Args: cty.ObjectVal(map[string]cty.Value{
			"prompt":            cty.StringVal("Tell me a story"),
			"max_tokens":        cty.NumberIntVal(1000),
			"temperature":       cty.NumberFloatVal(0),
			"top_p":             cty.NilVal,
			"completions_count": cty.NumberIntVal(1),
		}),
		Config: cty.ObjectVal(map[string]cty.Value{
			"api_key":           cty.StringVal("testtoken"),
			"resource_endpoint": cty.StringVal("http://test"),
			"deployment_name":   cty.StringVal("test"),
			"api_version":       cty.StringVal("2024-02-01"),
		}),
		DataContext: dataCtx,
	})
	s.Nil(result)
	s.Equal(diagnostics.Diag{{
		Severity: hcl.DiagError,
		Summary:  "Failed to generate text",
		Detail:   "failed to generate text from model",
	}}, diags)
}
