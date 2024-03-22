package parser

import (
	"context"
	"fmt"

	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/hcldec"
	"github.com/zclconf/go-cty/cty"
	"github.com/zclconf/go-cty/cty/convert"

	"github.com/blackstork-io/fabric/parser/definitions"
	"github.com/blackstork-io/fabric/parser/evaluation"
	"github.com/blackstork-io/fabric/pkg/diagnostics"
	"github.com/blackstork-io/fabric/pkg/fabctx"
	"github.com/blackstork-io/fabric/pkg/utils"
	"github.com/blackstork-io/fabric/plugin"
	"github.com/blackstork-io/fabric/plugin/runner"
)

type pluginData struct {
	ConfigSpec     hcldec.Spec
	InvocationSpec hcldec.Spec
}

type Caller struct {
	plugins *runner.Runner
}

func NewPluginCaller(r *runner.Runner) *Caller {
	return &Caller{
		plugins: r,
	}
}

var _ evaluation.PluginCaller = (*Caller)(nil)

func (c *Caller) pluginData(kind, name string) (pluginData, diagnostics.Diag) {
	switch kind {
	case "data":
		plugin, diag := c.plugins.DataSource(name)
		if diag.HasErrors() {
			return pluginData{}, diagnostics.Diag(diag)
		}
		return pluginData{
			ConfigSpec:     plugin.Config,
			InvocationSpec: plugin.Args,
		}, nil
	case "content":
		plugin, diag := c.plugins.ContentProvider(name)
		if diag.HasErrors() {
			return pluginData{}, diagnostics.Diag(diag)
		}
		return pluginData{
			ConfigSpec:     plugin.Config,
			InvocationSpec: plugin.Args,
		}, nil
	default:
		return pluginData{}, diagnostics.Diag{
			{
				Severity: hcl.DiagError,
				Summary:  "Unknown plugin kind",
				Detail:   fmt.Sprintf("Unknown plugin kind '%s'", kind),
			},
		}
	}
}

func (c *Caller) callPlugin(ctx context.Context, kind, name string, config evaluation.Configuration, invocation evaluation.Invocation, dataCtx plugin.MapData) (res any, diags diagnostics.Diag) {
	data, diags := c.pluginData(kind, name)
	if diags.HasErrors() {
		return
	}

	acceptsConfig := !utils.IsNil(data.ConfigSpec)
	hasConfig := config.Exists()

	var configVal cty.Value
	if acceptsConfig {
		var diag diagnostics.Diag
		configVal, diag = config.ParseConfig(data.ConfigSpec)
		if !diags.Extend(diag) {
			typ := hcldec.ImpliedType(data.ConfigSpec)
			errs := configVal.Type().TestConformance(typ)
			if errs != nil {
				// Attempt a conversion
				var err error
				configVal, err = convert.Convert(configVal, typ)
				diags.AppendErr(err, "Error while serializing config")
			}
		}
	} else if hasConfig {
		diags.Append(&hcl.Diagnostic{
			Severity: hcl.DiagWarning,
			Summary:  "Plugin doesn't support configuration",
			Detail: fmt.Sprintf("Plugin '%s %s' does not support configuration, "+
				"but was provided with one. Remove it.",
				kind, name),
			Subject: config.Range().Ptr(),
			Context: invocation.Range().Ptr(),
		})
	}

	pluginArgs, diag := invocation.ParseInvocation(data.InvocationSpec)
	diags.Extend(diag)
	if data.InvocationSpec != nil {
		typ := hcldec.ImpliedType(data.InvocationSpec)
		errs := pluginArgs.Type().TestConformance(typ)
		if errs != nil {
			// Attempt a conversion
			var err error
			pluginArgs, err = convert.Convert(pluginArgs, typ)
			diags.AppendErr(err, "Error while serializing args")
		}
	}
	if diags.HasErrors() {
		return
	}

	switch kind {
	case "data":
		if fabctx.Get(ctx).IsLinting() {
			res = plugin.MapData{}
			return
		}
		source, diag := c.plugins.DataSource(name)
		if diags.ExtendHcl(diag) {
			return
		}
		data, diag := source.Execute(ctx, &plugin.RetrieveDataParams{
			Config: configVal,
			Args:   pluginArgs,
		})
		res = data
		diags.ExtendHcl(diag)
	case "content":
		if fabctx.Get(ctx).IsLinting() {
			res = ""
			return
		}
		provider, diag := c.plugins.ContentProvider(name)
		if diags.ExtendHcl(diag) {
			return
		}
		content, diag := provider.Execute(ctx, &plugin.ProvideContentParams{
			Config:      configVal,
			Args:        pluginArgs,
			DataContext: dataCtx,
		})
		res = content
		diags.ExtendHcl(diag)
	}
	return
}

func (c *Caller) CallContent(ctx context.Context, name string, config evaluation.Configuration, invocation evaluation.Invocation, dataCtx plugin.MapData) (result *plugin.Content, diag diagnostics.Diag) {
	var ok bool
	var res any
	res, diag = c.callPlugin(ctx, definitions.BlockKindContent, name, config, invocation, dataCtx)
	if diag.HasErrors() {
		return
	}
	result, ok = res.(*plugin.Content)
	if !ok {
		panic("Incorrect plugin result type")
	}
	return
}

func (c *Caller) ContentInvocationOrder(ctx context.Context, name string) (order plugin.InvocationOrder, diag diagnostics.Diag) {
	content, hclDiag := c.plugins.ContentProvider(name)
	if hclDiag.HasErrors() {
		return plugin.InvocationOrderUnspecified, diagnostics.Diag(hclDiag)
	}
	order = content.InvocationOrder
	return
}

func (c *Caller) CallData(ctx context.Context, name string, config evaluation.Configuration, invocation evaluation.Invocation) (result plugin.Data, diag diagnostics.Diag) {
	var ok bool
	var res any
	res, diag = c.callPlugin(ctx, definitions.BlockKindData, name, config, invocation, nil)
	if diag.HasErrors() {
		return
	}
	result, ok = res.(plugin.Data)
	if !ok {
		panic("Incorrect plugin result type")
	}
	return
}
