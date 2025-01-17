package builtin

import (
	"context"
	"fmt"
	"io"
	"log/slog"
	"os"

	"github.com/hashicorp/hcl/v2"
	"github.com/zclconf/go-cty/cty"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
	nooptrace "go.opentelemetry.io/otel/trace/noop"

	"github.com/blackstork-io/fabric/pkg/diagnostics"
	"github.com/blackstork-io/fabric/plugin"
	"github.com/blackstork-io/fabric/plugin/dataspec"
	"github.com/blackstork-io/fabric/plugin/plugindata"
	"github.com/blackstork-io/fabric/print"
	"github.com/blackstork-io/fabric/print/htmlprint"
	"github.com/blackstork-io/fabric/print/mdprint"
)

func makeStdoutPublisher(logger *slog.Logger, tracer trace.Tracer) *plugin.Publisher {
	if logger == nil {
		logger = slog.New(slog.NewTextHandler(io.Discard, nil))
	}
	if tracer == nil {
		tracer = nooptrace.Tracer{}
	}
	return &plugin.Publisher{
		Doc:  "Prints content to stdout",
		Tags: []string{},
		Args: &dataspec.RootSpec{
			Attrs: []*dataspec.AttrSpec{
				{
					Name:       "format",
					Doc:        "Format of the output",
					Type:       cty.String,
					DefaultVal: cty.StringVal("md"),
					OneOf:      []cty.Value{cty.StringVal("md"), cty.StringVal("html")},
				},
			},
		},
		PublishFunc: publishStdout(logger, tracer),
	}
}

func publishStdout(logger *slog.Logger, tracer trace.Tracer) plugin.PublishFunc {
	return func(ctx context.Context, params *plugin.PublishParams) diagnostics.Diag {
		if params.Document == nil {
			return diagnostics.Diag{{
				Severity: hcl.DiagError,
				Summary:  "Failed to parse document",
				Detail:   "document is required",
			}}
		}
		format := params.Args.GetAttrVal("format").AsString()

		var printer print.Printer
		switch format {
		case "md":
			printer = mdprint.New()
		case "html":
			printer = htmlprint.New()
		default:
			return diagnostics.Diag{{
				Severity: hcl.DiagError,
				Summary:  "Unsupported format",
				Detail:   "Only md and html formats are supported",
			}}
		}

		datactx := params.DataContext
		datactx["format"] = plugindata.String(format)

		printer = print.WithLogging(printer, logger, slog.String("format", format))
		printer = print.WithTracing(printer, tracer, attribute.String("format", format))

		logger.InfoContext(ctx, "Printing to stdout")

		err := printer.Print(ctx, os.Stdout, params.Document)
		fmt.Println()
		if err != nil {
			return diagnostics.Diag{{
				Severity: hcl.DiagError,
				Summary:  "Failed to write to a file",
				Detail:   err.Error(),
			}}
		}
		return nil
	}
}
