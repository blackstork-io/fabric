package main

import (
	"weave-cli/plugins"
	"weave-cli/plugins/content"
	"weave-cli/plugins/content/table"
	"weave-cli/plugins/content/text"
	"weave-cli/plugins/data"
	"weave-cli/plugins/data/plugin_a"
	"weave-cli/plugins/data/plugin_b"

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
