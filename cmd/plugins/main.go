package main

import (
	"github.com/hashicorp/go-plugin"

	plugininterface "github.com/blackstork-io/fabric/plugininterface/v1"
	"github.com/blackstork-io/fabric/plugins"
	"github.com/blackstork-io/fabric/plugins/content/text"
)

func main() {
	plugin.Serve(&plugin.ServeConfig{
		HandshakeConfig: plugins.Handshake,
		VersionedPlugins: map[int]plugin.PluginSet{
			plugininterface.RPCVersion: {
				plugins.RPCPluginName: &plugins.GoPlugin{Impl: NewMultiplugin([]plugininterface.PluginRPC{
					// TODO: add concrete plugininterface.PluginRPC impls here
					&text.Plugin{},
				})},
			},
		},
	})
}
