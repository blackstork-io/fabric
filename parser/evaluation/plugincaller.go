package evaluation

import (
	"context"

	"github.com/blackstork-io/fabric/pkg/diagnostics"
	"github.com/blackstork-io/fabric/plugin"
)

type DataCaller interface {
	CallData(ctx context.Context, name string, config Configuration, invocation Invocation) (result plugin.MapData, diag diagnostics.Diag)
}

type ContentCaller interface {
	CallContent(ctx context.Context, name string, config Configuration, invocation Invocation, context plugin.MapData) (result string, diag diagnostics.Diag)
}

type PluginCaller interface {
	DataCaller
	ContentCaller
}
