package main

import (
	"fmt"

	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/hcldec"
	goctyjson "github.com/zclconf/go-cty/cty/json"

	"github.com/blackstork-io/fabric/pkg/gobfix"
	"github.com/blackstork-io/fabric/pkg/utils"
	plugininterface "github.com/blackstork-io/fabric/plugininterface/v1"
)

type key struct {
	Kind    string
	Name    string
	Version plugininterface.Version
}

type val struct {
	Info *plugininterface.Plugin
	Call plugininterface.Callable
}

// Combines multiple plugins and presents them as an RPC interface.
type Multiplugin struct {
	plugins map[key]val
}

// Call implements plugininterface.PluginRPCSer.
func (m *Multiplugin) Call(args plugininterface.ArgsSer) (res plugininterface.ResultSer) {
	plugin, found := m.plugins[key{
		Kind:    args.Kind,
		Name:    args.Name,
		Version: args.Version,
	}]
	if !found {
		res.Diags = append(res.Diags, &plugininterface.RemoteDiag{
			Severity: hcl.DiagError,
			Summary:  "Plugin not found",
			Detail:   fmt.Sprintf("Plugin '%s %s' was not found", args.Kind, args.Name),
		})
		return
	}
	var err error
	argsDeser := plugininterface.Args{
		Kind:    args.Kind,
		Name:    args.Name,
		Version: args.Version,
		Context: args.Context,
	}
	if utils.IsNil(plugin.Info.InvocationSpec) {
		panic("InvocationSpec on plugins must not be nil")
	}
	argsDeser.Args, err = goctyjson.Unmarshal(args.Args, hcldec.ImpliedType(plugin.Info.InvocationSpec))
	if err != nil {
		res.Diags = append(res.Diags, &plugininterface.RemoteDiag{
			Severity: hcl.DiagError,
			Summary:  "Deserialization error",
			Detail:   err.Error(),
		})
		return
	}
	if !utils.IsNil(plugin.Info.ConfigSpec) {
		argsDeser.Config, err = goctyjson.Unmarshal(args.Config, hcldec.ImpliedType(plugin.Info.ConfigSpec))
		if err != nil {
			res.Diags = append(res.Diags, &plugininterface.RemoteDiag{
				Severity: hcl.DiagError,
				Summary:  "Deserialization error",
				Detail:   err.Error(),
			})
			return
		}
	}
	resUnser := plugin.Call(argsDeser)
	res.Result = resUnser.Result
	for _, v := range resUnser.Diags {
		res.Diags = append(res.Diags, &plugininterface.RemoteDiag{
			Severity: v.Severity,
			Summary:  v.Summary,
			Detail:   v.Detail,
		})
	}
	return
}

// GetPlugins implements plugininterface.PluginRPCSer.
func (m *Multiplugin) GetPlugins() (res []plugininterface.PluginSer) {
	for _, val := range m.plugins {
		ser := plugininterface.PluginSer{
			Namespace:      val.Info.Namespace,
			Kind:           val.Info.Kind,
			Name:           val.Info.Name,
			Version:        val.Info.Version,
			ConfigSpec:     gobfix.FromHcl(val.Info.ConfigSpec),
			InvocationSpec: gobfix.FromHcl(val.Info.InvocationSpec),
		}
		res = append(res, ser)
	}
	return
}

var _ plugininterface.PluginRPCSer = (*Multiplugin)(nil)

// TODO: this can accept something less generic than whole plugininterface.PluginRPC
func NewMultiplugin(plugins []plugininterface.PluginRPC) (res *Multiplugin) {
	res = &Multiplugin{
		plugins: map[key]val{},
	}
	for _, pluginsRPC := range plugins {
		plugins := pluginsRPC.GetPlugins()
		for _, plugin := range plugins {
			k := key{
				Kind:    plugin.Kind,
				Name:    plugin.Name,
				Version: plugin.Version,
			}
			res.plugins[k] = val{
				Info: &plugin,
				Call: pluginsRPC.Call,
			}
		}
	}
	return
}
