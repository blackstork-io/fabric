package runner

import (
	"fmt"
	"log/slog"

	"github.com/hashicorp/hcl/v2"

	"github.com/blackstork-io/fabric/plugin"
)

type Runner struct {
	pluginMap    map[string]loadedPlugin
	dataMap      map[string]loadedDataSource
	contentMap   map[string]loadedContentProvider
	publisherMap map[string]loadedPublisher
}

func Load(binaryMap map[string]string, builtin *plugin.Schema, logger *slog.Logger) (*Runner, hcl.Diagnostics) {
	loader := makeLoader(binaryMap, builtin, logger)
	if diags := loader.loadAll(); diags.HasErrors() {
		return nil, diags
	}
	return &Runner{
		pluginMap:    loader.pluginMap,
		dataMap:      loader.dataMap,
		contentMap:   loader.contentMap,
		publisherMap: loader.publisherMap,
	}, nil
}

func (m *Runner) DataSource(name string) (*plugin.DataSource, hcl.Diagnostics) {
	source, has := m.dataMap[name]
	if !has {
		return nil, hcl.Diagnostics{{
			Severity: hcl.DiagError,
			Summary:  fmt.Sprintf("Missing data source '%s'", name),
			Detail:   fmt.Sprintf("'%s' not found in any plugin", name),
		}}
	}
	return source.DataSource, nil
}

func (m *Runner) ContentProvider(name string) (*plugin.ContentProvider, hcl.Diagnostics) {
	provider, has := m.contentMap[name]
	if !has {
		return nil, hcl.Diagnostics{{
			Severity: hcl.DiagError,
			Summary:  fmt.Sprintf("Missing content provider '%s'", name),
			Detail:   fmt.Sprintf("'%s' not found in any plugin", name),
		}}
	}
	return provider.ContentProvider, nil
}

func (m *Runner) Publisher(name string) (*plugin.Publisher, hcl.Diagnostics) {
	publisher, has := m.publisherMap[name]
	if !has {
		return nil, hcl.Diagnostics{{
			Severity: hcl.DiagError,
			Summary:  fmt.Sprintf("Missing publisher '%s'", name),
			Detail:   fmt.Sprintf("'%s' not found in any plugin", name),
		}}
	}
	return publisher.Publisher, nil
}

func (m *Runner) Close() hcl.Diagnostics {
	var diags hcl.Diagnostics
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
