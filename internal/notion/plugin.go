package notion

import (
	"log/slog"

	"github.com/blackstork-io/fabric/plugin"
	"go.opentelemetry.io/otel/trace"
)

func Plugin(version string, logger *slog.Logger, tracer trace.Tracer) *plugin.Schema {
	return &plugin.Schema{
		Name:    "blackstork/notion",
		Version: version,
		DataSources: plugin.DataSources{
			"md": makeMarkdownDataSource(),
		},
		Publishers: plugin.Publishers{
			"notion_page": makeNotionPagePublisher(logger, tracer),
		},
	}
}
