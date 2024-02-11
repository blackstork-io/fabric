package pluginapiv1

import (
	"fmt"
	"log/slog"
	"os/exec"
	"path/filepath"
	"strings"

	goplugin "github.com/hashicorp/go-plugin"

	"github.com/blackstork-io/fabric/pkg/sloghclog"
	"github.com/blackstork-io/fabric/plugin"
)

func parsePluginInfo(path string) (name, version string, err error) {
	nameVer := filepath.Base(path)
	ext := filepath.Ext(path)

	parts := strings.SplitN(
		nameVer[:len(nameVer)-len(ext)],
		"@", 2,
	)
	if len(parts) != 2 {
		err = fmt.Errorf("plugin at '%s' must have a file name '<plugin_name>@<plugin_version>[.exe]'", path)
		return
	}
	name = parts[0]
	version = parts[1]
	return
}

func NewClient(loc string) (p *plugin.Schema, closefn func() error, err error) {
	pluginName, _, err := parsePluginInfo(loc)
	if err != nil {
		return
	}
	slog.Info("Loading plugin", "filename", loc)
	client := goplugin.NewClient(&goplugin.ClientConfig{
		HandshakeConfig: handshake,
		Plugins: map[string]goplugin.Plugin{
			pluginName: &grpcPlugin{},
		},
		Cmd: exec.Command(loc),
		AllowedProtocols: []goplugin.Protocol{
			goplugin.ProtocolGRPC,
		},
		Logger: sloghclog.Adapt(
			slog.Default(),
			sloghclog.Name("plugin."+pluginName),
			// disable code location reporting, it's always going to be incorrect
			// for remote plugin logs
			sloghclog.AddSource(false),
			sloghclog.Level(slog.LevelInfo), // debug is too noisy for plugins
		),
	})
	rpcClient, err := client.Client()
	if err != nil {
		return nil, nil, fmt.Errorf("failed to create plugin client: %v", err)
	}
	raw, err := rpcClient.Dispense(pluginName)
	if err != nil {
		rpcClient.Close()
		return nil, nil, fmt.Errorf("failed to dispense plugin: %v", err)
	}
	plg, ok := raw.(*plugin.Schema)
	if !ok {
		rpcClient.Close()
		return nil, nil, fmt.Errorf("unexpected plugin type: %T", raw)
	}
	return plg, rpcClient.Close, nil
}
