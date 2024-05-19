package runner

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/hashicorp/hcl/v2"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"

	"github.com/blackstork-io/fabric/pkg/diagnostics"
	"github.com/blackstork-io/fabric/plugin"
)

type Runner struct {
	pluginMap    map[string]loadedPlugin
	dataMap      map[string]loadedDataSource
	contentMap   map[string]loadedContentProvider
	publisherMap map[string]loadedPublisher
}

func Load(
	ctx context.Context,
	binaryMap map[string]string,
	builtin *plugin.Schema,
	logger *slog.Logger,
	tracer trace.Tracer,
) (_ *Runner, diags diagnostics.Diag) {
	ctx, span := tracer.Start(ctx, "runner.Load")
	defer func() {
		if diags.HasErrors() {
			span.RecordError(diags)
			span.SetStatus(codes.Error, diags.Error())
		}
		span.End()
	}()
	logger = logger.With("component", "runner")
	logger.DebugContext(ctx, "Loading plugins")
	loader := makeLoader(binaryMap, builtin, logger, tracer)
	if diags = loader.loadAll(ctx); diags.HasErrors() {
		return nil, diags
	}
	return &Runner{
		pluginMap:    loader.pluginMap,
		dataMap:      loader.dataMap,
		contentMap:   loader.contentMap,
		publisherMap: loader.publisherMap,
	}, nil
}

func (m *Runner) DataSource(name string) (*plugin.DataSource, bool) {
	source, ok := m.dataMap[name]
	if !ok {
		return nil, false
	}
	return source.DataSource, true
}

func (m *Runner) ContentProvider(name string) (*plugin.ContentProvider, bool) {
	provider, ok := m.contentMap[name]
	if !ok {
		return nil, false
	}
	return provider.ContentProvider, true
}

func (m *Runner) Publisher(name string) (*plugin.Publisher, bool) {
	publisher, ok := m.publisherMap[name]
	if !ok {
		return nil, false
	}
	return publisher.Publisher, true
}

func (m *Runner) Close() diagnostics.Diag {
	var diags diagnostics.Diag
	for _, p := range m.pluginMap {
		if err := p.closefn(); err != nil {
			diags = append(diags, &hcl.Diagnostic{
				Severity: hcl.DiagWarning,
				Summary:  fmt.Sprintf("Failed to close plugin '%s'", p.Name),
				Detail:   err.Error(),
			})
		}
	}
	return diags
}
