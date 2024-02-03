package runner

import (
	"fmt"

	"github.com/hashicorp/hcl/v2"

	"github.com/blackstork-io/fabric/plugin"
)

type VersionMap map[string]string

type Runner struct {
	cacheDir   string
	pluginMap  map[string]loadedPlugin
	dataMap    map[string]loadedDataSource
	contentMap map[string]loadedContentProvider
}

func Load(o ...Option) (*Runner, hcl.Diagnostics) {
	opts := defaultOptions
	for _, opt := range o {
		opt(&opts)
	}
	loader := makeLoader(opts.pluginDir, opts.builtin, opts.versionMap)
	if diags := loader.loadAll(); diags.HasErrors() {
		return nil, diags
	}
	return &Runner{
		cacheDir:   opts.pluginDir,
		pluginMap:  loader.pluginMap,
		dataMap:    loader.dataMap,
		contentMap: loader.contentMap,
	}, nil
}

func (m *Runner) DataSource(name string) (*plugin.DataSource, hcl.Diagnostics) {
	source, has := m.dataMap[name]
	if !has {
		return nil, hcl.Diagnostics{{
			Severity: hcl.DiagError,
			Summary:  "Data source not found",
			Detail:   fmt.Sprintf("data source '%s' not found in any plugin", name),
		}}
	}
	return source.DataSource, nil
}

func (m *Runner) ContentProvider(name string) (*plugin.ContentProvider, hcl.Diagnostics) {
	provider, has := m.contentMap[name]
	if !has {
		return nil, hcl.Diagnostics{{
			Severity: hcl.DiagError,
			Summary:  "Content provider not found",
			Detail:   fmt.Sprintf("content provider '%s' not found in any plugin", name),
		}}
	}
	return provider.ContentProvider, nil
}

func (m *Runner) Close() hcl.Diagnostics {
	var diags hcl.Diagnostics
	for _, p := range m.pluginMap {
		if err := p.closefn(); err != nil {
			diags = append(diags, &hcl.Diagnostic{
				Severity: hcl.DiagWarning,
				Summary:  "Failed to close plugin",
				Detail:   fmt.Sprintf("failed to close plugin %s@%s: %v", p.Name, p.Version, err),
			})
		}
	}
	return diags
}
