package resolver

import (
	"io"
	"log/slog"
)

// options for the resolver.
type options struct {
	logger  *slog.Logger
	sources []Source
}

var defaultOptions = options{
	logger:  slog.New(slog.NewTextHandler(io.Discard, nil)),
	sources: []Source{},
}

// Option is a functional option for the resolver.
type Option func(*options)

// WithLogger sets the logger for the resolver.
func WithLogger(logger *slog.Logger) Option {
	return func(o *options) {
		o.logger = logger
	}
}

// WithSources sets the sources for the resolver.
func WithSources(sources ...Source) Option {
	return func(o *options) {
		o.sources = sources
	}
}
