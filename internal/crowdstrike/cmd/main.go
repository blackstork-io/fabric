package main

import (
	"github.com/blackstork-io/fabric/internal/crowdstrike"
	pluginapiv1 "github.com/blackstork-io/fabric/plugin/pluginapi/v1"
)

var version string

func main() {
	pluginapiv1.Serve(
		crowdstrike.Plugin(version, nil),
	)
}
