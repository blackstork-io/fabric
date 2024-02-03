package openai

import (
	"context"
	"encoding/json"
	"errors"

	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/hcldec"
	"github.com/zclconf/go-cty/cty"

	"github.com/blackstork-io/fabric/internal/openai/client"
	"github.com/blackstork-io/fabric/plugin"
)

func makeOpenAITextContentSchema(loader ClientLoadFn) *plugin.ContentProvider {
	return &plugin.ContentProvider{
		Config: &hcldec.ObjectSpec{
			"system_prompt": &hcldec.AttrSpec{
				Name:     "system_prompt",
				Type:     cty.String,
				Required: false,
			},
			"api_key": &hcldec.AttrSpec{
				Name:     "api_key",
				Type:     cty.String,
				Required: true,
			},
			"organization_id": &hcldec.AttrSpec{
				Name:     "organization_id",
				Type:     cty.String,
				Required: false,
			},
		},
		Args: &hcldec.ObjectSpec{
			"prompt": &hcldec.AttrSpec{
				Name:     "prompt",
				Type:     cty.String,
				Required: true,
			},
			"model": &hcldec.AttrSpec{
				Name:     "model",
				Type:     cty.String,
				Required: false,
			},
		},
		ContentFunc: genOpenAIText(loader),
	}
}

func genOpenAIText(loader ClientLoadFn) plugin.ProvideContentFunc {
	return func(ctx context.Context, params *plugin.ProvideContentParams) (*plugin.Content, hcl.Diagnostics) {
		cli, err := makeClient(loader, params.Config)
		if err != nil {
			return nil, hcl.Diagnostics{{
				Severity: hcl.DiagError,
				Summary:  "Failed to create client",
				Detail:   err.Error(),
			}}
		}
		result, err := renderText(ctx, cli, params.Config, params.Args, params.DataContext)
		if err != nil {
			return nil, hcl.Diagnostics{{
				Severity: hcl.DiagError,
				Summary:  "Failed to generate text",
				Detail:   err.Error(),
			}}
		}
		return &plugin.Content{
			Markdown: result,
		}, nil
	}
}

func renderText(ctx context.Context, cli client.Client, cfg, args cty.Value, datactx plugin.MapData) (string, error) {
	prompt := args.GetAttr("prompt")
	if prompt.IsNull() || prompt.AsString() == "" {
		return "", errors.New("prompt is required in invocation")
	}
	model := args.GetAttr("model")
	if model.IsNull() || model.AsString() == "" {
		model = cty.StringVal(defaultModel)
	}
	params := client.ChatCompletionParams{
		Model: model.AsString(),
	}
	systemPrompt := cfg.GetAttr("system_prompt")
	if !systemPrompt.IsNull() && systemPrompt.AsString() != "" {
		params.Messages = append(params.Messages, client.ChatCompletionMessage{
			Role:    "system",
			Content: systemPrompt.AsString(),
		})
	}
	content := prompt.AsString()
	if datactx != nil {
		if queryResult, ok := datactx[queryResultKey]; ok {
			data, err := json.MarshalIndent(queryResult, "", "  ")
			if err != nil {
				return "", err
			}
			content += "\n```\n" + string(data) + "\n```"
		}
	}
	params.Messages = append(params.Messages, client.ChatCompletionMessage{
		Role:    "user",
		Content: content,
	})
	result, err := cli.GenerateChatCompletion(ctx, &params)
	if err != nil {
		return "", err
	}
	if len(result.Choices) < 1 {
		return "", errors.New("no choices")
	}
	return result.Choices[0].Message.Content, nil
}
