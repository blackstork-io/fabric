package stixview

import (
	"github.com/blackstork-io/fabric/plugin"
)

func Plugin(version string) *plugin.Schema {
	return &plugin.Schema{
		Name:    "blackstork/stixview",
		Version: version,
		ContentProviders: plugin.ContentProviders{
			"stixview": makeStixViewContentProvider(),
		},
	}
}
