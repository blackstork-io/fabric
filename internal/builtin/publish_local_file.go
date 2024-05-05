package builtin

import (
	"bytes"
	"context"
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/Masterminds/sprig/v3"
	"github.com/hashicorp/hcl/v2"
	"github.com/zclconf/go-cty/cty"

	"github.com/blackstork-io/fabric/plugin"
	"github.com/blackstork-io/fabric/plugin/dataspec"
	"github.com/blackstork-io/fabric/plugin/dataspec/constraint"
	"github.com/blackstork-io/fabric/printer"
	"github.com/blackstork-io/fabric/printer/htmlprint"
	"github.com/blackstork-io/fabric/printer/mdprint"
	"github.com/blackstork-io/fabric/printer/pdfprint"
)

func makeLocalFilePublisher() *plugin.Publisher {
	return &plugin.Publisher{
		Doc:  "Publishes content to local file",
		Tags: []string{},
		Args: dataspec.ObjectSpec{
			&dataspec.AttrSpec{
				Name:        "path",
				Doc:         "Path to the file",
				Type:        cty.String,
				ExampleVal:  cty.StringVal("dist/output.md"),
				Constraints: constraint.Required,
			},
		},
		AllowedFormats: []plugin.OutputFormat{plugin.OutputFormatMD, plugin.OutputFormatHTML, plugin.OutputFormatPDF},
		PublishFunc:    publishLocalFile,
	}
}

func publishLocalFile(ctx context.Context, params *plugin.PublishParams) hcl.Diagnostics {
	document, _ := parseScope(params.DataContext)
	if document == nil {
		return hcl.Diagnostics{{
			Severity: hcl.DiagError,
			Summary:  "Failed to parse document",
			Detail:   "document is required",
		}}
	}
	datactx := params.DataContext
	datactx["format"] = plugin.StringData(params.Format.String())

	var printer printer.Printer
	switch params.Format {
	case plugin.OutputFormatMD:
		printer = mdprint.New()
	case plugin.OutputFormatHTML:
		printer = htmlprint.New()
	case plugin.OutputFormatPDF:
		printer = pdfprint.New()
	default:
		return hcl.Diagnostics{{
			Severity: hcl.DiagError,
			Summary:  "Unsupported format",
			Detail:   "Only md and html formats are supported",
		}}
	}

	pathAttr := params.Args.GetAttr("path")
	if pathAttr.IsNull() || pathAttr.AsString() == "" {
		return hcl.Diagnostics{{
			Severity: hcl.DiagError,
			Summary:  "Failed to parse arguments",
			Detail:   "path is required",
		}}
	}
	path, err := templatePath(pathAttr.AsString(), datactx)
	if err != nil {
		return hcl.Diagnostics{{
			Severity: hcl.DiagError,
			Summary:  "Failed to render path",
			Detail:   err.Error(),
		}}
	}
	slog.Info("Writing to file", "path", path)
	dir := filepath.Dir(path)
	err = os.MkdirAll(dir, 0o755)
	if err != nil {
		return hcl.Diagnostics{{
			Severity: hcl.DiagError,
			Summary:  "Failed to create directory",
			Detail:   err.Error(),
		}}
	}
	fs, err := os.Create(path)
	if err != nil {
		return hcl.Diagnostics{{
			Severity: hcl.DiagError,
			Summary:  "Failed to create file",
			Detail:   err.Error(),
		}}
	}
	defer fs.Close()
	err = printer.Print(fs, document)
	if err != nil {
		return hcl.Diagnostics{{
			Severity: hcl.DiagError,
			Summary:  "Failed to write to file",
			Detail:   err.Error(),
		}}
	}
	return nil
}

func templatePath(pattern string, datactx plugin.MapData) (string, error) {
	tmpl, err := template.New("pattern").Funcs(sprig.FuncMap()).Parse(pattern)
	if err != nil {
		return "", fmt.Errorf("failed to parse text template: %w", err)
	}

	var buf bytes.Buffer
	err = tmpl.Execute(&buf, datactx.Any())
	if err != nil {
		return "", fmt.Errorf("failed to execute text template: %w", err)
	}
	return filepath.Abs(filepath.Clean(strings.TrimSpace(buf.String())))
}
