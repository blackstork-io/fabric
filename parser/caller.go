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
	"github.com/blackstork-io/fabric/pkg/utils"
	"github.com/blackstork-io/fabric/plugin"
	"github.com/blackstork-io/fabric/plugin/runner"
)

// Stub implementation of plugin caller
// TODO: attach to plugin discovery mechanism

type pluginData struct {
	ConfigSpec     hcldec.Spec
	InvocationSpec hcldec.Spec
}
type PluginCaller interface {
	CallContent(name string, config evaluation.Configuration, invocation evaluation.Invocation, context map[string]any) (result string, diag diagnostics.Diag)
	CallData(name string, config evaluation.Configuration, invocation evaluation.Invocation) (result map[string]any, diag diagnostics.Diag)
}

type Caller struct {
	plugins *runner.Runner
}

func NewPluginCaller(r *runner.Runner) *Caller {
	return &Caller{
		plugins: r,
	}
}

var _ PluginCaller = (*Caller)(nil)

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

func (c *Caller) callPlugin(kind, name string, config evaluation.Configuration, invocation evaluation.Invocation, dataCtx map[string]any) (res any, diags diagnostics.Diag) {
	data, diags := c.pluginData(kind, name)
	if diags.HasErrors() {
		return
	}

	dataCtxAny, err := plugin.ParseDataMapAny(dataCtx)
	if err != nil {
		diags.Add("Error while parsing context", err.Error())
		return
	}

	acceptsConfig := !utils.IsNil(data.ConfigSpec)
	hasConfig := config.Exists()

	var configVal cty.Value
	if acceptsConfig {
		var stdDiag diagnostics.Diag
		configVal, stdDiag = config.ParseConfig(data.ConfigSpec)
		if !diags.Extend(stdDiag) {
			typ := hcldec.ImpliedType(data.ConfigSpec)
			errs := configVal.Type().TestConformance(typ)
			if errs != nil {
				// Attempt a conversion
				var err error
				configVal, err = convert.Convert(configVal, typ)
				if err != nil {
					diags.AppendErr(err, "Error while serializing config")
				}
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
	if diag.HasErrors() {
		return
	}
	if data.InvocationSpec != nil {
		typ := hcldec.ImpliedType(data.InvocationSpec)
		errs := pluginArgs.Type().TestConformance(typ)
		if errs != nil {
			// Attempt a conversion
			var err error
			pluginArgs, err = convert.Convert(pluginArgs, typ)
			if err != nil {
				diag.AppendErr(err, "Error while serializing args")

				return nil, diag
			}
		}
	}

	var result struct {
		Result any
		Diags  hcl.Diagnostics
	}
	switch kind {
	case "data":
		source, diags := c.plugins.DataSource(name)
		if diags.HasErrors() {
			return nil, diagnostics.Diag(diags)
		}
		data, diags := source.Execute(context.Background(), &plugin.RetrieveDataParams{
			Config: configVal,
			Args:   pluginArgs,
		})
		if data != nil {
			result.Result = data.Any()
		}
		result.Diags = diags
	case "content":
		provider, diags := c.plugins.ContentProvider(name)
		if diags.HasErrors() {
			return nil, diagnostics.Diag(diags)
		}
		content, diags := provider.Execute(context.Background(), &plugin.ProvideContentParams{
			Config:      configVal,
			Args:        pluginArgs,
			DataContext: dataCtxAny,
		})
		result.Result = ""
		if content != nil {
			result.Result = content.Markdown
		}
		result.Diags = diags
	}

	for _, d := range result.Diags {
		diags = append(diags, d)
	}
	res = result.Result
	return
}

func (c *Caller) CallContent(name string, config evaluation.Configuration, invocation evaluation.Invocation, context map[string]any) (result string, diag diagnostics.Diag) {
	var ok bool
	var res any
	res, diag = c.callPlugin(definitions.BlockKindContent, name, config, invocation, context)
	result, ok = res.(string)
	if !diag.HasErrors() && !ok {
		diag.Append(&hcl.Diagnostic{
			Severity: hcl.DiagError,
			Summary:  "Incorrect result type",
			Detail:   "Plugin returned incorrect data type. Please contact fabric developers about this issue",
			Subject:  invocation.DefRange().Ptr(),
		})
	}
	return
}

func (c *Caller) CallData(name string, config evaluation.Configuration, invocation evaluation.Invocation) (result map[string]any, diag diagnostics.Diag) {
	var ok bool
	var res any
	res, diag = c.callPlugin(definitions.BlockKindData, name, config, invocation, nil)
	result, ok = res.(map[string]any)
	if !diag.HasErrors() && !ok {
		diag.Append(&hcl.Diagnostic{
			Severity: hcl.DiagError,
			Summary:  "Incorrect result type",
			Detail:   "Plugin returned incorrect data type. Please contact fabric developers about this issue",
			Subject:  invocation.DefRange().Ptr(),
		})
	}
	return
}
