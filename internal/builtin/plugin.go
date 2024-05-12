package builtin

import (
	"log/slog"

	"go.opentelemetry.io/otel/trace"

	"github.com/blackstork-io/fabric/plugin"
)

const Name = "blackstork/builtin"

func Plugin(version string, logger *slog.Logger, tracer trace.Tracer) *plugin.Schema {
	return &plugin.Schema{
		Name:    Name,
		Version: version,
		DataSources: plugin.DataSources{
			"csv":    makeCSVDataSource(),
			"txt":    makeTXTDataSource(),
			"rss":    makeRSSDataSource(),
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
		Publishers: plugin.Publishers{
			"local_file": makeLocalFilePublisher(logger, tracer),
		},
	}
}
