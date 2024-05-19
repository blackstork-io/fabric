package print

import (
	"context"
	"io"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"

	"github.com/blackstork-io/fabric/plugin"
)

type tracing struct {
	next   Printer
	tracer trace.Tracer
	attrs  []attribute.KeyValue
}

// WithTracing wraps a printer with tracing instrumentation.
func WithTracing(next Printer, tracer trace.Tracer, attrs ...attribute.KeyValue) Printer {
	return tracing{
		next:   next,
		tracer: tracer,
		attrs:  attrs,
	}
}

func (p tracing) Print(ctx context.Context, w io.Writer, el plugin.Content) (err error) {
	ctx, span := p.tracer.Start(ctx, "Printer.Print", trace.WithAttributes(p.attrs...))
	defer func() {
		if err != nil {
			span.RecordError(err)
			span.SetStatus(codes.Error, err.Error())
		} else {
			span.SetStatus(codes.Ok, "success")
		}
		span.End()
	}()
	return p.next.Print(ctx, w, el)
}
