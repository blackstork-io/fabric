package plugin

import (
	"context"

	"github.com/hashicorp/hcl/v2"
)

type Schema struct {
	Name             string
	Version          string
	DataSources      DataSources
	ContentProviders ContentProviders
}

func (p *Schema) Validate() hcl.Diagnostics {
	var diags hcl.Diagnostics
	if p.Name == "" {
		diags = append(diags, &hcl.Diagnostic{
			Severity: hcl.DiagError,
			Summary:  "Incomplete PluginSchema",
			Detail:   "Name not defined",
		})
	}
	if p.DataSources != nil {
		diags = append(diags, p.DataSources.Validate()...)
	}
	if p.ContentProviders != nil {
		diags = append(diags, p.ContentProviders.Validate()...)
	}
	if p.DataSources == nil && p.ContentProviders == nil {
		diags = append(diags, &hcl.Diagnostic{
			Severity: hcl.DiagError,
			Summary:  "Incomplete PluginSchema",
			Detail:   "No data sources or content providers defined",
		})
	}
	return diags
}

func (p *Schema) RetrieveData(ctx context.Context, name string, params *RetrieveDataParams) (Data, hcl.Diagnostics) {
	if p == nil {
		return nil, hcl.Diagnostics{{
			Severity: hcl.DiagError,
			Summary:  "No schema",
			Detail:   "No schema defined",
		}}
	}
	if p.DataSources == nil {
		return nil, hcl.Diagnostics{{
			Severity: hcl.DiagError,
			Summary:  "No data sources",
			Detail:   "No data sources defined in schema",
		}}
	}
	source, ok := p.DataSources[name]
	if !ok || source == nil {
		return nil, hcl.Diagnostics{{
			Severity: hcl.DiagError,
			Summary:  "Data source not found",
			Detail:   "Data source " + name + " not found in schema",
		}}
	}
	return source.Execute(ctx, params)
}

func (p *Schema) ProvideContent(ctx context.Context, name string, params *ProvideContentParams) (*ContentResult, hcl.Diagnostics) {
	if p == nil {
		return nil, hcl.Diagnostics{{
			Severity: hcl.DiagError,
			Summary:  "No schema",
			Detail:   "No schema defined",
		}}
	}
	if p.ContentProviders == nil {
		return nil, hcl.Diagnostics{{
			Severity: hcl.DiagError,
			Summary:  "No content providers",
			Detail:   "No content providers defined in schema",
		}}
	}
	provider, ok := p.ContentProviders[name]
	if !ok || provider == nil {
		return nil, hcl.Diagnostics{{
			Severity: hcl.DiagError,
			Summary:  "Content provider not found",
			Detail:   "Content provider '" + name + "' not found in schema",
		}}
	}
	result, diags := provider.Execute(ctx, params)
	if diags.HasErrors() {
		return nil, diags
	}
	result.Content.setMeta(&ContentMeta{
		Provider: name,
		Plugin:   p.Name,
		Version:  p.Version,
	})
	return result, diags
}
