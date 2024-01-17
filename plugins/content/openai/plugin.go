package openai

import (
	"context"
	"encoding/json"
	"errors"

	"github.com/Masterminds/semver/v3"
	"github.com/blackstork-io/fabric/plugininterface/v1"
	"github.com/blackstork-io/fabric/plugins/content/openai/client"
	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/hcldec"
	"github.com/zclconf/go-cty/cty"
)

var Version = semver.MustParse("0.1.0")

const (
	defaultModel   = "gpt-3.5-turbo"
	queryResultKey = "query_result"
)

type ClientLoadFn func(opts ...client.Option) client.Client

var DefaultClientLoader ClientLoadFn = func(opts ...client.Option) client.Client {
	return client.New(opts...)
}

type Plugin struct {
	ClientLoader ClientLoadFn
}

func (Plugin) GetPlugins() []plugininterface.Plugin {
	return []plugininterface.Plugin{
		{
			Namespace: "blackstork",
			Kind:      "content",
			Name:      "openai_text",
			Version:   plugininterface.Version(*Version),
			ConfigSpec: &hcldec.ObjectSpec{
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
			InvocationSpec: &hcldec.ObjectSpec{
				"prompt": &hcldec.AttrSpec{
					Name:     "prompt",
					Type:     cty.String,
					Required: true,
				},
			},
		},
	}
}

func (p Plugin) makeClient(cfg cty.Value) (client.Client, error) {
	opts := []client.Option{}
	apiKey := cfg.GetAttr("api_key")
	if apiKey.IsNull() || apiKey.AsString() == "" {
		return nil, errors.New("api_key is required in configuration")
	}
	opts = append(opts, client.WithAPIKey(apiKey.AsString()))
	orgID := cfg.GetAttr("organization_id")
	if !orgID.IsNull() && orgID.AsString() != "" {
		opts = append(opts, client.WithOrgID(orgID.AsString()))
	}
	return p.ClientLoader(opts...), nil
}

func (p Plugin) generate(cli client.Client, cfg cty.Value, args cty.Value, datactx map[string]any) (string, error) {
	prompt := args.GetAttr("prompt")
	if prompt.IsNull() || prompt.AsString() == "" {
		return "", errors.New("prompt is required in invocation")
	}
	params := client.ChatCompletionParams{
		Model: defaultModel,
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
	result, err := cli.GenerateChatCompletion(context.Background(), &params)
	if err != nil {
		return "", err
	}
	if len(result.Choices) < 1 {
		return "", errors.New("no choices")
	}
	return result.Choices[0].Message.Content, nil
}

func (p Plugin) Call(args plugininterface.Args) plugininterface.Result {
	cli, err := p.makeClient(args.Config)
	if err != nil {
		return plugininterface.Result{
			Diags: hcl.Diagnostics{{
				Severity: hcl.DiagError,
				Summary:  "Failed to create client",
				Detail:   err.Error(),
			}},
		}
	}
	result, err := p.generate(cli, args.Config, args.Args, args.Context)
	if err != nil {
		return plugininterface.Result{
			Diags: hcl.Diagnostics{{
				Severity: hcl.DiagError,
				Summary:  "Failed to generate text",
				Detail:   err.Error(),
			}},
		}
	}
	return plugininterface.Result{
		Result: result,
	}
}
