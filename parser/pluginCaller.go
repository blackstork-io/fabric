package parser

import (
	"fmt"
	"sync"

	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/hcldec"
	"github.com/sanity-io/litter"
	"golang.org/x/exp/maps"

	"github.com/blackstork-io/fabric/parser/definitions"
	"github.com/blackstork-io/fabric/parser/evaluation"
	"github.com/blackstork-io/fabric/pkg/diagnostics"
	plugin "github.com/blackstork-io/fabric/pluginInterface/v1"
)

// Stub implementation of plugin caller
// TODO: attach to plugin discovery mechanism

type pluginKey struct {
	Kind string
	Name string
}

type pluginData struct {
	rpc            plugin.PluginRPC
	Version        plugin.Version
	ConfigSpec     hcldec.Spec
	InvocationSpec hcldec.Spec
}
type PluginCaller interface {
	CallContent(name string, config evaluation.Configuration, invocation evaluation.Invocation, context map[string]any) (result string, diag diagnostics.Diag)
	CallData(name string, config evaluation.Configuration, invocation evaluation.Invocation) (result map[string]any, diag diagnostics.Diag)
}

type Caller struct {
	// TODO: make a ptr? why?
	plugins map[pluginKey]pluginData
}

var _ PluginCaller = (*Caller)(nil)

func (c *Caller) callPlugin(kind, name string, config evaluation.Configuration, invocation evaluation.Invocation, context map[string]any) (res any, diags diagnostics.Diag) {
	var diag diagnostics.Diag
	key := pluginKey{
		Kind: kind,
		Name: name,
	}
	data, found := c.plugins[key]
	if !found {
		diags.Add("Plugin not found", fmt.Sprintf("Plugin '%s %s' is missing!", kind, name))
		return
	}

	args := plugin.Args{
		Kind:    key.Kind,
		Name:    key.Name,
		Version: data.Version,
		Context: context,
		// Config ans Args to be filled
	}

	needsConfig := data.ConfigSpec != nil
	hasConfig := config != nil

	if needsConfig == hasConfig { // happy path
		if hasConfig {
			args.Config, diag = config.ParseConfig(data.ConfigSpec)
			diags.Extend(diag)
		}
	} else if !hasConfig { // config is needed but absent
		diags.Append(&hcl.Diagnostic{
			Severity: hcl.DiagError,
			Summary:  "Plugin requires configuration",
			Detail: fmt.Sprintf("Plugin '%s %s' has no default configuration and "+
				"no configuration was provided at the plugin invocation. "+
				"Provide an inline config block or a config attribute",
				kind, name),
			Subject: invocation.MissingItemRange().Ptr(),
			Context: invocation.Range().Ptr(),
		})
	} else { // config is present but not needed
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

	args.Args, diag = invocation.ParseInvocation(data.InvocationSpec)
	diag.Extend(diag)
	if diag.HasErrors() {
		return
	}
	result := data.rpc.Call(args)
	diags.ExtendHcl(result.Diags)
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

type MockCaller struct {
	once sync.Once
}

// CallContent implements PluginCaller.
func (c *MockCaller) CallContent(name string, config evaluation.Configuration, invocation evaluation.Invocation, context map[string]any) (result string, diag diagnostics.Diag) {
	c.once.Do(func() {
		result = litter.Sdump(context) + "\n"
	})
	attrs, _ := invocation.(*evaluation.BlockInvocation).JustAttributes()
	result += litter.Sdump("Call to content:", name, config, maps.Keys(attrs))
	return
}

// CallData implements PluginCaller.
func (*MockCaller) CallData(name string, config evaluation.Configuration, invocation evaluation.Invocation) (result map[string]any, diag diagnostics.Diag) {
	attrs, _ := invocation.(*evaluation.BlockInvocation).JustAttributes()

	result = map[string]any{
		"result": litter.Sdump("Call to data:", name, config, maps.Keys(attrs)),
	}
	return
}

var _ PluginCaller = (*MockCaller)(nil)
