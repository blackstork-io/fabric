package runner

import (
	"context"
	"fmt"
	"log/slog"
	"os"

	"github.com/hashicorp/hcl/v2"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"

	"github.com/blackstork-io/fabric/pkg/diagnostics"
	"github.com/blackstork-io/fabric/plugin"
	"github.com/blackstork-io/fabric/plugin/ast/nodes"
	pluginapiv1 "github.com/blackstork-io/fabric/plugin/pluginapi/v1"
)

type loadedPlugin struct {
	closefn func() error
	*plugin.Schema
}

type loadedDataSource struct {
	plugin *plugin.Schema
	*plugin.DataSource
}

type loadedContentProvider struct {
	plugin *plugin.Schema
	*plugin.ContentProvider
}

type loadedPublisher struct {
	plugin *plugin.Schema
	*plugin.Publisher
}

type loadedNodeRenderer struct {
	plugin   *plugin.Schema
	renderer plugin.NodeRendererFunc
}

type loader struct {
	logger          *slog.Logger
	tracer          trace.Tracer
	binaryMap       map[string]string
	builtin         *plugin.Schema
	pluginMap       map[string]loadedPlugin
	dataMap         map[string]loadedDataSource
	contentMap      map[string]loadedContentProvider
	publisherMap    map[string]loadedPublisher
	nodeRendererMap map[string]loadedNodeRenderer
}

func makeLoader(
	binaryMap map[string]string,
	builtin *plugin.Schema,
	logger *slog.Logger,
	tracer trace.Tracer,
) *loader {
	return &loader{
		tracer:          tracer,
		logger:          logger,
		binaryMap:       binaryMap,
		builtin:         builtin,
		pluginMap:       make(map[string]loadedPlugin),
		dataMap:         make(map[string]loadedDataSource),
		contentMap:      make(map[string]loadedContentProvider),
		publisherMap:    make(map[string]loadedPublisher),
		nodeRendererMap: make(map[string]loadedNodeRenderer),
	}
}

func nopCloser() error {
	return nil
}

func (l *loader) loadAll(ctx context.Context) (diags diagnostics.Diag) {
	ctx, span := l.tracer.Start(ctx, "loader.loadAll")
	l.logger.DebugContext(ctx, "Loading all plugins")
	defer func() {
		if diags.HasErrors() {
			span.RecordError(diags)
			span.SetStatus(codes.Error, diags.Error())
		} else {
			span.SetStatus(codes.Ok, "success")
		}
		span.End()
	}()
	diags = l.registerPlugin(ctx, l.builtin, nopCloser)
	if diags.HasErrors() {
		diags = append(diags, l.closeAll()...)
		return diags
	}
	for name, binaryPath := range l.binaryMap {
		if diags := l.loadBinary(ctx, name, binaryPath); diags.HasErrors() {
			diags = append(diags, l.closeAll()...)
			return diags
		}
	}
	return nil
}

func (l *loader) closeAll() diagnostics.Diag {
	var diags diagnostics.Diag
	for _, p := range l.pluginMap {
		if err := p.closefn(); err != nil {
			diags = append(diags, &hcl.Diagnostic{
				Severity: hcl.DiagWarning,
				Summary:  "Failed to close plugin",
				Detail:   fmt.Sprintf("Failed to close plugin %s@%s: %v", p.Name, p.Version, err),
			})
		}
	}
	return diags
}

func (l *loader) registerDataSource(ctx context.Context, name string, schema *plugin.Schema, ds *plugin.DataSource) diagnostics.Diag {
	l.logger.DebugContext(ctx, "Registering data source", "name", name, "plugin", schema.Name, "version", schema.Version)
	if found, has := l.dataMap[name]; has {
		return diagnostics.Diag{{
			Severity: hcl.DiagError,
			Summary:  "Duplicate data source",
			Detail:   fmt.Sprintf("Data source %s provided by plugin %s@%s and %s@%s", name, schema.Name, schema.Version, found.plugin.Name, found.plugin.Version),
		}}
	}
	l.dataMap[name] = loadedDataSource{
		plugin:     schema,
		DataSource: ds,
	}
	return nil
}

func (l *loader) registerContentProvider(ctx context.Context, name string, schema *plugin.Schema, cp *plugin.ContentProvider) diagnostics.Diag {
	l.logger.DebugContext(ctx, "Registering content provider", "name", name, "plugin", schema.Name, "version", schema.Version)
	if found, has := l.contentMap[name]; has {
		return diagnostics.Diag{{
			Severity: hcl.DiagError,
			Summary:  "Duplicate content provider",
			Detail:   fmt.Sprintf("Content provider %s provided by plugin %s@%s and %s@%s", name, schema.Name, schema.Version, found.plugin.Name, found.plugin.Version),
		}}
	}
	l.contentMap[name] = loadedContentProvider{
		plugin: schema,
		ContentProvider: &plugin.ContentProvider{
			Config: cp.Config,
			Args:   cp.Args,
			Doc:    cp.Doc,
			Tags:   cp.Tags,
			ContentFunc: func(ctx context.Context, params *plugin.ProvideContentParams) (*plugin.ContentElement, diagnostics.Diag) {
				return schema.ProvideContent(ctx, name, params)
			},
		},
	}
	return nil
}

func (l *loader) registerPublisher(ctx context.Context, name string, schema *plugin.Schema, pub *plugin.Publisher) diagnostics.Diag {
	l.logger.DebugContext(ctx, "Registering publisher", "name", name, "plugin", schema.Name, "version", schema.Version)
	if found, has := l.publisherMap[name]; has {
		return diagnostics.Diag{{
			Severity: hcl.DiagError,
			Summary:  "Duplicate publisher",
			Detail:   fmt.Sprintf("Publisher %s provided by plugin %s@%s and %s@%s", name, schema.Name, schema.Version, found.plugin.Name, found.plugin.Version),
		}}
	}
	l.publisherMap[name] = loadedPublisher{
		plugin:    schema,
		Publisher: pub,
	}
	return nil
}

func (l *loader) registerNodeRenderer(ctx context.Context, nodeType string, schema *plugin.Schema, renderer plugin.NodeRendererFunc) diagnostics.Diag {
	l.logger.DebugContext(ctx, "Registering node renderer", "name", nodeType, "plugin", schema.Name, "version", schema.Version)
	fqn := nodes.FullyQualifiedCustomNodeTypeURL(schema.Name, nodeType)
	if found, has := l.publisherMap[fqn]; has {
		return diagnostics.Diag{{
			Severity: hcl.DiagError,
			Summary:  "Duplicate node renderer",
			Detail:   fmt.Sprintf("Node renderer %s provided by plugin %s@%s and %s@%s", fqn, schema.Name, schema.Version, found.plugin.Name, found.plugin.Version),
		}}
	}

	l.nodeRendererMap[fqn] = loadedNodeRenderer{
		plugin:   schema,
		renderer: renderer,
	}
	return nil
}

func (l *loader) registerPlugin(ctx context.Context, schema *plugin.Schema, closefn func() error) (diags diagnostics.Diag) {
	l.logger.DebugContext(ctx, "Registering a plugin", "name", schema.Name, "version", schema.Version)
	if diags.Extend(schema.Validate()) {
		l.logger.ErrorContext(ctx, "Validation errors while registering a plugin", "name", schema.Name, "version", schema.Version)
		return
	}
	schema = plugin.WithLogging(schema, l.logger)
	schema = plugin.WithTracing(schema, l.tracer)
	if found, has := l.pluginMap[schema.Name]; has {
		diags.Append(&hcl.Diagnostic{
			Severity: hcl.DiagError,
			Summary:  fmt.Sprintf("Plugin %s conflict", schema.Name),
			Detail:   fmt.Sprintf("%s@%s and %s@%s have the same schema name", schema.Name, schema.Version, found.Name, found.Version),
		})
		if err := closefn(); err != nil {
			diags.AppendErr(err, fmt.Sprintf("Failed to close plugin %s@%s", found.Name, found.Version))
		}
		return
	}
	plugin := loadedPlugin{
		closefn: closefn,
		Schema:  schema,
	}
	l.pluginMap[schema.Name] = plugin
	for name, source := range schema.DataSources {
		if diags.Extend(l.registerDataSource(ctx, name, schema, source)) {
			return
		}
	}
	for name, provider := range schema.ContentProviders {
		if diags.Extend(l.registerContentProvider(ctx, name, schema, provider)) {
			return
		}
	}
	for name, publisher := range schema.Publishers {
		if diags.Extend(l.registerPublisher(ctx, name, schema, publisher)) {
			return
		}
	}
	for nodeType, renderer := range schema.NodeRenderers {
		if diags.Extend(l.registerNodeRenderer(ctx, nodeType, schema, renderer)) {
			return
		}
	}
	return
}

func (l *loader) loadBinary(ctx context.Context, name, binaryPath string) (diags diagnostics.Diag) {
	ctx, span := l.tracer.Start(ctx, "loader.loadBinary", trace.WithAttributes(
		attribute.String("name", name),
	))
	l.logger.InfoContext(ctx, "Loading plugin", "name", name, "path", binaryPath)
	defer func() {
		if diags.HasErrors() {
			span.RecordError(diags)
			span.SetStatus(codes.Error, diags.Error())
		} else {
			span.SetStatus(codes.Ok, "success")
		}
		span.End()
	}()
	if info, err := os.Stat(binaryPath); os.IsNotExist(err) {
		return diagnostics.Diag{{
			Severity: hcl.DiagError,
			Summary:  fmt.Sprintf("Plugin %s binary not found", name),
			Detail:   fmt.Sprintf("Executable not found at: %s", binaryPath),
		}}
	} else if info.IsDir() {
		return diagnostics.Diag{{
			Severity: hcl.DiagError,
			Summary:  fmt.Sprintf("Plugin %s binary path is a directory", name),
			Detail:   fmt.Sprintf("Path %s is a directory", binaryPath),
		}}
	}
	p, close, err := pluginapiv1.NewClient(name, binaryPath, l.logger)
	if err != nil {
		return diagnostics.Diag{{
			Severity: hcl.DiagError,
			Summary:  fmt.Sprintf("Failed to load plugin %s", name),
			Detail:   err.Error(),
		}}
	}
	return l.registerPlugin(ctx, p, close)
}
