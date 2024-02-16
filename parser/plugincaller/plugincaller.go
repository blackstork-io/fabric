package plugincaller

import (
	"context"

	"github.com/blackstork-io/fabric/parser/evaluation"
	"github.com/blackstork-io/fabric/pkg/diagnostics"
	"github.com/blackstork-io/fabric/plugin"
)

type DataCaller interface {
	CallData(ctx context.Context, name string, config evaluation.Configuration, invocation evaluation.Invocation) (result plugin.MapData, diag diagnostics.Diag)
}

type ContentCaller interface {
	CallContent(ctx context.Context, name string, config evaluation.Configuration, invocation evaluation.Invocation, context plugin.MapData) (result string, diag diagnostics.Diag)
}

type PluginCaller interface {
	DataCaller
	ContentCaller
}
