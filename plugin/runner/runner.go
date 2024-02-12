package runner

import (
	"fmt"
	"strings"
	"unicode"

	"github.com/Masterminds/semver/v3"
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

func validatePluginName(name string) hcl.Diagnostics {
	parts := strings.Split(name, "/")
	if len(parts) != 2 {
		return hcl.Diagnostics{{
			Severity: hcl.DiagError,
			Summary:  "Invalid plugin name",
			Detail:   fmt.Sprintf("plugin name '%s' is not in the form '<namespace>/<name>'", name),
		}}
	}
	for _, r := range parts[0] {
		if !unicode.IsLetter(r) && !unicode.IsDigit(r) && r != '-' {
			return hcl.Diagnostics{{
				Severity: hcl.DiagError,
				Summary:  "Invalid plugin name",
				Detail:   fmt.Sprintf("plugin name '%s' contains invalid character: '%c'", name, r),
			}}
		}
	}
	for _, r := range parts[1] {
		if !unicode.IsLetter(r) && !unicode.IsDigit(r) && r != '-' {
			return hcl.Diagnostics{{
				Severity: hcl.DiagError,
				Summary:  "Invalid plugin name",
				Detail:   fmt.Sprintf("plugin name '%s' contains invalid character: '%c'", name, r),
			}}
		}
	}
	return nil
}

func validatePluginVersionMap(versionMap VersionMap) (diags hcl.Diagnostics) {
	for name, version := range versionMap {
		diags = validatePluginName(name).Extend(diags)
		if version == "" {
			diags = diags.Append(&hcl.Diagnostic{
				Severity: hcl.DiagError,
				Summary:  "Missing plugin version",
				Detail:   fmt.Sprintf("Missing plugin version for '%s'", name),
			})
			continue
		}
		_, err := semver.NewConstraint(version)
		if err != nil {
			diags = diags.Append(&hcl.Diagnostic{
				Severity: hcl.DiagError,
				Summary:  "Invalid plugin version",
				Detail:   fmt.Sprintf("Invalid version constraint for '%s': %s", name, err),
			})
		}
	}
	return diags
}

func Load(o ...Option) (*Runner, hcl.Diagnostics) {
	opts := defaultOptions
	for _, opt := range o {
		opt(&opts)
	}
	if diags := validatePluginVersionMap(opts.versionMap); diags.HasErrors() {
		return nil, diags
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
