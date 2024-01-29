package parser

import (
	"fmt"
	"log"
	"os/exec"

	"github.com/hashicorp/go-plugin"
	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/hcldec"
	"github.com/sanity-io/litter"

	"github.com/blackstork-io/fabric/parser/definitions"
	"github.com/blackstork-io/fabric/parser/evaluation"
	"github.com/blackstork-io/fabric/pkg/diagnostics"
	"github.com/blackstork-io/fabric/pkg/utils"
	plugininterface "github.com/blackstork-io/fabric/plugininterface/v1"
	"github.com/blackstork-io/fabric/plugins"
)

// Stub implementation of plugin caller
// TODO: attach to plugin discovery mechanism

type pluginKey struct {
	Kind string
	Name string
}

type pluginData struct {
	rpc            plugininterface.PluginRPC
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

	args := plugininterface.Args{
		Kind:    key.Kind,
		Name:    key.Name,
		Version: data.Version,
		Context: context,
		// Config ans Args to be filled
	}

	// TODO: check that nil interface values are checked like this everywhere
	needsConfig := !utils.IsNil(data.ConfigSpec)
	hasConfig := !utils.IsNil(config)
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

func (c *Caller) LoadPluginBinary(pluginPath string) (diag diagnostics.Diag) {
	// TODO: setup pluggin logging?
	// hclog.DefaultOutput = io.Discard

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

	pluginRPC, ok := rawPlugin.(plugininterface.PluginRPC)
	if !ok {
		diag.Add("RPC plugin doesn't conform to spec", "Contact Fabric developers about this")
		return
	}
	plugins := pluginRPC.GetPlugins()

	for _, pl := range plugins {
		log.Println("discovered", pl.Namespace, pl.Kind, pl.Name, pl.Version.Cast().String())
		litter.Dump(pl)
		c.plugins[pluginKey{
			Kind: pl.Kind,
			Name: pl.Name,
		}] = pluginData{
			rpc:            pluginRPC,
			Version:        pl.Version,
			ConfigSpec:     pl.ConfigSpec,
			InvocationSpec: pl.InvocationSpec,
		}
	}
	return
}
