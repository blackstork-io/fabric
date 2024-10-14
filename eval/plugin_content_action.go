package eval

import (
	"context"
	"fmt"

	"github.com/hashicorp/hcl/v2"

	"github.com/blackstork-io/fabric/parser/definitions"
	"github.com/blackstork-io/fabric/pkg/diagnostics"
	"github.com/blackstork-io/fabric/plugin"
	"github.com/blackstork-io/fabric/plugin/dataspec"
	"github.com/blackstork-io/fabric/plugin/dataspec/deferred"
	"github.com/blackstork-io/fabric/plugin/plugindata"
)

type PluginContentAction struct {
	*PluginAction
	Provider     *plugin.ContentProvider
	Vars         *definitions.ParsedVars
	RequiredVars []string
}

func (action *PluginContentAction) RenderContent(ctx context.Context, dataCtx plugindata.Map, doc, parent *plugin.ContentSection, contentID uint32) (res *plugin.ContentResult, diags diagnostics.Diag) {
	contentMap := plugindata.Map{}
	if action.PluginAction.Meta != nil {
		contentMap[definitions.BlockKindMeta] = action.PluginAction.Meta.AsPluginData()
	}
	docData := dataCtx[definitions.BlockKindDocument]
	docData.(plugindata.Map)[definitions.BlockKindContent] = doc.AsData()
	dataCtx[definitions.BlockKindDocument] = docData
	dataCtx[definitions.BlockKindContent] = contentMap
	diag := ApplyVars(ctx, action.Vars, dataCtx)
	if diags.Extend(diag) {
		return
	}

	if len(action.RequiredVars) > 0 {
		diag := verifyRequiredVars(dataCtx, action.RequiredVars, action.Source.Block)
		if diags.Extend(diag) {
			return
		}
	}
	evaluatedBlock, diag := dataspec.EvalBlockCopy(ctx, action.Args, dataCtx)
	if diags.Extend(diag) {
		return
	}

	res, diag = action.Provider.Execute(ctx, &plugin.ProvideContentParams{
		Config:      action.Config,
		Args:        evaluatedBlock,
		DataContext: dataCtx,
		ContentID:   contentID,
	})
	if diags.Extend(diag) {
		return
	}
	if res.Location == nil {
		res.Location = &plugin.Location{
			Index: contentID,
		}
	}
	parent.Add(res.Content, res.Location)
	return res, diags
}

func LoadPluginContentAction(ctx context.Context, providers ContentProviders, node *definitions.ParsedPlugin) (_ *PluginContentAction, diags diagnostics.Diag) {
	cp, ok := providers.ContentProvider(node.PluginName)
	if !ok {
		return nil, diagnostics.Diag{{
			Severity: hcl.DiagError,
			Summary:  "Missing content provider",
			Detail:   fmt.Sprintf("'%s' not found in any plugin", node.PluginName),
		}}
	}
	var cfg *dataspec.Block
	if cp.Config != nil {
		cfg, diags = node.Config.ParseConfig(ctx, cp.Config)
		if diags.HasErrors() {
			return nil, diags
		}
	} else if node.Config.Exists() {
		diags.Append(&hcl.Diagnostic{
			Severity: hcl.DiagWarning,
			Summary:  "ContentProvider doesn't support configuration",
			Detail: fmt.Sprintf("ContentProvider '%s' does not support configuration, "+
				"but was provided with one. Remove it.", node.PluginName),
			Subject: node.Config.Range().Ptr(),
			Context: node.Invocation.Range().Ptr(),
		})
		return nil, diags
	}

	args, diag := dataspec.DecodeBlock(
		deferred.WithQueryFuncs(ctx),
		node.Invocation.Block,
		cp.Args,
	)
	if diags.Extend(diag) {
		return nil, diags
	}
	return &PluginContentAction{
		PluginAction: &PluginAction{
			Source:     node.Source,
			PluginName: node.PluginName,
			BlockName:  node.BlockName,
			Meta:       node.Meta,
			Config:     cfg,
			Args:       args,
		},
		Provider:     cp,
		Vars:         node.Vars,
		RequiredVars: node.RequiredVars,
	}, diags
}
