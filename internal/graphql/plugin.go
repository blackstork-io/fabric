package graphql

import (
	"github.com/blackstork-io/fabric/plugin"
)

func Plugin(version string) *plugin.Schema {
	return &plugin.Schema{
		Name:    "blackstork/graphql",
		Version: version,
		DataSources: plugin.DataSources{
			"graphql": makeGraphQLDataSource(),
		},
	}
}
