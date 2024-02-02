package plugins

import (
	"github.com/hashicorp/go-plugin"

	"github.com/blackstork-io/fabric/plugininterface/v1"
)

const RPCPluginName = "FabricPlugin"

var Handshake = plugin.HandshakeConfig{
	MagicCookieKey:   "PLUGINS_FOR",
	MagicCookieValue: "fabric",
}

var PluginMap = map[int]plugin.PluginSet{
	plugininterface.RPCVersion: {
		RPCPluginName: &GoPlugin{},
	},
}
