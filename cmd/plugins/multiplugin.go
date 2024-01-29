package main

import (
	"fmt"
	"log"

	"github.com/hashicorp/hcl/v2"
	"github.com/sanity-io/litter"

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

// Call implements plugininterface.PluginRPC.
func (m *Multiplugin) Call(args plugininterface.Args) (res plugininterface.Result) {
	plugin, found := m.plugins[key{
		Kind:    args.Kind,
		Name:    args.Name,
		Version: args.Version,
	}]
	if !found {
		res.Diags = res.Diags.Append(&hcl.Diagnostic{
			Severity: hcl.DiagError,
			Summary:  "Plugin not found",
			Detail:   fmt.Sprintf("Plugin '%s %s' was not found", args.Kind, args.Name),
		})
		return
	}
	return plugin.Call(args)
}

// GetPlugins implements plugininterface.PluginRPC.
func (m *Multiplugin) GetPlugins() (res []plugininterface.Plugin) {
	for _, val := range m.plugins {
		res = append(res, *val.Info)
	}
	return
}

var _ plugininterface.PluginRPC = (*Multiplugin)(nil)

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
	log.Println("res.plugins:", litter.Sdump(res.plugins))
	return
}
