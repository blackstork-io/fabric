package builtin

import (
	"github.com/blackstork-io/fabric/plugin"
)

func Plugin(version string) *plugin.Schema {
	return &plugin.Schema{
		Name:    "builtin",
		Version: version,
		DataSources: plugin.DataSources{
			"csv":    makeCSVDataSource(),
			"txt":    makeTXTDataSource(),
			"json":   makeJSONDataSource(),
			"inline": makeInlineDataSource(),
		},
		ContentProviders: plugin.ContentProviders{
			"text":        makeTextContentProvider(),
			"image":       makeImageContentProvider(),
			"list":        makeListContentProvider(),
			"table":       makeTableContentProvider(),
			"frontmatter": makeFrontMatterContentProvider(),
		},
	}
}