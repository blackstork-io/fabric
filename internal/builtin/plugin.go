package builtin

import (
	"github.com/blackstork-io/fabric/plugin"
)

const Name = "blackstork/builtin"

func Plugin(version string) *plugin.Schema {
	return &plugin.Schema{
		Name:    Name,
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
			"title":       makeTitleContentProvider(),
			"code":        makeCodeContentProvider(),
			"blockquote":  makeBlockQuoteContentProvider(),
			"image":       makeImageContentProvider(),
			"list":        makeListContentProvider(),
			"table":       makeTableContentProvider(),
			"frontmatter": makeFrontMatterContentProvider(),
		},
	}
}
