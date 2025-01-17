package eval

import (
	"context"

	"github.com/hashicorp/hcl/v2"

	"github.com/blackstork-io/fabric/parser/definitions"
	"github.com/blackstork-io/fabric/pkg/diagnostics"
	"github.com/blackstork-io/fabric/plugin"
	"github.com/blackstork-io/fabric/plugin/plugindata"
)

type Content struct {
	Section *Section
	Plugin  *PluginContentAction
	Dynamic *Dynamic
}

func (action *Content) RenderContent(ctx context.Context, dataCtx plugindata.Map) (plugin.Content, diagnostics.Diag) {
	if action.Section != nil {
		return action.Section.RenderContent(ctx, dataCtx)
	}
	if action.Plugin != nil {
		return action.Plugin.RenderContent(ctx, dataCtx)
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
	case node.Dynamic != nil:
		block.Dynamic, diags = LoadDynamic(ctx, providers, node.Dynamic)
	default:
		diags.Append(&hcl.Diagnostic{
			Severity: hcl.DiagError,
			Summary:  "Unsupported content block",
			Detail:   "Content block must be either 'content', 'section' or 'dynamic'",
		})
	}
	if diags.HasErrors() {
		return nil, diags
	}
	return &block, diags
}
