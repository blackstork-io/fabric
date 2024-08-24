package notion

import (
	"bytes"
	"context"
	"io"
	"log/slog"

	"github.com/blackstork-io/fabric/internal/builtin"
	"github.com/blackstork-io/fabric/pkg/diagnostics"
	"github.com/blackstork-io/fabric/plugin"
	"github.com/blackstork-io/fabric/plugin/dataspec"
	"github.com/blackstork-io/fabric/plugin/dataspec/constraint"
	"github.com/blackstork-io/fabric/plugin/plugindata"
	"github.com/blackstork-io/fabric/print/mdprint"
	"github.com/brittonhayes/notionmd"
	"github.com/dstotijn/go-notion"
	"github.com/hashicorp/hcl/v2"
	"github.com/zclconf/go-cty/cty"
	"go.opentelemetry.io/otel/trace"
	nooptrace "go.opentelemetry.io/otel/trace/noop"
)

func makeNotionPagePublisher(logger *slog.Logger, tracer trace.Tracer) *plugin.Publisher {
	if logger == nil {
		logger = slog.New(slog.NewTextHandler(io.Discard, nil))
	}
	if tracer == nil {
		tracer = nooptrace.Tracer{}
	}

	return &plugin.Publisher{
		Doc:  "Publishes content to a Notion page",
		Tags: []string{},
		Args: &dataspec.RootSpec{
			Attrs: []*dataspec.AttrSpec{
				{
					Name:        "title",
					Doc:         "Title of the Notion page",
					Type:        cty.String,
					ExampleVal:  cty.StringVal("My Notion Page"),
					Constraints: constraint.Required,
				},
				{
					Name:        "parent_page_id",
					Doc:         "Notion parent page ID",
					Type:        cty.String,
					ExampleVal:  cty.StringVal("1234567890"),
					Constraints: constraint.Required,
				},
				{
					Name:        "api_key",
					Doc:         "Notion API key",
					Type:        cty.String,
					ExampleVal:  cty.StringVal("secret_1234567890"),
					Constraints: constraint.Required,
					Secret:      true,
				},
			},
		},
		AllowedFormats: []plugin.OutputFormat{plugin.OutputFormatMD},
		PublishFunc:    publishNotionPage(logger, tracer),
	}
}

func publishNotionPage(logger *slog.Logger, _ trace.Tracer) plugin.PublishFunc {
	return func(ctx context.Context, params *plugin.PublishParams) diagnostics.Diag {
		document, _ := builtin.ParseScope(params.DataContext)
		if document == nil {
			return diagnostics.Diag{{
				Severity: hcl.DiagError,
				Summary:  "Failed to parse document",
				Detail:   "document is required",
			}}
		}

		datactx := params.DataContext
		datactx["format"] = plugindata.String(params.Format.String())

		titleAttr := params.Args.GetAttrVal("title")
		if titleAttr.IsNull() || titleAttr.AsString() == "" {
			return diagnostics.Diag{{
				Severity: hcl.DiagError,
				Summary:  "Failed to parse arguments",
				Detail:   "title is required",
			}}
		}

		parentPageIDAttr := params.Args.GetAttrVal("parent_page_id")
		if parentPageIDAttr.IsNull() || parentPageIDAttr.AsString() == "" {
			return diagnostics.Diag{{
				Severity: hcl.DiagError,
				Summary:  "Failed to parse arguments",
				Detail:   "parent_page_id is required",
			}}
		}

		apiKeyAttr := params.Args.GetAttrVal("api_key")
		if apiKeyAttr.IsNull() || apiKeyAttr.AsString() == "" {
			return diagnostics.Diag{{
				Severity: hcl.DiagError,
				Summary:  "Failed to parse arguments",
				Detail:   "api_key is required",
			}}
		}

		writer := bytes.NewBuffer([]byte{})

		printer := mdprint.New()
		err := printer.Print(ctx, writer, document)
		if err != nil {
			return diagnostics.Diag{{
				Severity: hcl.DiagError,
				Summary:  "Failed to print content",
				Detail:   err.Error(),
			}}
		}

		blocks, err := notionmd.Convert(writer.String())
		if err != nil {
			return diagnostics.Diag{{
				Severity: hcl.DiagError,
				Summary:  "Failed to convert content to Notion blocks",
				Detail:   err.Error(),
			}}
		}

		// Publish to Notion
		logger.InfoContext(ctx, "Publishing to Notion", "title", titleAttr.AsString())
		client := notion.NewClient(apiKeyAttr.AsString())
		page, err := client.CreatePage(ctx, notion.CreatePageParams{
			ParentType: notion.ParentTypePage,
			ParentID:   parentPageIDAttr.AsString(),
			Title: []notion.RichText{
				{
					Type: notion.RichTextTypeText,
					Text: &notion.Text{
						Content: titleAttr.AsString(),
					},
				},
			},
			Children: blocks,
		})
		if err != nil {
			return diagnostics.Diag{{
				Severity: hcl.DiagError,
				Summary:  "Failed to create a Notion page",
				Detail:   err.Error(),
			}}
		}

		logger.InfoContext(ctx, "Published to Notion", "page_id", page.ID)

		return nil
	}
}
