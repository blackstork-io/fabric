package main

import (
	"github.com/blackstork-io/fabric/plugins"
	"github.com/blackstork-io/fabric/plugins/content"
	"github.com/blackstork-io/fabric/plugins/content/table"
	"github.com/blackstork-io/fabric/plugins/content/text"
	"github.com/blackstork-io/fabric/plugins/data"
	"github.com/blackstork-io/fabric/plugins/data/plugin_a"
	"github.com/blackstork-io/fabric/plugins/data/plugin_b"

	"github.com/hashicorp/go-plugin"
)

func main() {
	plugin.Serve(&plugin.ServeConfig{
		HandshakeConfig: plugins.Handshake,
		Plugins: plugin.PluginSet{
			"data.plugin_a": &data.GoPlugin{Impl: &plugin_a.Impl{}},
			"data.plugin_b": &data.GoPlugin{Impl: &plugin_b.Impl{}},
			"content.table": &content.GoPlugin{Impl: &table.Impl{}},
			"content.text":  &content.GoPlugin{Impl: &text.Impl{}},
		},
	})
}
