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
					Constraints: constraint.Required,
				},
			},
		},
		AllowedFormats: []plugin.OutputFormat{plugin.OutputFormatMD, plugin.OutputFormatHTML, plugin.OutputFormatPDF},
		PublishFunc:    publishLocalFile(logger, tracer),
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
		datactx["format"] = plugindata.String(params.Format.String())

		var printer print.Printer
		switch params.Format {
		case plugin.OutputFormatMD:
			printer = mdprint.New()
		case plugin.OutputFormatHTML:
			printer = htmlprint.New()
		case plugin.OutputFormatPDF:
			printer = pdfprint.New()
		default:
			return diagnostics.Diag{{
				Severity: hcl.DiagError,
				Summary:  "Unsupported format",
				Detail:   "Only md, html and pdf formats are supported",
			}}
		}
		printer = print.WithLogging(printer, logger, slog.String("format", params.Format.String()))
		printer = print.WithTracing(printer, tracer, attribute.String("format", params.Format.String()))
		pathAttr := params.Args.GetAttrVal("path")
		if pathAttr.IsNull() || pathAttr.AsString() == "" {
			return diagnostics.Diag{{
				Severity: hcl.DiagError,
				Summary:  "Failed to parse arguments",
				Detail:   "path is required",
			}}
		}
		path, err := templatePath(pathAttr.AsString(), datactx)
		if err != nil {
			return diagnostics.Diag{{
				Severity: hcl.DiagError,
				Summary:  "Failed to render a path value",
				Detail:   err.Error(),
			}}
		}
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
