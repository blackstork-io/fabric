package builtin

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"log/slog"
	"os"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/Masterminds/sprig/v3"
	"github.com/hashicorp/hcl/v2"
	"github.com/zclconf/go-cty/cty"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
	nooptrace "go.opentelemetry.io/otel/trace/noop"

	"github.com/blackstork-io/fabric/pkg/diagnostics"
	"github.com/blackstork-io/fabric/plugin"
	"github.com/blackstork-io/fabric/plugin/dataspec"
	"github.com/blackstork-io/fabric/plugin/dataspec/constraint"
	"github.com/blackstork-io/fabric/plugin/plugindata"
	"github.com/blackstork-io/fabric/print"
	"github.com/blackstork-io/fabric/print/htmlprint"
	"github.com/blackstork-io/fabric/print/mdprint"
	"github.com/blackstork-io/fabric/print/pdfprint"
)

func makeLocalFilePublisher(logger *slog.Logger, tracer trace.Tracer) *plugin.Publisher {
	if logger == nil {
		logger = slog.New(slog.NewTextHandler(io.Discard, nil))
	}
	if tracer == nil {
		tracer = nooptrace.Tracer{}
	}
	return &plugin.Publisher{
		Doc:  "Publishes content to local file",
		Tags: []string{},
		Args: &dataspec.RootSpec{
			Attrs: []*dataspec.AttrSpec{
				{
					Name:        "path",
					Doc:         "Path to the file",
					Type:        cty.String,
					ExampleVal:  cty.StringVal("dist/output.md"),
					Constraints: constraint.RequiredMeaningful,
				},
				{
					Name:       "format",
					Doc:        "Format of the file. If not provided, the format will be inferred from the file extension",
					Type:       cty.String,
					ExampleVal: cty.StringVal("md"),
					OneOf:      []cty.Value{cty.StringVal("md"), cty.StringVal("html"), cty.StringVal("pdf")},
				},
			},
		},
		PublishFunc: publishLocalFile(logger, tracer),
	}
}

func publishLocalFile(logger *slog.Logger, tracer trace.Tracer) plugin.PublishFunc {
	return func(ctx context.Context, params *plugin.PublishParams) diagnostics.Diag {
		document, _ := parseScope(params.DataContext)
		if document == nil {
			return diagnostics.Diag{{
				Severity: hcl.DiagError,
				Summary:  "Failed to parse document",
				Detail:   "document is required",
			}}
		}
		datactx := params.DataContext

		path, err := templatePath(params.Args.GetAttrVal("path").AsString(), datactx)
		if err != nil {
			return diagnostics.Diag{{
				Severity: hcl.DiagError,
				Summary:  "Failed to render a path value",
				Detail:   err.Error(),
			}}
		}
		var format string
		formatAttr := params.Args.GetAttrVal("format")
		if formatAttr.IsNull() {
			format = strings.ToLower(strings.TrimLeft(filepath.Ext(path), "."))
		} else {
			format = formatAttr.AsString()
		}
		datactx["format"] = plugindata.String(format)

		var printer print.Printer
		switch format {
		case "md":
			printer = mdprint.New()
		case "html":
			printer = htmlprint.New()
		case "pdf":
			printer = pdfprint.New()
		default:
			return diagnostics.Diag{{
				Severity: hcl.DiagError,
				Summary:  "Unsupported format",
				Detail:   "Only md, html and pdf formats are supported",
			}}
		}
		printer = print.WithLogging(printer, logger, slog.String("format", format))
		printer = print.WithTracing(printer, tracer, attribute.String("format", format))
		logger.InfoContext(ctx, "Writing to a file", "path", path)
		dir := filepath.Dir(path)
		err = os.MkdirAll(dir, 0o755)
		if err != nil {
			return diagnostics.Diag{{
				Severity: hcl.DiagError,
				Summary:  "Failed to create a directory",
				Detail:   err.Error(),
			}}
		}
		fs, err := os.Create(path)
		if err != nil {
			return diagnostics.Diag{{
				Severity: hcl.DiagError,
				Summary:  "Failed to create a file",
				Detail:   err.Error(),
			}}
		}
		defer fs.Close()
		err = printer.Print(ctx, fs, document)
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

func templatePath(pattern string, datactx plugindata.Map) (string, error) {
	tmpl, err := template.New("pattern").Funcs(sprig.FuncMap()).Parse(pattern)
	if err != nil {
		return "", fmt.Errorf("failed to parse a text template: %w", err)
	}

	var buf bytes.Buffer
	err = tmpl.Execute(&buf, datactx.Any())
	if err != nil {
		return "", fmt.Errorf("failed to execute a text template: %w", err)
	}
	return filepath.Abs(filepath.Clean(strings.TrimSpace(buf.String())))
}
