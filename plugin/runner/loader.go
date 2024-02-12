package runner

import (
	"fmt"
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
	resolver   *resolver
	versionMap VersionMap
	builtin    []*plugin.Schema
	pluginMap  map[string]loadedPlugin
	dataMap    map[string]loadedDataSource
	contentMap map[string]loadedContentProvider
}

func makeLoader(mirrorDir string, builtin []*plugin.Schema, pluginMap VersionMap) *loader {
	return &loader{
		resolver:   makeResolver(mirrorDir),
		versionMap: pluginMap,
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
	for _, p := range l.builtin {
		if diags := l.registerPlugin(p, nopCloser); diags.HasErrors() {
			diags = append(diags, l.closeAll()...)
			return diags
		}
	}
	for name, version := range l.versionMap {
		if diags := l.loadBinary(name, version); diags.HasErrors() {
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

func (l *loader) registerDataSource(name string, p *plugin.Schema, ds *plugin.DataSource) hcl.Diagnostics {
	if found, has := l.dataMap[name]; has {
		return hcl.Diagnostics{{
			Severity: hcl.DiagError,
			Summary:  "Duplicate data source",
			Detail:   fmt.Sprintf("Data source %s provided by plugin %s@%s and %s@%s", name, p.Name, p.Version, found.plugin.Name, found.plugin.Version),
		}}
	}
	l.dataMap[name] = loadedDataSource{p, ds}
	return nil
}

func (l *loader) registerContentProvider(name string, p *plugin.Schema, cp *plugin.ContentProvider) hcl.Diagnostics {
	if found, has := l.contentMap[name]; has {
		return hcl.Diagnostics{{
			Severity: hcl.DiagError,
			Summary:  "Duplicate content provider",
			Detail:   fmt.Sprintf("Content provider %s provided by plugin %s@%s and %s@%s", name, p.Name, p.Version, found.plugin.Name, found.plugin.Version),
		}}
	}
	l.contentMap[name] = loadedContentProvider{p, cp}
	return nil
}

func (l *loader) registerPlugin(p *plugin.Schema, closefn func() error) hcl.Diagnostics {
	if diags := p.Validate(); diags.HasErrors() {
		return diags
	}
	if found, has := l.pluginMap[p.Name]; has {
		diags := hcl.Diagnostics{{
			Severity: hcl.DiagError,
			Summary:  "Plugin conflict",
			Detail:   fmt.Sprintf("Plugin %s@%s and %s@%s have the same name", p.Name, p.Version, found.Name, found.Version),
		}}
		err := found.closefn()
		if err != nil {
			diags = append(diags, &hcl.Diagnostic{
				Severity: hcl.DiagError,
				Summary:  "Failed to close plugin",
				Detail:   fmt.Sprintf("Failed to close plugin %s@%s: %v", found.Name, found.Version, err),
			})
		}
		return diags
	}
	plugin := loadedPlugin{
		closefn: closefn,
		Schema:  p,
	}
	l.pluginMap[p.Name] = plugin
	for name, source := range p.DataSources {
		if diags := l.registerDataSource(name, p, source); diags.HasErrors() {
			return diags
		}
	}
	for name, provider := range p.ContentProviders {
		if diags := l.registerContentProvider(name, p, provider); diags.HasErrors() {
			return diags
		}
	}
	return nil
}

func (l *loader) loadBinary(name, version string) hcl.Diagnostics {
	loc, diags := l.resolver.resolve(name, version)
	if diags.HasErrors() {
		return diags
	}
	if info, err := os.Stat(loc); os.IsNotExist(err) {
		return hcl.Diagnostics{{
			Severity: hcl.DiagError,
			Summary:  "Plugin not found",
			Detail:   fmt.Sprintf("Plugin %s@%s not found at: %s", name, version, loc),
		}}
	} else if info.IsDir() {
		return hcl.Diagnostics{{
			Severity: hcl.DiagError,
			Summary:  "Plugin is a directory",
			Detail:   fmt.Sprintf("Plugin %s@%s is a directory at: %s", name, version, loc),
		}}
	}
	p, close, err := pluginapiv1.NewClient(loc)
	if err != nil {
		return hcl.Diagnostics{{
			Severity: hcl.DiagError,
			Summary:  "Failed to create plugin client",
			Detail:   fmt.Sprintf("Failed to create plugin client for %s@%s: %v", name, version, err),
		}}
	}
	return l.registerPlugin(p, close)
}
