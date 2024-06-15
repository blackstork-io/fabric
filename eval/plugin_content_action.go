package eval

import (
	"context"
	"fmt"

	"github.com/hashicorp/hcl/v2"
	"github.com/zclconf/go-cty/cty"

	"github.com/blackstork-io/fabric/parser/definitions"
	"github.com/blackstork-io/fabric/pkg/diagnostics"
	"github.com/blackstork-io/fabric/plugin"
)

type PluginContentAction struct {
	*PluginAction
	Provider *plugin.ContentProvider
	Vars     *definitions.ParsedVars
}

func (action *PluginContentAction) RenderContent(ctx context.Context, dataCtx plugin.MapData, doc, parent *plugin.ContentSection, contentID uint32) (res *plugin.ContentResult, diags diagnostics.Diag) {
	var diag diagnostics.Diag
	contentMap := plugin.MapData{}
	if action.PluginAction.Meta != nil {
		contentMap[definitions.BlockKindMeta] = action.PluginAction.Meta.AsJQData()
	}
	docData := dataCtx[definitions.BlockKindDocument]
	docData.(plugin.MapData)[definitions.BlockKindContent] = doc.AsData()
	dataCtx[definitions.BlockKindDocument] = docData
	dataCtx[definitions.BlockKindContent] = contentMap
	diag = ApplyVars(ctx, action.Vars, dataCtx)
	if diags.Extend(diag) {
		return
	}

	res, diag = action.Provider.Execute(ctx, &plugin.ProvideContentParams{
		Config:      action.Config,
		Args:        action.Args,
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
	var cfg cty.Value
	if cp.Config != nil && !cp.Config.IsEmpty() {
		cfg, diags = node.Config.ParseConfig(ctx, cp.Config)
		if diags.HasErrors() {
			return nil, diags
		}
	} else if (cp.Config == nil || cp.Config.IsEmpty()) && node.Config.Exists() {
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

	var args cty.Value
	args, diag := node.Invocation.ParseInvocation(ctx, cp.Args)
	if diags.Extend(diag) {
		return nil, diags
	}
	return &PluginContentAction{
		PluginAction: &PluginAction{
			PluginName: node.PluginName,
			BlockName:  node.BlockName,
			Meta:       node.Meta,
			Config:     cfg,
			Args:       args,
		},
		Provider: cp,
		Vars:     node.Vars,
	}, diags
}
