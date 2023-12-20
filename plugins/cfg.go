package plugins

import (
	"github.com/blackstork-io/fabric/plugins/content"
	"github.com/blackstork-io/fabric/plugins/data"

	"github.com/hashicorp/go-plugin"
)

var Handshake = plugin.HandshakeConfig{
	ProtocolVersion:  1,
	MagicCookieKey:   "PLUGINS_FOR",
	MagicCookieValue: "fabric",
}

var PluginMap = plugin.PluginSet{
	"data.plugin_a": &data.GoPlugin{},
	"data.plugin_b": &data.GoPlugin{},
	"content.table": &content.GoPlugin{},
	"content.text":  &content.GoPlugin{},
}
