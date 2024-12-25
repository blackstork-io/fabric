package misp

import (
	"bytes"
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/google/uuid"
	"github.com/hashicorp/hcl/v2"
	"github.com/zclconf/go-cty/cty"
	"go.opentelemetry.io/otel/attribute"
	nooptrace "go.opentelemetry.io/otel/trace/noop"

	"github.com/blackstork-io/fabric/internal/misp/client"
	"github.com/blackstork-io/fabric/pkg/diagnostics"
	"github.com/blackstork-io/fabric/plugin"
	"github.com/blackstork-io/fabric/plugin/dataspec"
	"github.com/blackstork-io/fabric/plugin/dataspec/constraint"
	"github.com/blackstork-io/fabric/plugin/plugindata"
	"github.com/blackstork-io/fabric/print"
	"github.com/blackstork-io/fabric/print/mdprint"
)

func makeMispEventReportsPublisher(loader ClientLoaderFn) *plugin.Publisher {
	return &plugin.Publisher{
		Doc:    "Publishes content to misp event reports",
		Tags:   []string{},
		Config: makeDataSourceConfig(),
		Args: &dataspec.RootSpec{
			Attrs: []*dataspec.AttrSpec{
				{
					Name:        "event_id",
					Type:        cty.String,
					Constraints: constraint.RequiredMeaningful,
				},
				{
					Name:        "name",
					Type:        cty.String,
					Constraints: constraint.RequiredMeaningful,
				},
				{
					Name: "distribution",
					Type: cty.String,
					OneOf: []cty.Value{
						cty.StringVal("0"),
						cty.StringVal("1"),
						cty.StringVal("2"),
						cty.StringVal("3"),
						cty.StringVal("4"),
						cty.StringVal("5"),
					},
				},
				{
					Name: "sharing_group_id",
					Type: cty.String,
				},
			},
		},
		AllowedFormats: []plugin.OutputFormat{plugin.OutputFormatMD},
		PublishFunc:    publishEventReport(loader),
	}
}

func parseContent(data plugindata.Map) (document *plugin.ContentSection) {
	documentMap, ok := data["document"]
	if !ok {
		return
	}
	contentMap, ok := documentMap.(plugindata.Map)["content"]
	if !ok {
		return
	}
	content, err := plugin.ParseContentData(contentMap.(plugindata.Map))
	if err != nil {
		return
	}
	document = content.(*plugin.ContentSection)
	return
}

func publishEventReport(loader ClientLoaderFn) plugin.PublishFunc {
	logger := slog.Default()
	tracer := nooptrace.Tracer{}
	return func(ctx context.Context, params *plugin.PublishParams) diagnostics.Diag {
		document := parseContent(params.DataContext)
		if document == nil {
			return diagnostics.Diag{{
				Severity: hcl.DiagError,
				Summary:  "Failed to parse document",
				Detail:   "document is required",
			}}
		}
		datactx := params.DataContext
		datactx["format"] = plugindata.String(params.Format.String())
		var printer print.Printer
		switch params.Format {
		case plugin.OutputFormatMD:
			printer = mdprint.New()
		default:
			return diagnostics.Diag{{
				Severity: hcl.DiagError,
				Summary:  "Unsupported format",
				Detail:   "Only md format is supported",
			}}
		}
		printer = print.WithLogging(printer, logger, slog.String("format", params.Format.String()))
		printer = print.WithTracing(printer, tracer, attribute.String("format", params.Format.String()))

		buff := bytes.NewBuffer(nil)
		err := printer.Print(ctx, buff, document)
		if err != nil {
			return diagnostics.Diag{{
				Severity: hcl.DiagError,
				Summary:  "Failed to render the report",
				Detail:   err.Error(),
			}}
		}

		cli := loader(params.Config)

		timestamp := fmt.Sprintf("%d", time.Now().Unix())
		report := client.AddEventReportRequest{
			Uuid:      uuid.New().String(),
			EventId:   params.Args.GetAttrVal("event_id").AsString(),
			Name:      params.Args.GetAttrVal("name").AsString(),
			Content:   buff.String(),
			Timestamp: &timestamp,
			Deleted:   false,
		}
		distribution := params.Args.GetAttrVal("distribution")
		if !distribution.IsNull() {
			distributionStr := distribution.AsString()
			report.Distribution = &distributionStr
		}
		sharingGroupId := params.Args.GetAttrVal("sharing_group_id")
		if !sharingGroupId.IsNull() {
			sharingGroupIdStr := sharingGroupId.AsString()
			report.SharingGroupId = &sharingGroupIdStr
		}

		slog.InfoContext(ctx, "Publish to misp event reports", "filename", report.Name)

		resp, err := cli.AddEventReport(ctx, report)
		if err != nil {
			return diagnostics.Diag{{
				Severity: hcl.DiagError,
				Summary:  "Failed to add event report",
				Detail:   err.Error(),
			}}
		}

		slog.InfoContext(ctx, "Successfully added report", "id", resp.EventReport.Id, "uuid", resp.EventReport.Uuid, "event_id", resp.EventReport.EventId, "name", resp.EventReport.Name)
		return nil
	}
}
