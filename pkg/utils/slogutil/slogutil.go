package slogutil

import (
	"context"
	"log/slog"
	"runtime"
)

// SourceRewriter is a slog.Handler that rewrites the source of log entries.
type SourceRewriter struct {
	slog.Handler
}

// WithAttrs implements slog.Handler.
func (sr SourceRewriter) WithAttrs(attrs []slog.Attr) slog.Handler {
	return SourceRewriter{Handler: sr.Handler.WithAttrs(attrs)}
}

// WithGroup implements slog.Handler.
func (sr SourceRewriter) WithGroup(name string) slog.Handler {
	return SourceRewriter{Handler: sr.Handler.WithGroup(name)}
}

// NewSourceRewriter returns a new SourceRewriter that wraps the given handler.
// The returned SourceRewriter applies the SourceOverride commands.
func NewSourceRewriter(h slog.Handler) SourceRewriter {
	return SourceRewriter{Handler: h}
}

const sourceKey = "Source Override. Use SourceRewriter handler to apply."

func (sr SourceRewriter) Handle(ctx context.Context, r slog.Record) error {
	var attrCountAfterOverride int
	r.Attrs(func(a slog.Attr) bool {
		if !(a.Key == sourceKey && a.Value.Kind() == slog.KindUint64) {
			attrCountAfterOverride++
		}
		return true
	})
	if attrCountAfterOverride == r.NumAttrs() {
		return sr.Handler.Handle(ctx, r)
	}
	attrs := make([]slog.Attr, 0, attrCountAfterOverride)
	newRecord := slog.Record{
		Time:    r.Time,
		Message: r.Message,
		Level:   r.Level,
		PC:      r.PC,
	}
	r.Attrs(func(a slog.Attr) bool {
		if a.Key == sourceKey && a.Value.Kind() == slog.KindUint64 {
			newRecord.PC = uintptr(a.Value.Uint64())
		} else {
			attrs = append(attrs, a)
		}
		return true
	})
	newRecord.AddAttrs(attrs...)
	return sr.Handler.Handle(ctx, newRecord)
}

var _ slog.Handler = SourceRewriter{}

// SourceOverride returns a slog.Attr with the source of the caller at the given offset.
// SourceOverride(0) is the caller of SourceOverride, SourceOverride(1) is the caller of the caller, etc.
func SourceOverride(offset int) slog.Attr {
	var pcs [1]uintptr
	// skip [runtime.Callers, this function, +offset functions]
	runtime.Callers(2+offset, pcs[:])
	return slog.Uint64(sourceKey, uint64(pcs[0]))
}
