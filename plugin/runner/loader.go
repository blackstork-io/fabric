package runner

import (
	"context"
	"fmt"
	"log/slog"
	"os"

	"github.com/hashicorp/hcl/v2"

	"github.com/blackstork-io/fabric/plugin"
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

type loader struct {
	logger     *slog.Logger
	binaryMap  map[string]string
	builtin    *plugin.Schema
	pluginMap  map[string]loadedPlugin
	dataMap    map[string]loadedDataSource
	contentMap map[string]loadedContentProvider
}

func makeLoader(binaryMap map[string]string, builtin *plugin.Schema, logger *slog.Logger) *loader {
	return &loader{
		logger:     logger,
		binaryMap:  binaryMap,
		builtin:    builtin,
		pluginMap:  make(map[string]loadedPlugin),
		dataMap:    make(map[string]loadedDataSource),
		contentMap: make(map[string]loadedContentProvider),
	}
}

func nopCloser() error {
	return nil
}

func (l *loader) loadAll() hcl.Diagnostics {
	if diags := l.registerPlugin(l.builtin, nopCloser); diags.HasErrors() {
		diags = append(diags, l.closeAll()...)
		return diags
	}
	for name, binaryPath := range l.binaryMap {
		if diags := l.loadBinary(name, binaryPath); diags.HasErrors() {
			diags = append(diags, l.closeAll()...)
			return diags
		}
	}
	return nil
}

func (l *loader) closeAll() hcl.Diagnostics {
	var diags hcl.Diagnostics
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

func (l *loader) registerDataSource(name string, schema *plugin.Schema, ds *plugin.DataSource) hcl.Diagnostics {
	if found, has := l.dataMap[name]; has {
		return hcl.Diagnostics{{
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

func (l *loader) registerContentProvider(name string, schema *plugin.Schema, cp *plugin.ContentProvider) hcl.Diagnostics {
	if found, has := l.contentMap[name]; has {
		return hcl.Diagnostics{{
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
			ContentFunc: func(ctx context.Context, params *plugin.ProvideContentParams) (*plugin.ContentResult, hcl.Diagnostics) {
				return schema.ProvideContent(ctx, name, params)
			},
			InvocationOrder: cp.InvocationOrder,
		},
	}
	return nil
}

func (l *loader) registerPlugin(schema *plugin.Schema, closefn func() error) hcl.Diagnostics {
	if diags := schema.Validate(); diags.HasErrors() {
		return diags
	}
	if found, has := l.pluginMap[schema.Name]; has {
		diags := hcl.Diagnostics{{
			Severity: hcl.DiagError,
			Summary:  fmt.Sprintf("Plugin %s conflict", schema.Name),
			Detail:   fmt.Sprintf("%s@%s and %s@%s have the same schema name", schema.Name, schema.Version, found.Name, found.Version),
		}}
		err := found.closefn()
		if err != nil {
			diags = append(diags, &hcl.Diagnostic{
				Severity: hcl.DiagError,
				Summary:  fmt.Sprintf("Failed to close plugin %s@%s", found.Name, found.Version),
				Detail:   err.Error(),
			})
		}
		return diags
	}
	plugin := loadedPlugin{
		closefn: closefn,
		Schema:  schema,
	}
	l.pluginMap[schema.Name] = plugin
	for name, source := range schema.DataSources {
		if diags := l.registerDataSource(name, schema, source); diags.HasErrors() {
			return diags
		}
	}
	for name, provider := range schema.ContentProviders {
		if diags := l.registerContentProvider(name, schema, provider); diags.HasErrors() {
			return diags
		}
	}
	return nil
}

func (l *loader) loadBinary(name, binaryPath string) hcl.Diagnostics {
	if info, err := os.Stat(binaryPath); os.IsNotExist(err) {
		return hcl.Diagnostics{{
			Severity: hcl.DiagError,
			Summary:  fmt.Sprintf("Plugin %s binary not found", name),
			Detail:   fmt.Sprintf("Executable not found at: %s", binaryPath),
		}}
	} else if info.IsDir() {
		return hcl.Diagnostics{{
			Severity: hcl.DiagError,
			Summary:  fmt.Sprintf("Plugin %s binary path is a directory", name),
			Detail:   fmt.Sprintf("Path %s is a directory", binaryPath),
		}}
	}
	l.logger.Info("Loading plugin", "name", name, "path", binaryPath)
	p, close, err := pluginapiv1.NewClient(name, binaryPath, l.logger)
	if err != nil {
		return hcl.Diagnostics{{
			Severity: hcl.DiagError,
			Summary:  fmt.Sprintf("Failed to load plugin %s", name),
			Detail:   err.Error(),
		}}
	}
	return l.registerPlugin(p, close)
}
