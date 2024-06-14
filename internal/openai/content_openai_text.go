package openai

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"strings"
	"text/template"

	"github.com/Masterminds/sprig/v3"
	"github.com/hashicorp/hcl/v2"
	"github.com/zclconf/go-cty/cty"

	"github.com/blackstork-io/fabric/internal/openai/client"
	"github.com/blackstork-io/fabric/pkg/diagnostics"
	"github.com/blackstork-io/fabric/plugin"
	"github.com/blackstork-io/fabric/plugin/dataspec"
	"github.com/blackstork-io/fabric/plugin/dataspec/constraint"
)

func makeOpenAITextContentSchema(loader ClientLoadFn) *plugin.ContentProvider {
	return &plugin.ContentProvider{
		Config: dataspec.ObjectSpec{
			&dataspec.AttrSpec{
				Name: "system_prompt",
				Type: cty.String,
			},
			&dataspec.AttrSpec{
				Name:        "api_key",
				Type:        cty.String,
				Constraints: constraint.RequiredNonNull,
				Secret:      true,
			},
			&dataspec.AttrSpec{
				Name: "organization_id",
				Type: cty.String,
			},
		},
		Args: dataspec.ObjectSpec{
			&dataspec.AttrSpec{
				Name:        "prompt",
				Type:        cty.String,
				Constraints: constraint.RequiredNonNull,
				ExampleVal:  cty.StringVal("Summarize the following text: {{.vars.text_to_summarize}}"),
			},
			&dataspec.AttrSpec{
				Name:        "model",
				Type:        cty.String,
				Constraints: constraint.Meaningful,
				DefaultVal:  cty.StringVal(defaultModel),
			},
		},
		ContentFunc: genOpenAIText(loader),
	}
}

func genOpenAIText(loader ClientLoadFn) plugin.ProvideContentFunc {
	return func(ctx context.Context, params *plugin.ProvideContentParams) (*plugin.ContentResult, diagnostics.Diag) {
		cli, err := makeClient(loader, params.Config)
		if err != nil {
			return nil, diagnostics.Diag{{
				Severity: hcl.DiagError,
				Summary:  "Failed to create client",
				Detail:   err.Error(),
			}}
		}
		result, err := renderText(ctx, cli, params.Config, params.Args, params.DataContext)
		if err != nil {
			return nil, diagnostics.Diag{{
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

func renderText(ctx context.Context, cli client.Client, cfg, args cty.Value, dataCtx plugin.MapData) (string, error) {
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
	content, err := templateText(args.GetAttr("prompt").AsString(), dataCtx)
	if err != nil {
		return "", err
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

func templateText(text string, dataCtx plugin.MapData) (string, error) {
	tmpl, err := template.New("text").Funcs(sprig.FuncMap()).Parse(text)
	if err != nil {
		return "", fmt.Errorf("failed to parse template: %w", err)
	}

	var buf bytes.Buffer
	err = tmpl.Execute(&buf, dataCtx.Any())
	if err != nil {
		return "", fmt.Errorf("failed to execute template: %w", err)
	}
	return strings.TrimSpace(buf.String()), nil
}
