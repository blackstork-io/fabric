package print

import (
	"context"
	"io"
	"log/slog"

	"github.com/blackstork-io/fabric/plugin/ast/nodes"
)

type logging struct {
	next   Printer
	logger *slog.Logger
	attrs  []slog.Attr
}

// WithLogging wraps the printer with logging instrumentation.
func WithLogging(next Printer, logger *slog.Logger, attrs ...slog.Attr) Printer {
	return &logging{
		next:   next,
		logger: logger,
		attrs:  attrs,
	}
}

func (p logging) Print(ctx context.Context, w io.Writer, el *nodes.Node) (err error) {
	p.logger.LogAttrs(ctx, slog.LevelDebug, "Printing content", p.attrs...)
	return p.next.Print(ctx, w, el)
}
