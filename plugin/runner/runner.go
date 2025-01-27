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
	pluginMap       map[string]loadedPlugin
	dataMap         map[string]loadedDataSource
	contentMap      map[string]loadedContentProvider
	publisherMap    map[string]loadedPublisher
	nodeRendererMap map[string]loadedNodeRenderer
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
		pluginMap:       loader.pluginMap,
		dataMap:         loader.dataMap,
		contentMap:      loader.contentMap,
		publisherMap:    loader.publisherMap,
		nodeRendererMap: loader.nodeRendererMap,
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

func (m *Runner) NodeRenderer(customNodeType string) (plugin.NodeRendererFunc, bool) {
	nodeRenderer, ok := m.nodeRendererMap[customNodeType]
	if !ok {
		return nil, false
	}
	return nodeRenderer.renderer, true
}

func (m *Runner) AllNodeRenderers() map[string]struct{} {
	renderers := make(map[string]struct{}, len(m.nodeRendererMap))
	for name := range m.nodeRendererMap {
		renderers[name] = struct{}{}
	}
	return renderers
}

func (m *Runner) Close() (diags diagnostics.Diag) {
	for _, p := range m.pluginMap {
		if err := p.closefn(); err != nil {
			diags.AppendErr(
				err,
				fmt.Sprintf("Failed to close plugin '%s'", p.Name),
			)
		}
	}
	return diags.Refine(diagnostics.OverrideSeverity(hcl.DiagWarning))
}
