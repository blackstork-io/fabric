package plugin

import (
	"context"

	"github.com/hashicorp/hcl/v2"

	"github.com/blackstork-io/fabric/pkg/diagnostics"
	"github.com/blackstork-io/fabric/plugin/dataspec"
	"github.com/blackstork-io/fabric/plugin/plugindata"
)

type FormatFunc func(ctx context.Context, params *FormatParams) diagnostics.Diag

type FormatParams struct {
	DocumentName string
	Config       *dataspec.Block
	Args         *dataspec.Block
	DataContext  plugindata.Map
}

type Formatter struct {
	Doc        string
	Format     string
	FileExt    string
	FormatFunc FormatFunc
	Args       *dataspec.RootSpec
	Config     *dataspec.RootSpec
}

func (formatter *Formatter) Validate() diagnostics.Diag {
	var diags diagnostics.Diag
	if formatter.FormatFunc == nil {
		diags = append(diags, &hcl.Diagnostic{
			Severity: hcl.DiagError,
			Summary:  "Incomplete Formatter schema",
			Detail:   "Formatter function not loaded",
		})
	}
	if formatter.Args == nil {
		diags = append(diags, &hcl.Diagnostic{
			Severity: hcl.DiagError,
			Summary:  "Incomplete Formatter schema",
			Detail:   "Missing args schema",
		})
	}
	return diags
}

func (formatter *Formatter) Execute(ctx context.Context, params *FormatParams) (diags diagnostics.Diag) {
	if formatter == nil {
		return diagnostics.Diag{{
			Severity: hcl.DiagError,
			Summary:  "Missing Formatter schema",
		}}
	}
	if formatter.FormatFunc == nil {
		return diagnostics.Diag{{
			Severity: hcl.DiagError,
			Summary:  "Incomplete Formatter schema",
			Detail:   "Format function not loaded",
		}}
	}
	return formatter.FormatFunc(ctx, params)
}

type Formatters map[string]*Formatter

func (formatters Formatters) Validate() diagnostics.Diag {
	var diags diagnostics.Diag
	for name, formatter := range formatters {
		if formatter == nil {
			diags = append(diags, &hcl.Diagnostic{
				Severity: hcl.DiagError,
				Summary:  "Incomplete Formatter schema",
				Detail:   "Formatter '" + name + "' not loaded",
			})
		} else {
			diags = append(diags, formatter.Validate()...)
		}
	}
	return diags
}
