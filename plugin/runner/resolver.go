package runner

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/Masterminds/semver/v3"
	"github.com/hashicorp/hcl/v2"
)

type resolver struct {
	mirrorDir string
}

func makeResolver(mirrorDir string) *resolver {
	return &resolver{
		mirrorDir: mirrorDir,
	}
}

func (r *resolver) resolve(name, version string) (loc string, diags hcl.Diagnostics) {
	nameSpace, pluginName, err := r.parseName(name)
	if err != nil {
		return "", hcl.Diagnostics{{
			Severity: hcl.DiagError,
			Summary:  "Invalid plugin name",
			Detail:   fmt.Sprintf("Invalid plugin name '%s': %s", name, err),
		}}
	}
	constraint, err := semver.NewConstraint(version)
	if err != nil {
		return "", hcl.Diagnostics{{
			Severity: hcl.DiagError,
			Summary:  "Failed to resolve plugin",
			Detail:   fmt.Sprintf("Invalid version constraint: %s", err),
		}}
	}
	entry, err := os.ReadDir(filepath.Join(r.mirrorDir, nameSpace))
	if err != nil {
		return "", hcl.Diagnostics{{
			Severity: hcl.DiagError,
			Summary:  "Failed to resolve plugin",
			Detail:   fmt.Sprintf("Failed to read directory for namespace '%s': %s", nameSpace, err),
		}}
	}
	matched := make(map[string]*semver.Version)
	for _, e := range entry {
		if e.IsDir() {
			continue
		}
		parts := strings.SplitN(e.Name(), "@", 2)
		if len(parts) != 2 || parts[0] != pluginName {
			continue
		}
		v, err := semver.NewVersion(parts[1])
		if err != nil {
			continue
		}
		if !constraint.Check(v) {
			continue
		}
		matched[parts[1]] = v
	}
	if len(matched) == 0 {
		return "", hcl.Diagnostics{{
			Severity: hcl.DiagError,
			Summary:  "Failed to resolve plugin binary",
			Detail:   fmt.Sprintf("No plugin matches version constraint for %s@%s", name, version),
		}}
	}
	// find latest version that matches version constraint
	var latestVerStr string
	var latestVer *semver.Version
	for str, ver := range matched {
		if latestVer == nil {
			latestVerStr = str
			latestVer = ver
			continue
		}
		if ver.Compare(latestVer) > 0 {
			latestVerStr = str
			latestVer = ver
		}
	}
	return filepath.Join(r.mirrorDir, nameSpace, fmt.Sprintf("%s@%s", pluginName, latestVerStr)), nil
}

func (r *resolver) parseName(name string) (string, string, error) {
	parts := strings.SplitN(name, "/", 2)
	if len(parts) != 2 {
		return "", "", fmt.Errorf("plugin name '%s' is not in the form '<namespace>/<name>'", name)
	}
	return parts[0], parts[1], nil
}
