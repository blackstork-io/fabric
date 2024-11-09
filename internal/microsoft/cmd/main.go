package main

import (
	"github.com/blackstork-io/fabric/internal/microsoft"
	"github.com/blackstork-io/fabric/internal/microsoft/client"
	pluginapiv1 "github.com/blackstork-io/fabric/plugin/pluginapi/v1"
)

var version string

func main() {
	pluginapiv1.Serve(
		microsoft.Plugin(
			version,
			microsoft.MakeDefaultAzureClientLoader(client.AcquireAzureToken),
			client.MakeAzureOpenAIClientLoader(),
			microsoft.MakeDefaultMicrosoftGraphClientLoader(client.AcquireAzureToken),
			microsoft.MakeDefaultMicrosoftSecurityClientLoader(client.AcquireAzureToken),
		),
	)
}
