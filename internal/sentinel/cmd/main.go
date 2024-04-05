package main

import (
	"github.com/blackstork-io/fabric/internal/sentinel"
	pluginapiv1 "github.com/blackstork-io/fabric/plugin/pluginapi/v1"
)

var version string

func main() {
	pluginapiv1.Serve(
		sentinel.Plugin(version, sentinel.DefaultClientLoader),
	)
}
