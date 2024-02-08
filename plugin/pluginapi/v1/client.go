package pluginapiv1

import (
	"fmt"
	"log/slog"
	"os/exec"
	"path"
	"strings"

	goplugin "github.com/hashicorp/go-plugin"

	"github.com/blackstork-io/fabric/pkg/sloghclog"
	"github.com/blackstork-io/fabric/plugin"
)

func NewClient(loc string) (p *plugin.Schema, closefn func() error, err error) {
	base := path.Base(loc)
	if base == "" {
		return nil, nil, fmt.Errorf("invalid plugin location: %s", loc)
	}
	parts := strings.SplitN(base, "@", 2)
	if len(parts) != 2 {
		return nil, nil, fmt.Errorf("invalid plugin name: %s", base)
	}
	client := goplugin.NewClient(&goplugin.ClientConfig{
		HandshakeConfig: handshake,
		Plugins: map[string]goplugin.Plugin{
			parts[0]: &grpcPlugin{},
		},
		Cmd: exec.Command("sh", "-c", loc),
		AllowedProtocols: []goplugin.Protocol{
			goplugin.ProtocolGRPC,
		},
		Logger: sloghclog.Adapt(
			slog.Default(),
			sloghclog.Name("plugin."+parts[0]),
			// disable code location reporting, it's always going to be incorrect
			// for remote plugin logs
			sloghclog.AddSource(false),
		),
	})
	rpcClient, err := client.Client()
	if err != nil {
		return nil, nil, fmt.Errorf("failed to create plugin client: %v", err)
	}
	raw, err := rpcClient.Dispense(parts[0])
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
