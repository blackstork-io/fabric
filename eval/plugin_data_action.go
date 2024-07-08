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

type PluginDataAction struct {
	*PluginAction
	Source   *plugin.DataSource
	SrcRange hcl.Range
}

func (action *PluginDataAction) FetchData(ctx context.Context) (plugin.Data, diagnostics.Diag) {
	res, diags := action.Source.Execute(ctx, &plugin.RetrieveDataParams{
		Config: action.Config,
		Args:   action.Args,
	})
	diags.DefaultSubject(action.SrcRange.Ptr())
	return res, diags
}

func LoadDataAction(ctx context.Context, sources DataSources, node *definitions.ParsedPlugin) (_ *PluginDataAction, diags diagnostics.Diag) {
	defer func() {
		diags.DefaultSubject(node.Invocation.Range().Ptr())
	}()

	ds, ok := sources.DataSource(node.PluginName)
	if !ok {
		return nil, diagnostics.Diag{{
			Severity: hcl.DiagError,
			Summary:  "Missing datasource",
			Detail:   fmt.Sprintf("'%s' not found in any plugin", node.PluginName),
		}}
	}
	var cfg cty.Value
	if ds.Config != nil && !ds.Config.IsEmpty() {
		cfg, diags = node.Config.ParseConfig(ctx, ds.Config)
		if diags.HasErrors() {
			return nil, diags
		}
	} else if (ds.Config == nil || ds.Config.IsEmpty()) && node.Config.Exists() {
		diags.Append(&hcl.Diagnostic{
			Severity: hcl.DiagWarning,
			Summary:  "DataSource doesn't support configuration",
			Detail: fmt.Sprintf("DataSource '%s' does not support configuration, "+
				"but was provided with one. Remove it.", node.PluginName),
			Subject: node.Config.Range().Ptr(),
			Context: node.Invocation.Range().Ptr(),
		})
	}
	var args cty.Value
	args, diag := node.Invocation.ParseInvocation(ctx, ds.Args)
	if diags.Extend(diag) {
		return nil, diags
	}
	return &PluginDataAction{
		PluginAction: &PluginAction{
			PluginName: node.PluginName,
			BlockName:  node.BlockName,
			Meta:       node.Meta,
			Config:     cfg,
			Args:       args,
		},
		Source: ds,
	}, diags
}
