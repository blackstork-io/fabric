package openai

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"html/template"
	"strings"

	"github.com/Masterminds/sprig/v3"
	"github.com/hashicorp/hcl/v2"
	"github.com/zclconf/go-cty/cty"

	"github.com/blackstork-io/fabric/internal/openai/client"
	"github.com/blackstork-io/fabric/plugin"
	"github.com/blackstork-io/fabric/plugin/dataspec"
)

func makeOpenAITextContentSchema(loader ClientLoadFn) *plugin.ContentProvider {
	return &plugin.ContentProvider{
		Config: dataspec.ObjectSpec{
			&dataspec.AttrSpec{
				Name:       "system_prompt",
				Type:       cty.String,
				Required:   false,
				ExampleVal: cty.StringVal("You are a security report summarizer"),
			},
			&dataspec.AttrSpec{
				Name:       "api_key",
				Type:       cty.String,
				Required:   true,
				ExampleVal: cty.StringVal("OPENAI_API_KEY"),
			},
			&dataspec.AttrSpec{
				Name:       "organization_id",
				Type:       cty.String,
				Required:   false,
				ExampleVal: cty.StringVal("YOUR_ORG_ID"),
			},
		},
		Args: dataspec.ObjectSpec{
			&dataspec.AttrSpec{
				Name:       "prompt",
				Type:       cty.String,
				Required:   true,
				Doc:        `Go template of the prompt for an OpenAI model`,
				ExampleVal: cty.StringVal("This is the report to be summarized: "),
			},
			&dataspec.AttrSpec{
				Name:       "model",
				Type:       cty.String,
				Required:   false,
				DefaultVal: cty.StringVal("gpt-3.5-turbo"),
			},
		},
		ContentFunc: genOpenAIText(loader),
		Doc:         `Produces a chat completion result from an OpenAI model`,
	}
}

func genOpenAIText(loader ClientLoadFn) plugin.ProvideContentFunc {
	return func(ctx context.Context, params *plugin.ProvideContentParams) (*plugin.ContentResult, hcl.Diagnostics) {
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
		return &plugin.ContentResult{
			Content: &plugin.ContentElement{
				Markdown: result,
			},
		}, nil
	}
}

func renderText(ctx context.Context, cli client.Client, cfg, args cty.Value, datactx plugin.MapData) (string, error) {
	prompt := args.GetAttr("prompt")
	if prompt.IsNull() || prompt.AsString() == "" {
		return "", errors.New("prompt is required in invocation")
	}

	params := client.ChatCompletionParams{
		Model: args.GetAttr("model").AsString(),
	}
	systemPrompt := cfg.GetAttr("system_prompt")
	if !systemPrompt.IsNull() && systemPrompt.AsString() != "" {
		params.Messages = append(params.Messages, client.ChatCompletionMessage{
			Role:    "system",
			Content: systemPrompt.AsString(),
		})
	}
	content := prompt.AsString()
	content, err := templateText(content, datactx)
	if err != nil {
		return "", err
	}
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

func templateText(text string, datactx plugin.MapData) (string, error) {
	tmpl, err := template.New("text").Funcs(sprig.FuncMap()).Parse(text)
	if err != nil {
		return "", fmt.Errorf("failed to parse template: %w", err)
	}

	var buf bytes.Buffer
	err = tmpl.Execute(&buf, datactx.Any())
	if err != nil {
		return "", fmt.Errorf("failed to execute template: %w", err)
	}
	return strings.TrimSpace(buf.String()), nil
}
