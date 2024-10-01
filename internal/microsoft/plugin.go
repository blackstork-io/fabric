package microsoft

import (
	"context"
	"fmt"
	"net/url"
	"os"

	"github.com/Azure/azure-sdk-for-go/sdk/ai/azopenai"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/AzureAD/microsoft-authentication-library-for-go/apps/confidential"

	"github.com/blackstork-io/fabric/internal/microsoft/client"
	"github.com/blackstork-io/fabric/plugin"
	"github.com/blackstork-io/fabric/plugin/dataspec"
	"github.com/blackstork-io/fabric/plugin/plugindata"
)

type ClientLoadFn func() client.Client

var DefaultClientLoader ClientLoadFn = client.New

type AzureOpenaiClientLoadFn func(azureOpenAIKey string, azureOpenAIEndpoint string) (client AzureOpenaiClient, err error)

type AzureOpenaiClient interface {
	GetCompletions(
		ctx context.Context,
		body azopenai.CompletionsOptions,
		options *azopenai.GetCompletionsOptions,
	) (azopenai.GetCompletionsResponse, error)
}

var DefaultAzureOpenAIClientLoader AzureOpenaiClientLoadFn = func(azureOpenAIKey string, azureOpenAIEndpoint string) (client AzureOpenaiClient, err error) {
	keyCredential := azcore.NewKeyCredential(azureOpenAIKey)
	client, err = azopenai.NewClientWithKeyCredential(azureOpenAIEndpoint, keyCredential, nil)
	return
}

type MicrosoftGraphClient interface {
	QueryGraph(
		ctx context.Context,
		endpoint string,
		queryParams url.Values,
		size int,
		onlyObjects bool,
	) (objects plugindata.Data, err error)
	QueryGraphObject(
		ctx context.Context,
		endpoint string,
	) (result plugindata.Data, err error)
}

type AcquireTokenFn func(ctx context.Context, tenantId string, clientId string, cred confidential.Credential) (string, error)

type MicrosoftGraphClientLoadFn func(ctx context.Context, apiVersion string, cfg *dataspec.Block) (client MicrosoftGraphClient, err error)

func MakeDefaultMicrosoftGraphClientLoader(tokenFn AcquireTokenFn) MicrosoftGraphClientLoadFn {
	return func(ctx context.Context, apiVersion string, cfg *dataspec.Block) (cli MicrosoftGraphClient, err error) {
		if cfg == nil {
			return nil, fmt.Errorf("configuration is required")
		}
		tenantId := cfg.GetAttrVal("tenant_id").AsString()
		clientId := cfg.GetAttrVal("client_id").AsString()
		clientSecretAttr := cfg.GetAttrVal("client_secret")
		if !clientSecretAttr.IsNull() {
			cred, err := confidential.NewCredFromSecret(clientSecretAttr.AsString())
			if err != nil {
				return nil, err
			}
			accessToken, err := tokenFn(ctx, tenantId, clientId, cred)
			if err != nil {
				return nil, err
			}
			return client.NewGraphClient(accessToken, apiVersion), nil
		}

		// if client_secret is not provided, try to use private_key
		privateKeyFileAttr := cfg.GetAttrVal("private_key_file")
		privateKeyAttr := cfg.GetAttrVal("private_key")

		if !privateKeyFileAttr.IsNull() || !privateKeyAttr.IsNull() {
			var pemData []byte
			if !privateKeyAttr.IsNull() {
				pemData = []byte(privateKeyAttr.AsString())
			} else {
				pemData, err = os.ReadFile(privateKeyFileAttr.AsString())
				if err != nil {
					return nil, fmt.Errorf("failed to read private key file: %w", err)
				}
			}

			keyPassphrase := ""
			keyPassphraseAttr := cfg.GetAttrVal("key_passphrase")
			if !keyPassphraseAttr.IsNull() {
				keyPassphrase = keyPassphraseAttr.AsString()
			}

			certs, privateKey, err := confidential.CertFromPEM(pemData, keyPassphrase)
			if err != nil {
				return nil, fmt.Errorf("failed to parse private key: %w", err)
			}
			cred, err := confidential.NewCredFromCert(certs, privateKey)
			if err != nil {
				return nil, fmt.Errorf("failed to create credential from cert: %w", err)
			}
			accessToken, err := tokenFn(ctx, tenantId, clientId, cred)
			if err != nil {
				return nil, err
			}
			return client.NewGraphClient(accessToken, apiVersion), nil
		}

		return nil, fmt.Errorf("missing credentials to authenticate. client_secret or private_key is required")
	}
}

func Plugin(
	version string,
	loader ClientLoadFn,
	openAiClientLoader AzureOpenaiClientLoadFn,
	graphClientLoader MicrosoftGraphClientLoadFn,
) *plugin.Schema {
	if loader == nil {
		loader = DefaultClientLoader
	}
	return &plugin.Schema{
		Doc:     "The `microsoft` plugin for Microsoft services.",
		Name:    "blackstork/microsoft",
		Version: version,
		DataSources: plugin.DataSources{
			"microsoft_sentinel_incidents": makeMicrosoftSentinelIncidentsDataSource(loader),
			"microsoft_graph":              makeMicrosoftGraphDataSource(graphClientLoader),
		},
		ContentProviders: plugin.ContentProviders{
			"azure_openai_text": makeAzureOpenAITextContentSchema(openAiClientLoader),
		},
	}
}

func makeClient(ctx context.Context, loader ClientLoadFn, cfg *dataspec.Block) (client.Client, error) {
	if cfg == nil {
		return nil, fmt.Errorf("configuration is required")
	}
	cli := loader()
	res, err := cli.GetClientCredentialsToken(ctx, &client.GetClientCredentialsTokenReq{
		TenantID:     cfg.GetAttrVal("tenant_id").AsString(),
		ClientID:     cfg.GetAttrVal("client_id").AsString(),
		ClientSecret: cfg.GetAttrVal("client_secret").AsString(),
	})
	if err != nil {
		return nil, err
	}
	cli.UseAuth(res.AccessToken)
	return cli, nil
}
