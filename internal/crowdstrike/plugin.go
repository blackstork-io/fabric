package crowdstrike

import (
	"github.com/blackstork-io/fabric/plugin"
)

func Plugin(version string) *plugin.Schema {
	return &plugin.Schema{
		Name:        "blackstork/crowdstrike",
		Version:     version,
		DataSources: plugin.DataSources{},
	}
}
