package plugins

import (
	"github.com/blackstork-io/fabric/pkg/utils"
)

// Stub implementation of plugin caller
// TODO: attach to plugin discovery mechanism

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
