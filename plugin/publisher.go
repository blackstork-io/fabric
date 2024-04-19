package plugin

import (
	"context"

	"github.com/hashicorp/hcl/v2"
	"github.com/zclconf/go-cty/cty"

	"github.com/blackstork-io/fabric/plugin/dataspec"
)

type OutputFormat int

const (
	OutputFormatUnspecified OutputFormat = iota
	OutputFormatMD
	OutputFormatHTML
	OutputFormatPDF
)

func (f OutputFormat) String() string {
	switch f {
	case OutputFormatMD:
		return "md"
	case OutputFormatHTML:
		return "html"
	case OutputFormatPDF:
		return "pdf"
	default:
		return "unknown"
	}
}

func (f OutputFormat) Ext() string {
	return "." + f.String()
}

type PublishFunc func(ctx context.Context, params *PublishParams) hcl.Diagnostics

type PublishParams struct {
	Config      cty.Value
	Args        cty.Value
	DataContext MapData
	Format      OutputFormat
}

type Publisher struct {
	Doc            string
	Tags           []string
	PublishFunc    PublishFunc
	Args           dataspec.RootSpec
	Config         dataspec.RootSpec
	AllowedFormats []OutputFormat
}

func (pub *Publisher) Validate() hcl.Diagnostics {
	var diags hcl.Diagnostics
	if pub.PublishFunc == nil {
		diags = append(diags, &hcl.Diagnostic{
			Severity: hcl.DiagError,
			Summary:  "Incomplete Publisher shema",
			Detail:   "Publisher function not loaded",
		})
	}
	if pub.Args == nil {
		diags = append(diags, &hcl.Diagnostic{
			Severity: hcl.DiagError,
			Summary:  "Incomplete Publisher shema",
			Detail:   "Missing args schema",
		})
	}
	return diags
}

func (pub *Publisher) Execute(ctx context.Context, params *PublishParams) hcl.Diagnostics {
	if pub == nil {
		return hcl.Diagnostics{{
			Severity: hcl.DiagError,
			Summary:  "Missing Publisher schema",
		}}
	}
	if pub.PublishFunc == nil {
		return hcl.Diagnostics{{
			Severity: hcl.DiagError,
			Summary:  "Incomplete Publisher schema",
			Detail:   "Publish function not loaded",
		}}
	}
	return pub.PublishFunc(ctx, params)
}

type Publishers map[string]*Publisher

func (pubs Publishers) Validate() hcl.Diagnostics {
	var diags hcl.Diagnostics
	for name, provider := range pubs {
		if provider == nil {
			diags = append(diags, &hcl.Diagnostic{
				Severity: hcl.DiagError,
				Summary:  "Incomplete Publisher schema",
				Detail:   "Publisher '" + name + "' not loaded",
			})
		} else {
			diags = append(diags, provider.Validate()...)
		}
	}
	return diags
}
