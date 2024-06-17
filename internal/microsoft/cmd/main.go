package main

import (
	"github.com/blackstork-io/fabric/internal/microsoft"
	pluginapiv1 "github.com/blackstork-io/fabric/plugin/pluginapi/v1"
)

var version string

func main() {
	pluginapiv1.Serve(
		microsoft.Plugin(version, microsoft.DefaultClientLoader),
	)
}
