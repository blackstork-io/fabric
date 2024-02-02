package main

import (
	"github.com/hashicorp/go-plugin"

	"github.com/blackstork-io/fabric/plugininterface/v1"
	"github.com/blackstork-io/fabric/plugins"
)

func ServePlugins(pluginImpls ...plugininterface.PluginRPC) {
	plugin.Serve(&plugin.ServeConfig{
		HandshakeConfig: plugins.Handshake,
		VersionedPlugins: map[int]plugin.PluginSet{
			plugininterface.RPCVersion: {
				plugins.RPCPluginName: &plugins.GoPlugin{Impl: NewMultiplugin(pluginImpls)},
			},
		},
	})
}
