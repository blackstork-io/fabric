package engine

import (
	"context"
	"fmt"
	"io"
	"io/fs"
	"log/slog"
	"path/filepath"
	"slices"
	"strings"

	"github.com/hashicorp/hcl/v2"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"

	"github.com/blackstork-io/fabric/eval"
	"github.com/blackstork-io/fabric/parser"
	"github.com/blackstork-io/fabric/parser/definitions"
	"github.com/blackstork-io/fabric/pkg/diagnostics"
	"github.com/blackstork-io/fabric/plugin"
	"github.com/blackstork-io/fabric/plugin/resolver"
	"github.com/blackstork-io/fabric/plugin/runner"
)

// Engine is the main entry point for the fabric engine. It is responsible for
// loading and evaluating fabric files, installing plugins, and fetching data.
// It is also responsible for managing the plugin resolver and runner.
type Engine struct {
	builtin  *plugin.Schema
	logger   *slog.Logger
	tracer   trace.Tracer
	config   *definitions.GlobalConfig
	blocks   *parser.DefinedBlocks
	runner   *runner.Runner
	lockFile *resolver.LockFile
	resolver *resolver.Resolver
	fileMap  map[string]*hcl.File
}

// New creates a new Engine instance with the provided options.
func New(options ...Option) *Engine {
	opts := defaultOptions
	for _, opt := range options {
		opt(&opts)
	}
	return &Engine{
		builtin: opts.builtin,
		logger:  opts.logger,
		tracer:  opts.tracer,
		config: &definitions.GlobalConfig{
			PluginRegistry: &definitions.PluginRegistry{
				BaseURL:   opts.registryBaseURL,
				MirrorDir: "",
			},
			CacheDir: opts.cacheDir,
		},
	}
}

func (e *Engine) PluginResolver() *resolver.Resolver {
	return e.resolver
}

func (e *Engine) PluginRunner() *runner.Runner {
	return e.runner
}

func (e *Engine) LockFile() *resolver.LockFile {
	return e.lockFile
}

func (e *Engine) FileMap() map[string]*hcl.File {
	return e.fileMap
}

func (e *Engine) Install(ctx context.Context, upgrade bool) (diags diagnostics.Diag) {
	ctx, span := e.tracer.Start(ctx, "Engine.Install", trace.WithAttributes(
		attribute.Bool("upgrade", upgrade),
	))
	e.logger.InfoContext(ctx, "Installing plugins", "upgrade", upgrade)
	defer func() {
		if diags.HasErrors() {
			span.RecordError(diags)
			span.SetStatus(codes.Error, diags.Error())
		}
		span.End()
	}()
	if e.resolver == nil {
		diags.Append(&hcl.Diagnostic{
			Severity: hcl.DiagError,
			Summary:  "Plugin resolver is not loaded",
			Detail:   "Load plugin resolver before installing",
		})
		return
	}
	lockFile, diag := e.resolver.Install(ctx, e.lockFile, upgrade)
	if diags.Extend(diag) {
		return
	}
	e.lockFile = lockFile
	err := resolver.SaveLockFileTo(defaultLockFile, lockFile)
	if err != nil {
		diags.Append(&hcl.Diagnostic{
			Severity: hcl.DiagError,
			Summary:  "Failed to save lock file",
			Detail:   err.Error(),
		})
	}
	return
}

func (e *Engine) ParseDir(ctx context.Context, sourceDir fs.FS) (diags diagnostics.Diag) {
	ctx, span := e.tracer.Start(ctx, "Engine.ParseDir")
	e.logger.InfoContext(ctx, "Parsing fabric files")
	defer func() {
		if diags.HasErrors() {
			span.RecordError(diags)
			span.SetStatus(codes.Error, diags.Error())
		}
		span.End()
	}()
	e.blocks, e.fileMap, diags = parser.ParseDir(sourceDir)
	if diags.HasErrors() {
		return
	}
	if e.blocks.GlobalConfig != nil {
		e.config.Merge(e.blocks.GlobalConfig)
	}
	return
}

func (e *Engine) Lint(ctx context.Context, fullLint bool) (diags diagnostics.Diag) {
	ctx, span := e.tracer.Start(ctx, "Engine.Lint", trace.WithAttributes(
		attribute.Bool("fullLint", fullLint),
	))
	e.logger.InfoContext(ctx, "Linting all documents", "full_lint", fullLint)
	defer func() {
		if diags.HasErrors() {
			span.RecordError(diags)
			span.SetStatus(codes.Error, diags.Error())
		}
		span.End()
	}()
	for _, doc := range e.blocks.Documents {
		e.logger.DebugContext(ctx, "Linting document", "document", doc.Name)
		parsedDoc, diag := e.blocks.ParseDocument(doc)
		diags.Extend(diag)
		if fullLint {
			_, diag = eval.LoadDocument(e.runner, parsedDoc)
			diags.Extend(diag)
		}
	}
	return diags
}

func (e *Engine) LoadPluginResolver(ctx context.Context, includeRemote bool) (diags diagnostics.Diag) {
	ctx, span := e.tracer.Start(ctx, "Engine.LoadPluginResolver", trace.WithAttributes(
		attribute.String("includeRemote", fmt.Sprint(includeRemote)),
	))
	e.logger.InfoContext(ctx, "Loading plugin resolver", "includeRemote", includeRemote)
	defer func() {
		if diags.HasErrors() {
			span.RecordError(diags)
			span.SetStatus(codes.Error, diags.Error())
		}
		span.End()
	}()
	pluginDir := filepath.Join(e.config.CacheDir, "plugins")
	sources := []resolver.Source{
		resolver.NewLocal(pluginDir, e.logger, e.tracer),
	}
	if e.config.PluginRegistry != nil {
		if e.config.PluginRegistry.MirrorDir != "" {
			sources = append(sources, resolver.NewLocal(e.config.PluginRegistry.MirrorDir, e.logger, e.tracer))
		}
		if includeRemote && e.config.PluginRegistry.BaseURL != "" {
			sources = append(sources, resolver.NewRemote(resolver.RemoteOptions{
				BaseURL:     e.config.PluginRegistry.BaseURL,
				DownloadDir: pluginDir,
				UserAgent:   fmt.Sprintf("fabric/%s", "version"),
			}))
		}
	}
	var err error
	e.lockFile, err = resolver.ReadLockFileFrom(defaultLockFile)
	if err != nil {
		return diagnostics.Diag{{
			Severity: hcl.DiagError,
			Summary:  "Failed to read lock file",
			Detail:   err.Error(),
		}}
	}
	resolve, diags := resolver.NewResolver(e.config.PluginVersions,
		resolver.WithSources(sources...),
		resolver.WithLogger(e.logger),
		resolver.WithTracer(e.tracer),
	)
	e.resolver = resolve
	return diags
}

func (e *Engine) LoadPluginRunner(ctx context.Context) (diags diagnostics.Diag) {
	ctx, span := e.tracer.Start(ctx, "Engine.LoadPluginRunner")
	e.logger.InfoContext(ctx, "Loading plugin runner")
	defer func() {
		if diags.HasErrors() {
			span.RecordError(diags)
			span.SetStatus(codes.Error, diags.Error())
		}
		span.End()
	}()
	binaryMap, diag := e.resolver.Resolve(ctx, e.lockFile)
	if diags.Extend(diag) {
		return diag
	}
	e.runner, diag = runner.Load(ctx, binaryMap, e.builtin, e.logger, e.tracer)
	diag.Extend(diag)
	return diag
}

func (e *Engine) PrintDiagnostics(output io.Writer, diags diagnostics.Diag, colorize bool) {
	diagnostics.PrintDiags(output, diags, e.fileMap, colorize)
}

func (e *Engine) loadGlobalData(ctx context.Context, source, name string) (_ *eval.PluginDataAction, diags diagnostics.Diag) {
	ctx, span := e.tracer.Start(ctx, "Engine.loadGlobalData", trace.WithAttributes(
		attribute.String("datasource", source),
		attribute.String("name", name),
	))
	e.logger.InfoContext(ctx, "Loading global data", "datasource", source, "name", name)
	defer func() {
		if diags.HasErrors() {
			span.RecordError(diags)
			span.SetStatus(codes.Error, diags.Error())
		}
		span.End()
	}()
	if e.blocks == nil {
		return nil, diagnostics.Diag{{
			Severity: hcl.DiagError,
			Summary:  "No files parsed",
			Detail:   "Parse files before selecting",
		}}
	}
	if e.runner == nil {
		return nil, diagnostics.Diag{{
			Severity: hcl.DiagError,
			Summary:  "Plugin runner is not loaded",
			Detail:   "Load plugin runner before evaluating",
		}}
	}
	data, ok := e.blocks.Plugins[definitions.Key{
		PluginKind: definitions.BlockKindData,
		PluginName: source,
		BlockName:  name,
	}]
	if !ok {
		return nil, diagnostics.Diag{{
			Severity: hcl.DiagError,
			Summary:  "Data source not found",
			Detail:   fmt.Sprintf("Data source named '%s' not found", name),
		}}
	}
	parsedData, diag := e.blocks.ParsePlugin(data)
	if diags.Extend(diag) {
		return nil, diags
	}
	loadedData, diag := eval.LoadDataAction(e.runner, parsedData)
	if diags.Extend(diag) {
		return nil, diags
	}
	return loadedData, diags
}

func (e *Engine) loadDocumentData(ctx context.Context, doc, source, name string) (_ *eval.PluginDataAction, diags diagnostics.Diag) {
	ctx, span := e.tracer.Start(ctx, "Engine.loadDocumentData", trace.WithAttributes(
		attribute.String("document", doc),
		attribute.String("datasource", source),
		attribute.String("name", name),
	))
	e.logger.InfoContext(ctx, "Loading document data", "document", doc, "datasource", source, "name", name)
	defer func() {
		if diags.HasErrors() {
			span.RecordError(diags)
			span.SetStatus(codes.Error, diags.Error())
		}
		span.End()
	}()
	if e.blocks == nil {
		return nil, diagnostics.Diag{{
			Severity: hcl.DiagError,
			Summary:  "No files parsed",
			Detail:   "Parse files before selecting",
		}}
	}
	if e.runner == nil {
		return nil, diagnostics.Diag{{
			Severity: hcl.DiagError,
			Summary:  "Plugin runner is not loaded",
			Detail:   "Load plugin runner before evaluating",
		}}
	}
	docBlock, ok := e.blocks.Documents[doc]
	if !ok {
		return nil, diagnostics.Diag{{
			Severity: hcl.DiagError,
			Summary:  "Document not found",
			Detail:   fmt.Sprintf("Definition for document named '%s' not found", doc),
		}}
	}
	docParsed, diag := e.blocks.ParseDocument(docBlock)
	if diags.Extend(diag) {
		return nil, diags
	}
	idx := slices.IndexFunc(docParsed.Data, func(p *definitions.ParsedPlugin) bool {
		return p.PluginName == source && p.BlockName == name
	})
	if idx < 0 {
		return nil, diagnostics.Diag{{
			Severity: hcl.DiagError,
			Summary:  "Data source not found",
			Detail:   fmt.Sprintf("Data source named '%s' not found", name),
		}}
	}
	loadedData, diag := eval.LoadDataAction(e.runner, docParsed.Data[idx])
	if diags.Extend(diag) {
		return nil, diags
	}
	return loadedData, diags
}

var ErrInvalidDataTarget = diagnostics.Diag{{
	Severity: hcl.DiagError,
	Summary:  "Invalid data target",
	Detail:   "Target must be in the format 'document.<doc-name>.data.<plugin-name>.<block-name>' or 'data.<plugin-name>.<block-name>'",
}}

func (e *Engine) FetchData(ctx context.Context, target string) (_ plugin.Data, diags diagnostics.Diag) {
	ctx, span := e.tracer.Start(ctx, "Engine.FetchData", trace.WithAttributes(
		attribute.String("target", target),
	))
	e.logger.InfoContext(ctx, "Fetching data", "target", target)
	defer func() {
		if diags.HasErrors() {
			span.RecordError(diags)
			span.SetStatus(codes.Error, diags.Error())
		}
		span.End()
	}()
	head, base, ok := strings.Cut(target, ".")
	if !ok {
		return nil, ErrInvalidDataTarget
	}
	var loadedData *eval.PluginDataAction
	var diag diagnostics.Diag
	switch head {
	case "document":
		parts := strings.Split(base, ".")
		if len(parts) != 4 {
			return nil, ErrInvalidDataTarget
		}
		if parts[1] != "data" {
			return nil, ErrInvalidDataTarget
		}
		loadedData, diag = e.loadDocumentData(ctx, parts[0], parts[2], parts[3])
	case "data":
		parts := strings.Split(base, ".")
		if len(parts) != 2 {
			return nil, ErrInvalidDataTarget
		}
		loadedData, diag = e.loadGlobalData(ctx, parts[0], parts[1])
	default:
		return nil, ErrInvalidDataTarget
	}
	if diags.Extend(diag) {
		return nil, diags
	}
	return loadedData.FetchData(ctx)
}

func (e *Engine) RenderContent(ctx context.Context, target string) (_ plugin.Content, _ plugin.Data, diags diagnostics.Diag) {
	ctx, span := e.tracer.Start(ctx, "Engine.RenderContent", trace.WithAttributes(
		attribute.String("target", target),
	))
	e.logger.InfoContext(ctx, "Rendering content", "target", target)
	defer func() {
		if diags.HasErrors() {
			span.RecordError(diags)
			span.SetStatus(codes.Error, diags.Error())
		}
		span.End()
	}()
	doc, diag := e.loadDocument(ctx, target)
	if diags.Extend(diag) {
		return nil, nil, diags
	}
	content, data, diag := doc.RenderContent(ctx)
	if diags.Extend(diag) {
		return nil, nil, diags
	}
	return content, data, diags
}

func (e *Engine) Publish(ctx context.Context, target string) (_ plugin.Content, _ plugin.Data, diags diagnostics.Diag) {
	ctx, span := e.tracer.Start(ctx, "Engine.Publish", trace.WithAttributes(
		attribute.String("target", target),
	))
	e.logger.InfoContext(ctx, "Publishing document", "target", target)
	defer func() {
		if diags.HasErrors() {
			span.RecordError(diags)
			span.SetStatus(codes.Error, diags.Error())
		}
		span.End()
	}()
	doc, diag := e.loadDocument(ctx, target)
	if diags.Extend(diag) {
		return nil, nil, diags
	}
	content, data, diag := doc.Publish(ctx)
	if diags.Extend(diag) {
		return nil, nil, diags
	}
	return content, data, diags
}

func (e *Engine) loadDocument(ctx context.Context, name string) (_ *eval.Document, diags diagnostics.Diag) {
	ctx, span := e.tracer.Start(ctx, "Engine.loadDocument", trace.WithAttributes(
		attribute.String("target", name),
	))
	e.logger.InfoContext(ctx, "Loading document", "target", name)
	defer func() {
		if diags.HasErrors() {
			span.RecordError(diags)
			span.SetStatus(codes.Error, diags.Error())
		}
		span.End()
	}()
	if e.runner == nil {
		diags.Append(&hcl.Diagnostic{
			Severity: hcl.DiagError,
			Summary:  "Plugin runner is not loaded",
			Detail:   "Load plugin runner before evaluating",
		})
		return nil, diags
	}
	doc, ok := e.blocks.Documents[name]
	if !ok {
		diags.Append(&hcl.Diagnostic{
			Severity: hcl.DiagError,
			Summary:  "Document not found",
			Detail:   fmt.Sprintf("Definition for document named '%s' not found", name),
		})
		return nil, diags
	}
	parsedDoc, diag := e.blocks.ParseDocument(doc)
	if diags.Extend(diag) {
		return nil, diags
	}
	loadedDoc, diag := eval.LoadDocument(e.runner, parsedDoc)
	if diags.Extend(diag) {
		return nil, diags
	}
	return loadedDoc, diags
}

func (e *Engine) Cleanup() diagnostics.Diag {
	if e.runner != nil {
		return e.runner.Close()
	}
	return nil
}
