package builtin

import (
	"github.com/blackstork-io/fabric/plugin"
)

func Plugin(version string) *plugin.Schema {
	return &plugin.Schema{
		Name:    "blackstork/builtin",
		Version: version,
		DataSources: plugin.DataSources{
			"csv":    makeCSVDataSource(),
			"txt":    makeTXTDataSource(),
			"json":   makeJSONDataSource(),
			"inline": makeInlineDataSource(),
		},
		ContentProviders: plugin.ContentProviders{
			"toc":         makeTOCContentProvider(),
			"text":        makeTextContentProvider(),
			"image":       makeImageContentProvider(),
			"list":        makeListContentProvider(),
			"table":       makeTableContentProvider(),
			"frontmatter": makeFrontMatterContentProvider(),
		},
	}
}
