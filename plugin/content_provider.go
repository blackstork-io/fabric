package plugin

import (
	"context"

	"github.com/hashicorp/hcl/v2"

	"github.com/blackstork-io/fabric/eval/dataquery"
	"github.com/blackstork-io/fabric/pkg/diagnostics"
	"github.com/blackstork-io/fabric/plugin/dataspec"
	"github.com/blackstork-io/fabric/plugin/plugindata"
)

type ContentProviders map[string]*ContentProvider

func (cp ContentProviders) Validate() diagnostics.Diag {
	var diags diagnostics.Diag
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
	// first non-empty line is treated as a short description
	Doc             string
	Tags            []string
	ContentFunc     ProvideContentFunc
	Args            *dataspec.RootSpec
	Config          *dataspec.RootSpec
	InvocationOrder InvocationOrder
}

func (cg *ContentProvider) Validate() diagnostics.Diag {
	var diags diagnostics.Diag
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

func (cg *ContentProvider) Execute(ctx context.Context, params *ProvideContentParams) (_ *ContentResult, diags diagnostics.Diag) {
	if cg == nil {
		return nil, diagnostics.Diag{{
			Severity: hcl.DiagError,
			Summary:  "Missing ContentProvider schema",
		}}
	}
	if cg.ContentFunc == nil {
		return nil, diagnostics.Diag{{
			Severity: hcl.DiagError,
			Summary:  "Incomplete ContentProviderSchema",
			Detail:   "content provider function not loaded",
		}}
	}
	diags.Extend(dataquery.EvaluateDeferredBlock(ctx, params.DataContext, params.Config))
	diags.Extend(dataquery.EvaluateDeferredBlock(ctx, params.DataContext, params.Args))
	if diags.HasErrors() {
		return
	}
	return cg.ContentFunc(ctx, params)
}

type ProvideContentParams struct {
	Config      *dataspec.Block
	Args        *dataspec.Block
	DataContext plugindata.Map
	ContentID   uint32
}

type ProvideContentFunc func(ctx context.Context, params *ProvideContentParams) (*ContentResult, diagnostics.Diag)
