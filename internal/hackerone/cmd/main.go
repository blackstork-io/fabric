package main

import (
	"github.com/blackstork-io/fabric/internal/hackerone"
	pluginapiv1 "github.com/blackstork-io/fabric/plugin/pluginapi/v1"
)

var version string

func main() {
	pluginapiv1.Serve(
		hackerone.Plugin(version, hackerone.DefaultClientLoader),
	)
}
