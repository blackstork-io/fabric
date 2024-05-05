package plugin

import (
	"context"
	"fmt"
	"slices"

	"github.com/hashicorp/hcl/v2"
)

type Schema struct {
	Name             string
	Version          string
	Doc              string
	Tags             []string
	DataSources      DataSources
	ContentProviders ContentProviders
	Publishers       Publishers
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
	if p.Publishers != nil {
		diags = append(diags, p.Publishers.Validate()...)
	}
	if p.DataSources == nil && p.ContentProviders == nil && p.Publishers == nil {
		diags = append(diags, &hcl.Diagnostic{
			Severity: hcl.DiagError,
			Summary:  "Incomplete PluginSchema",
			Detail:   "No data sources, content providers or publishers defined",
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
			Detail:   fmt.Sprintf("Data source '%s' not found in schema", name),
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
			Detail:   fmt.Sprintf("Content provider '%s' not found in schema", name),
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

func (p *Schema) Publish(ctx context.Context, name string, params *PublishParams) hcl.Diagnostics {
	if p == nil {
		return hcl.Diagnostics{{
			Severity: hcl.DiagError,
			Summary:  "No schema",
			Detail:   "No schema defined",
		}}
	}
	if p.Publishers == nil {
		return hcl.Diagnostics{{
			Severity: hcl.DiagError,
			Summary:  "No publishers",
			Detail:   "No publishers defined in schema",
		}}
	}
	publisher, ok := p.Publishers[name]
	if !ok || publisher == nil {
		return hcl.Diagnostics{{
			Severity: hcl.DiagError,
			Summary:  "Publisher not found",
			Detail:   fmt.Sprintf("Publisher '%s' not found in schema", name),
		}}
	}
	if !slices.Contains(publisher.AllowedFormats, params.Format) {
		return hcl.Diagnostics{{
			Severity: hcl.DiagError,
			Summary:  "Invalid format",
			Detail:   fmt.Sprintf("Publisher '%s' does not support format '%s'", name, params.Format),
		}}
	}
	return publisher.Execute(ctx, params)
}
