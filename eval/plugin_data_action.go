package eval

import (
	"context"
	"fmt"

	"github.com/hashicorp/hcl/v2"

	"github.com/blackstork-io/fabric/parser/definitions"
	"github.com/blackstork-io/fabric/pkg/diagnostics"
	"github.com/blackstork-io/fabric/plugin"
	"github.com/blackstork-io/fabric/plugin/dataspec"
	"github.com/blackstork-io/fabric/plugin/plugindata"
)

type PluginDataAction struct {
	*PluginAction
	Source   *plugin.DataSource
	SrcRange hcl.Range
}

func (action *PluginDataAction) FetchData(ctx context.Context) (plugindata.Data, diagnostics.Diag) {
	res, diags := action.Source.Execute(ctx, &plugin.RetrieveDataParams{
		Config: action.Config,
		Args:   action.Args,
	})
	diags.Refine(diagnostics.DefaultSubject(action.SrcRange))
	return res, diags
}

func LoadDataAction(ctx context.Context, sources DataSources, node *definitions.ParsedPlugin) (_ *PluginDataAction, diags diagnostics.Diag) {
	defer func() {
		diags.Refine(diagnostics.DefaultSubject(node.Invocation.Range()))
	}()

	ds, ok := sources.DataSource(node.PluginName)
	if !ok {
		return nil, diagnostics.Diag{{
			Severity: hcl.DiagError,
			Summary:  "Missing datasource",
			Detail:   fmt.Sprintf("'%s' not found in any plugin", node.PluginName),
		}}
	}
	var cfgBlock *dataspec.Block
	if ds.Config != nil {
		cfgBlock, diags = node.Config.ParseConfig(ctx, ds.Config)
		if diags.HasErrors() {
			return nil, diags
		}
	} else if node.Config.Exists() {
		diags.Append(&hcl.Diagnostic{
			Severity: hcl.DiagWarning,
			Summary:  "DataSource doesn't support configuration",
			Detail: fmt.Sprintf("DataSource '%s' does not support configuration, "+
				"but was provided with one. Remove it.", node.PluginName),
			Subject: node.Config.Range().Ptr(),
			Context: node.Invocation.Range().Ptr(),
		})
	}
	args, diag := dataspec.DecodeAndEvalBlock(ctx, node.Invocation.Block, ds.Args, nil)
	if diags.Extend(diag) {
		return nil, diags
	}
	return &PluginDataAction{
		PluginAction: &PluginAction{
			PluginName: node.PluginName,
			BlockName:  node.BlockName,
			Meta:       node.Meta,
			Config:     cfgBlock,
			Args:       args,
		},
		Source: ds,
	}, diags
}
