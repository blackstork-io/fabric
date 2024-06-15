package eval

import (
	"context"

	"github.com/hashicorp/hcl/v2"

	"github.com/blackstork-io/fabric/parser/definitions"
	"github.com/blackstork-io/fabric/pkg/diagnostics"
	"github.com/blackstork-io/fabric/plugin"
)

type Content struct {
	Section *Section
	Plugin  *PluginContentAction
}

func (action *Content) InvocationOrder() plugin.InvocationOrder {
	if action.Plugin != nil {
		return action.Plugin.Provider.InvocationOrder
	}
	return plugin.InvocationOrderUnspecified
}

func (action *Content) RenderContent(ctx context.Context, dataCtx plugin.MapData, doc, parent *plugin.ContentSection, contentID uint32) (*plugin.ContentResult, diagnostics.Diag) {
	if action.Section != nil {
		return action.Section.RenderContent(ctx, dataCtx, doc, parent, contentID)
	}
	if action.Plugin != nil {
		return action.Plugin.RenderContent(ctx, dataCtx, doc, parent, contentID)
	}
	return nil, diagnostics.Diag{{
		Severity: hcl.DiagError,
		Summary:  "Content block not found",
	}}
}

func LoadContent(ctx context.Context, providers ContentProviders, node *definitions.ParsedContent) (_ *Content, diags diagnostics.Diag) {
	var block Content
	switch {
	case node.Plugin != nil:
		block.Plugin, diags = LoadPluginContentAction(ctx, providers, node.Plugin)
	case node.Section != nil:
		block.Section, diags = LoadSection(ctx, providers, node.Section)
	default:
		diags.Append(&hcl.Diagnostic{
			Severity: hcl.DiagError,
			Summary:  "Unsupported content block",
			Detail:   "Content block must be either 'content' or 'section'",
		})
	}
	if diags.HasErrors() {
		return nil, diags
	}
	return &block, diags
}
