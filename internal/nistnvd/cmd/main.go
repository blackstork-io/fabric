package main

import (
	"github.com/blackstork-io/fabric/internal/nistnvd"
	pluginapiv1 "github.com/blackstork-io/fabric/plugin/pluginapi/v1"
)

var version string

func main() {
	pluginapiv1.Serve(
		nistnvd.Plugin(version, nistnvd.DefaultClientLoader),
	)
}
