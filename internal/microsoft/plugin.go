package microsoft

import (
	"context"
	"net/url"

	"github.com/blackstork-io/fabric/internal/microsoft/client"
	"github.com/blackstork-io/fabric/plugin"
	"github.com/blackstork-io/fabric/plugin/dataspec"
	"github.com/blackstork-io/fabric/plugin/plugindata"
)

type AzureClient interface {
	QueryObjects(
		ctx context.Context,
		endpoint string,
		queryParams url.Values,
		size int,
	) (objects plugindata.List, err error)
}

type AzureClientLoadFn func(ctx context.Context, cfg *dataspec.Block) (cli AzureClient, err error)

func MakeDefaultAzureClientLoader(tokenFn client.AcquireTokenFn) AzureClientLoadFn {
	return func(ctx context.Context, cfg *dataspec.Block) (cli AzureClient, err error) {
		scopes := []string{"https://management.azure.com/.default"}
		token, err := client.AcquireTokenWithCreds(ctx, tokenFn, cfg, scopes)
		if err != nil {
			return nil, err
		}
		return client.NewAzureClient(token), nil
	}
}

type MicrosoftGraphClient interface {
	QueryObjects(
		ctx context.Context,
		endpoint string,
		queryParams url.Values,
		size int,
	) (objects plugindata.List, err error)

	QueryObject(
		ctx context.Context,
		endpoint string,
	) (object plugindata.Data, err error)
}

type MicrosoftSecurityClient interface {
	QueryObjects(
		ctx context.Context,
		endpoint string,
		queryParams url.Values,
		size int,
	) (objects plugindata.List, err error)

	QueryObject(
		ctx context.Context,
		endpoint string,
	) (object plugindata.Data, err error)
}

type MicrosoftGraphClientLoadFn func(ctx context.Context, apiVersion string, cfg *dataspec.Block) (client MicrosoftGraphClient, err error)

type MicrosoftSecurityClientLoadFn func(ctx context.Context, cfg *dataspec.Block) (client MicrosoftSecurityClient, err error)

func MakeDefaultMicrosoftGraphClientLoader(tokenFn client.AcquireTokenFn) MicrosoftGraphClientLoadFn {
	return func(ctx context.Context, apiVersion string, cfg *dataspec.Block) (cli MicrosoftGraphClient, err error) {
		scopes := []string{"https://graph.microsoft.com/.default"}
		token, err := client.AcquireTokenWithCreds(ctx, tokenFn, cfg, scopes)
		if err != nil {
			return nil, err
		}
		return client.NewGraphClient(token, apiVersion), nil
	}
}

func MakeDefaultMicrosoftSecurityClientLoader(tokenFn client.AcquireTokenFn) MicrosoftSecurityClientLoadFn {
	return func(ctx context.Context, cfg *dataspec.Block) (cli MicrosoftSecurityClient, err error) {
		scopes := []string{"https://api.securitycenter.microsoft.com/.default"}
		token, err := client.AcquireTokenWithCreds(ctx, tokenFn, cfg, scopes)
		if err != nil {
			return nil, err
		}
		return client.NewSecurityClient(token), nil
	}
}

func Plugin(
	version string,
	azureClientLoader AzureClientLoadFn,
	openAIClientLoader client.AzureOpenAIClientLoadFn,
	graphClientLoader MicrosoftGraphClientLoadFn,
	securityClientLoader MicrosoftSecurityClientLoadFn,
) *plugin.Schema {
	return &plugin.Schema{
		Doc:     "Plugin for Microsoft services.",
		Name:    "blackstork/microsoft",
		Version: version,
		DataSources: plugin.DataSources{
			"microsoft_sentinel_incidents": makeMicrosoftSentinelIncidentsDataSource(azureClientLoader),
			"microsoft_graph":              makeMicrosoftGraphDataSource(graphClientLoader),
			"microsoft_security":           makeMicrosoftSecurityDataSource(securityClientLoader),
		},
		ContentProviders: plugin.ContentProviders{
			"azure_openai_text": makeAzureOpenAITextContentSchema(openAIClientLoader),
		},
	}
}
