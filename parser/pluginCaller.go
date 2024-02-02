package parser

import (
	"fmt"
	"io"
	"os/exec"

	"github.com/hashicorp/go-hclog"
	"github.com/hashicorp/go-plugin"
	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/hcldec"
	"github.com/zclconf/go-cty/cty"
	goctyjson "github.com/zclconf/go-cty/cty/json"

	"github.com/blackstork-io/fabric/parser/definitions"
	"github.com/blackstork-io/fabric/parser/evaluation"
	"github.com/blackstork-io/fabric/pkg/diagnostics"
	"github.com/blackstork-io/fabric/pkg/gobfix"
	"github.com/blackstork-io/fabric/pkg/utils"
	"github.com/blackstork-io/fabric/plugininterface/v1"
	"github.com/blackstork-io/fabric/plugins"
)

// Stub implementation of plugin caller
// TODO: attach to plugin discovery mechanism

type pluginKey struct {
	Kind string
	Name string
}

type pluginData struct {
	rpc            plugininterface.PluginRPCSer
	Version        plugininterface.Version
	ConfigSpec     hcldec.Spec
	InvocationSpec hcldec.Spec
}
type PluginCaller interface {
	CallContent(name string, config evaluation.Configuration, invocation evaluation.Invocation, context map[string]any) (result string, diag diagnostics.Diag)
	CallData(name string, config evaluation.Configuration, invocation evaluation.Invocation) (result map[string]any, diag diagnostics.Diag)
}

type Caller struct {
	plugins map[pluginKey]pluginData
}

func NewPluginCaller() *Caller {
	return &Caller{
		plugins: map[pluginKey]pluginData{},
	}
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

	args := plugininterface.ArgsSer{
		Kind:    key.Kind,
		Name:    key.Name,
		Version: data.Version,
		Context: context,
		// Config ans Args to be filled
	}

	var err error

	// TODO: check that nil interface values are checked like this everywhere
	needsConfig := !utils.IsNil(data.ConfigSpec)
	hasConfig := !utils.IsNil(config)

	switch {
	case needsConfig && hasConfig: // happy path
		var configVal cty.Value
		configVal, diag = config.ParseConfig(data.ConfigSpec)
		if !diags.Extend(diag) {
			// serialize only if no errors
			args.Config, err = goctyjson.Marshal(configVal, hcldec.ImpliedType(data.ConfigSpec))
			diags.AppendErr(err, "Error while serializing config")
		}
	case !needsConfig && !hasConfig:
		// happy path, do nothing
	case needsConfig && !hasConfig:
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
	case !needsConfig && hasConfig:
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
	if !diag.Extend(diag) {
		// serialize only if no errors
		args.Args, err = goctyjson.Marshal(pluginArgs, hcldec.ImpliedType(data.InvocationSpec))
		diags.AppendErr(err, "Error while serializing value")
	}

	if diag.HasErrors() {
		return
	}
	result := data.rpc.Call(args)
	for _, d := range result.Diags {
		diags.Append(&hcl.Diagnostic{
			Severity: d.Severity,
			Summary:  d.Summary,
			Detail:   d.Detail,
			Subject:  invocation.DefRange().Ptr(),
		})
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

func (c *Caller) LoadPluginBinary(pluginPath string) (diag diagnostics.Diag) {
	// TODO: setup pluggin logging?
	hclog.DefaultOutput = io.Discard

	client := plugin.NewClient(&plugin.ClientConfig{
		HandshakeConfig:  plugins.Handshake,
		VersionedPlugins: plugins.PluginMap,
		Cmd:              exec.Command(pluginPath),
		Managed:          true,
		// Logger:          hclog.de,
	})
	defer func() {
		if diag.HasErrors() {
			client.Kill()
		}
	}()

	// Connect via RPC
	rpcClient, err := client.Client()
	if diag.AppendErr(err, "Plugin connection error") {
		return
	}

	rawPlugin, err := rpcClient.Dispense(plugins.RPCPluginName)
	if diag.AppendErr(err, "Plugin RPC error") {
		return
	}

	pluginRPCSer, ok := rawPlugin.(plugininterface.PluginRPCSer)
	if !ok {
		diag.Add("RPC plugin doesn't conform to spec", "Contact Fabric developers about this")
		return
	}

	plugins := pluginRPCSer.GetPlugins()

	for _, plSer := range plugins {
		c.plugins[pluginKey{
			Kind: plSer.Kind,
			Name: plSer.Name,
		}] = pluginData{
			rpc:            pluginRPCSer,
			Version:        plSer.Version,
			ConfigSpec:     gobfix.ToHcl(plSer.ConfigSpec),
			InvocationSpec: gobfix.ToHcl(plSer.InvocationSpec),
		}
	}
	return
}
