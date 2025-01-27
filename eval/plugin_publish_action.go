package eval

import (
	"context"
	"fmt"

	"github.com/hashicorp/hcl/v2"

	"github.com/blackstork-io/fabric/parser/definitions"
	"github.com/blackstork-io/fabric/pkg/diagnostics"
	"github.com/blackstork-io/fabric/plugin"
	"github.com/blackstork-io/fabric/plugin/ast/nodes"
	"github.com/blackstork-io/fabric/plugin/dataspec"
	"github.com/blackstork-io/fabric/plugin/plugindata"
)

type PluginPublishAction struct {
	*PluginAction
	Publisher     *plugin.Publisher
	nodeRenderers NodeRenderers
}

func (block *PluginPublishAction) Publish(ctx context.Context, dataCtx plugindata.Map, documentName string, document *nodes.Node) (diags diagnostics.Diag) {
	// TODO: ask publisher for support info, ignore natively supported nodes
	publisherInfo, diag := block.Publisher.Info(ctx, &plugin.PublisherInfoParams{
		Config: block.Config,
		Args:   block.Args,
	})
	if diags.Extend(diag) {
		return
	}

	customRenderer := newCustomNodeRenderer(document, block.PluginName, publisherInfo, block.nodeRenderers)
	document, diag = customRenderer.renderNodes(ctx)
	if diags.Extend(diag) {
		return
	}

	return block.Publisher.Execute(ctx, &plugin.PublishParams{
		Config:       block.Config,
		Args:         block.Args,
		DataContext:  dataCtx,
		DocumentName: documentName,
		Document:     document,
	})
}

func LoadPluginPublishAction(ctx context.Context, plugins Plugins, node *definitions.ParsedPlugin) (_ *PluginPublishAction, diags diagnostics.Diag) {
	p, ok := plugins.Publisher(node.PluginName)
	if !ok {
		return nil, diagnostics.Diag{{
			Severity: hcl.DiagError,
			Summary:  "Missing publisher",
			Detail:   fmt.Sprintf("'%s' not found in any plugin", node.PluginName),
		}}
	}
	var cfg *dataspec.Block
	if p.Config != nil {
		cfg, diags = node.Config.ParseConfig(ctx, p.Config)
		if diags.HasErrors() {
			return nil, diags
		}
	} else if node.Config.Exists() {
		diags.Append(&hcl.Diagnostic{
			Severity: hcl.DiagWarning,
			Summary:  "Publisher doesn't support configuration",
			Detail: fmt.Sprintf("Publisher '%s' does not support configuration, "+
				"but was provided with one. Remove it.", node.PluginName),
			Subject: node.Config.Range().Ptr(),
			Context: node.Invocation.Range().Ptr(),
		})
		return nil, diags
	}

	args, diag := dataspec.DecodeAndEvalBlock(ctx, node.Invocation.Block, p.Args, nil)
	if diags.Extend(diag) {
		return nil, diags
	}
	return &PluginPublishAction{
		PluginAction: &PluginAction{
			PluginName: node.PluginName,
			BlockName:  node.BlockName,
			Meta:       node.Meta,
			Config:     cfg,
			Args:       args,
		},
		nodeRenderers: plugins,
	}, diags
}
