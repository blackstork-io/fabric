package client

import (
	"context"

	"github.com/Azure/azure-sdk-for-go/sdk/ai/azopenai"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
)

type AzureOpenAIClientLoadFn func(azureOpenAIKey string, azureOpenAIEndpoint string) (client AzureOpenAIClient, err error)

type AzureOpenAIClient interface {
	GetCompletions(
		ctx context.Context,
		body azopenai.CompletionsOptions,
		options *azopenai.GetCompletionsOptions,
	) (azopenai.GetCompletionsResponse, error)
}

func MakeAzureOpenAIClientLoader() AzureOpenAIClientLoadFn {
	return func(azureOpenAIKey string, azureOpenAIEndpoint string) (client AzureOpenAIClient, err error) {
		keyCredential := azcore.NewKeyCredential(azureOpenAIKey)
		client, err = azopenai.NewClientWithKeyCredential(azureOpenAIEndpoint, keyCredential, nil)
		return
	}
}
