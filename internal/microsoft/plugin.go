package microsoft

import (
	"context"
	"fmt"

	"github.com/Azure/azure-sdk-for-go/sdk/ai/azopenai"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/zclconf/go-cty/cty"

	"github.com/blackstork-io/fabric/internal/microsoft/client"
	"github.com/blackstork-io/fabric/plugin"
)

type ClientLoadFn func() client.Client

var DefaultClientLoader ClientLoadFn = client.New

type AzureOpenaiClientLoadFn func(azureOpenAIKey string, azureOpenAIEndpoint string) (client AzureOpenaiClient, err error)

type AzureOpenaiClient interface {
	GetCompletions(ctx context.Context, body azopenai.CompletionsOptions, options *azopenai.GetCompletionsOptions) (azopenai.GetCompletionsResponse, error)
}

var DefaultAzureOpenAIClientLoader AzureOpenaiClientLoadFn = func(azureOpenAIKey string, azureOpenAIEndpoint string) (client AzureOpenaiClient, err error) {
	keyCredential := azcore.NewKeyCredential(azureOpenAIKey)
	client, err = azopenai.NewClientWithKeyCredential(azureOpenAIEndpoint, keyCredential, nil)
	return
}

func Plugin(version string, loader ClientLoadFn, openAiClientLoader AzureOpenaiClientLoadFn) *plugin.Schema {
	if loader == nil {
		loader = DefaultClientLoader
	}
	return &plugin.Schema{
		Doc:     "The `microsoft` plugin for Microsoft services.",
		Name:    "blackstork/microsoft",
		Version: version,
		DataSources: plugin.DataSources{
			"microsoft_sentinel_incidents": makeMicrosoftSentinelIncidentsDataSource(loader),
		},
		ContentProviders: plugin.ContentProviders{
			"azure_openai_text": makeAzureOpenAITextContentSchema(openAiClientLoader),
		},
	}
}

func makeClient(ctx context.Context, loader ClientLoadFn, cfg cty.Value) (client.Client, error) {
	if cfg.IsNull() {
		return nil, fmt.Errorf("configuration is required")
	}
	cli := loader()
	res, err := cli.GetClientCredentialsToken(ctx, &client.GetClientCredentialsTokenReq{
		TenantID:     cfg.GetAttr("tenant_id").AsString(),
		ClientID:     cfg.GetAttr("client_id").AsString(),
		ClientSecret: cfg.GetAttr("client_secret").AsString(),
	})
	if err != nil {
		return nil, err
	}
	cli.UseAuth(res.AccessToken)
	return cli, nil
}
