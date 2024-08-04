package microsoft

import (
	"bytes"
	"context"
	"fmt"
	"strings"
	"text/template"

	"github.com/Azure/azure-sdk-for-go/sdk/ai/azopenai"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	"github.com/Masterminds/sprig/v3"
	"github.com/blackstork-io/fabric/pkg/diagnostics"
	"github.com/blackstork-io/fabric/plugin"
	"github.com/blackstork-io/fabric/plugin/dataspec"
	"github.com/blackstork-io/fabric/plugin/dataspec/constraint"
	"github.com/hashicorp/hcl/v2"
	"github.com/zclconf/go-cty/cty"
)

func makeAzureOpenAITextContentSchema(loader AzureOpenAIClientLoadFn) *plugin.ContentProvider {
	return &plugin.ContentProvider{
		Config: dataspec.ObjectSpec{
			&dataspec.AttrSpec{
				Name:        "api_key",
				Type:        cty.String,
				Constraints: constraint.RequiredNonNull,
				Secret:      true,
			},
			&dataspec.AttrSpec{
				Name:        "resource_endpoint",
				Type:        cty.String,
				Constraints: constraint.RequiredNonNull,
			},
			&dataspec.AttrSpec{
				Name:        "deployment_name",
				Type:        cty.String,
				Constraints: constraint.RequiredNonNull,
			},
			&dataspec.AttrSpec{
				Name:       "api_version",
				Type:       cty.String,
				DefaultVal: cty.StringVal("2024-02-01"),
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
				Name:       "max_tokens",
				Type:       cty.Number,
				DefaultVal: cty.NumberIntVal(1000),
			},
			&dataspec.AttrSpec{
				Name:       "temperature",
				Type:       cty.Number,
				DefaultVal: cty.NumberFloatVal(0),
			},
			&dataspec.AttrSpec{
				Name: "top_p",
				Type: cty.Number,
			},
			&dataspec.AttrSpec{
				Name:       "completions_count",
				Type:       cty.Number,
				DefaultVal: cty.NumberIntVal(1),
			},
		},
		ContentFunc: genOpenAIText(loader),
	}
}

func genOpenAIText(loader AzureOpenAIClientLoadFn) plugin.ProvideContentFunc {
	return func(ctx context.Context, params *plugin.ProvideContentParams) (*plugin.ContentResult, diagnostics.Diag) {
		apiKey := params.Config.GetAttr("api_key").AsString()
		resourceEndpoint := params.Config.GetAttr("resource_endpoint").AsString()
		client, err := loader(apiKey, resourceEndpoint)
		if err != nil {
			return nil, diagnostics.Diag{{
				Severity: hcl.DiagError,
				Summary:  "Failed to create client",
				Detail:   err.Error(),
			}}
		}
		result, err := renderText(ctx, client, params.Config, params.Args, params.DataContext)
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

func renderText(ctx context.Context, cli AzureOpenAIClient, cfg, args cty.Value, dataCtx plugin.MapData) (string, error) {

	params := azopenai.CompletionsOptions{}
	params.DeploymentName = to.Ptr(cfg.GetAttr("deployment_name").AsString())

	maxTokens, _ := args.GetAttr("max_tokens").AsBigFloat().Int64()
	params.MaxTokens = to.Ptr(int32(maxTokens))

	temperature, _ := args.GetAttr("temperature").AsBigFloat().Float32()
	params.Temperature = to.Ptr(temperature)

	completionsCount, _ := args.GetAttr("completions_count").AsBigFloat().Int64()
	params.N = to.Ptr(int32(completionsCount))

	topPAttr := args.GetAttr("top_p")
	if !topPAttr.IsNull() {
		topP, _ := topPAttr.AsBigFloat().Float32()
		params.TopP = to.Ptr(topP)
	}

	renderedPrompt, err := templateText(args.GetAttr("prompt").AsString(), dataCtx)
	if err != nil {
		return "", err
	}
	params.Prompt = []string{renderedPrompt}
	// TODO: use api version from config
	resp, err := cli.GetCompletions(ctx, params, nil)
	if err != nil {
		return "", err
	}
	if len(resp.Choices) == 0 {
		return "", nil
	}
	return *resp.Choices[0].Text, nil
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
