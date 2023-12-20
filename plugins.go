package main

import (
	"fmt"
	"io"
	"os/exec"
	"slices"
	"strings"
	"sync"
	"weave-cli/pkg/diagnostics"
	"weave-cli/plugins"
	plugContent "weave-cli/plugins/content"

	"github.com/hashicorp/go-hclog"
	"github.com/hashicorp/go-plugin"
	"github.com/hashicorp/hcl/v2"
	"golang.org/x/exp/maps"
)

type Plugins struct {
	content PluginType
	data    PluginType
	client  *plugin.Client
}

type PluginType struct {
	plugins map[string]any
	Names   func() string
}

func NewPluginType(plugins map[string]any) PluginType {
	return PluginType{
		plugins: plugins,
		Names:   memoizedKeys(&plugins),
	}
}

func (p *Plugins) ByKind(kind string) *PluginType {
	switch kind {
	case BK_CONTENT:
		return &p.content
	case BK_DATA:
		return &p.data
	default:
		panic(fmt.Errorf("unknown plugin kind '%s'", kind))
	}
}

func memoizedKeys[M ~map[string]V, V any](m *M) func() string {
	return sync.OnceValue(func() string {
		keys := maps.Keys(*m)
		slices.Sort(keys)
		return JoinSurround(", ", "'", keys...)
	})
}

type genericPlugin struct{}

// Execute implements content.Plugin.
func (*genericPlugin) Execute(attrs, content any) (string, error) {
	return "", nil
}

var _ plugContent.Plugin = (*genericPlugin)(nil)

func NewPlugins(pluginPath string) (p *Plugins, diag diagnostics.Diagnostics) {
	// TODO: setup pluggin logging?
	hclog.DefaultOutput = io.Discard
	client := plugin.NewClient(&plugin.ClientConfig{
		HandshakeConfig: plugins.Handshake,
		Plugins:         plugins.PluginMap,
		Cmd:             exec.Command(pluginPath),
		// Logger:          hclog.de,
	})
	defer func() {
		if diag.HasErrors() {
			client.Kill()
		}
	}()

	// Connect via RPC
	rpcClient, err := client.Client()
	if diag.FromErr(err, "Plugin connection error") {
		return
	}

	content := map[string]any{
		"generic": &genericPlugin{},
	}
	data := map[string]any{}

	for pluginName := range plugins.PluginMap {
		split := strings.SplitN(pluginName, ".", 2)
		if len(split) != 2 {
			diag.Append(&hcl.Diagnostic{
				Severity: hcl.DiagError,
				Summary:  "Invalid plugin name",
				Detail:   fmt.Sprintf("invalid name format for plugin '%s': missing dot", pluginName),
			})
			return
		}
		var tgtMap map[string]any
		switch split[0] {
		case BK_CONTENT:
			tgtMap = content
		case BK_DATA:
			tgtMap = data
		default:
			diag.Append(&hcl.Diagnostic{
				Severity: hcl.DiagError,
				Summary:  "Invalid plugin name",
				Detail:   fmt.Sprintf("invalid name for plugin '%s': '%s' is an invalid plugin kind", pluginName, split[0]),
			})
			return
		}
		// Request the plugin
		var rawPlugin any
		rawPlugin, err = rpcClient.Dispense(pluginName)
		if diag.FromErr(err, "Plugin RPC error") {
			return
		}
		tgtMap[split[1]] = rawPlugin
	}

	p = &Plugins{
		content: NewPluginType(content),
		data:    NewPluginType(data),
	}
	return
}

func (p *Plugins) Kill() {
	if p != nil && p.client != nil {
		p.client.Kill()
	}
}
