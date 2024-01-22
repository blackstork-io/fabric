package plugins

import (
	"github.com/hashicorp/go-plugin"

	"github.com/blackstork-io/fabric/pkg/utils"
)

// Stub implementation of plugin caller
// TODO: attach to plugin discovery mechanism

type Plugins struct {
	content PluginType
	data    PluginType
	client  *plugin.Client
}

type PluginType struct {
	plugins map[string]PluginRPCIntefrace
	Names   func() string
}

func NewPluginType(plugins map[string]PluginRPCIntefrace) PluginType {
	return PluginType{
		plugins: plugins,
		Names:   utils.MemoizedKeys(&plugins),
	}
}

type PluginRPCIntefrace interface{}
