package plugin

import (
	"context"

	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/hcldec"
	"github.com/zclconf/go-cty/cty"
)

type ContentProviders map[string]*ContentProvider

func (cp ContentProviders) Validate() hcl.Diagnostics {
	var diags hcl.Diagnostics
	for name, provider := range cp {
		if provider == nil {
			diags = append(diags, &hcl.Diagnostic{
				Severity: hcl.DiagError,
				Summary:  "Incomplete ContentProviderSchema",
				Detail:   "ContentProvider '" + name + "' not loaded",
			})
		} else {
			diags = append(diags, provider.Validate()...)
		}
	}
	return diags
}

type InvocationOrder int

const (
	InvocationOrderUnspecified InvocationOrder = iota
	InvocationOrderBegin
	InvocationOrderEnd
)

func (order InvocationOrder) Weight() int {
	switch order {
	case InvocationOrderBegin:
		return 0
	case InvocationOrderEnd:
		return 2
	default:
		return 1
	}
}

type ContentProvider struct {
	ContentFunc     ProvideContentFunc
	Args            hcldec.Spec
	Config          hcldec.Spec
	InvocationOrder InvocationOrder
}

func (cg *ContentProvider) Validate() hcl.Diagnostics {
	var diags hcl.Diagnostics
	if cg.ContentFunc == nil {
		diags = append(diags, &hcl.Diagnostic{
			Severity: hcl.DiagError,
			Summary:  "Incomplete ContentProviderSchema",
			Detail:   "ContentProvider function not loaded",
		})
	}
	if cg.Args == nil {
		diags = append(diags, &hcl.Diagnostic{
			Severity: hcl.DiagError,
			Summary:  "Incomplete ContentProviderSchema",
			Detail:   "Missing args schema",
		})
	}
	return diags
}

func (cg *ContentProvider) Execute(ctx context.Context, params *ProvideContentParams) (*Content, hcl.Diagnostics) {
	if cg == nil {
		return nil, hcl.Diagnostics{{
			Severity: hcl.DiagError,
			Summary:  "Missing ContentProvider schema",
		}}
	}
	if cg.ContentFunc == nil {
		return nil, hcl.Diagnostics{{
			Severity: hcl.DiagError,
			Summary:  "Incomplete ContentProviderSchema",
			Detail:   "content provider function not loaded",
		}}
	}
	return cg.ContentFunc(ctx, params)
}

type ProvideContentParams struct {
	Config      cty.Value
	Args        cty.Value
	DataContext MapData
}

type ProvideContentFunc func(ctx context.Context, params *ProvideContentParams) (*Content, hcl.Diagnostics)
