package resolver

import (
	"context"
	"fmt"
	"maps"
	"slices"

	"github.com/hashicorp/hcl/v2"
)

// Resolver resolves and installs plugins.
type Resolver struct {
	constraints ConstraintMap
	options
}

// NewResolver creates a new plugin resolver.
func NewResolver(constraints map[string]string, opts ...Option) (*Resolver, hcl.Diagnostics) {
	parsedVersions, err := ParseConstraintMap(constraints)
	if err != nil {
		return nil, hcl.Diagnostics{{
			Severity: hcl.DiagError,
			Summary:  "Failed to parse plugin versions",
			Detail:   err.Error(),
		}}
	}
	res := &Resolver{
		constraints: parsedVersions,
		options:     defaultOptions,
	}
	for _, opt := range opts {
		opt(&res.options)
	}
	return res, nil
}

// Install all plugins based the version constraints and return updated a lock file.
func (r *Resolver) Install(ctx context.Context, lockFile *LockFile, upgrade bool) (*LockFile, hcl.Diagnostics) {
	check := lockFile.Check(r.constraints)
	locks := []PluginLock{}
	// lookupMap is a map of plugins that are we look up based on the constraints
	lookupMap := make(ConstraintMap)
	if upgrade {
		// if upgrade is enabled we install all plugins based on the constraints
		maps.Copy(lookupMap, r.constraints)
	} else {
		// if upgrade is disabled we only install the missing and mismatched plugins based on the constraints
		maps.Copy(lookupMap, check.Missing)
		maps.Copy(lookupMap, check.Mismatch)
	}
	chain := makeSourceChain(r.sources...)
	// resolve the plugins by the latest version that matches the constraints
	for name, constraint := range lookupMap {
		r.logger.Info("Searching plugin", "name", name.String(), "constraints", constraint.String())
		list, err := chain.Lookup(ctx, name)
		if err != nil {
			return nil, hcl.Diagnostics{{
				Severity: hcl.DiagError,
				Summary:  fmt.Sprintf("Failed to lookup plugin '%s'", name),
				Detail:   err.Error(),
			}}
		}
		if len(list) == 0 {
			return nil, hcl.Diagnostics{{
				Severity: hcl.DiagError,
				Summary:  fmt.Sprintf("Plugin '%s' not found", name),
				Detail:   "Could not find version for the current platform",
			}}
		}
		// filter out the versions that do not match the constraint
		matches := slices.DeleteFunc(list, func(v Version) bool {
			return !constraint.Check(v.Version)
		})
		if len(matches) == 0 {
			return nil, hcl.Diagnostics{{
				Severity: hcl.DiagError,
				Summary:  fmt.Sprintf("Plugin '%s' not found", name),
				Detail:   fmt.Sprintf("No version of '%s' matches the constraint '%s'", name, constraint),
			}}
		}
		max := slices.MaxFunc(matches, func(a, b Version) int {
			return a.Compare(b)
		})
		r.logger.Info("Installing plugin", "name", name.String(), "version", max.String())
		var checksums []Checksum
		// check if the plugin with the same version is already in the lock file
		lockIdx := slices.IndexFunc(lockFile.Plugins, func(lock PluginLock) bool {
			return lock.Name == name && lock.Version.Equal(max.Version)
		})
		if lockIdx > -1 {
			// use the checksums from the lock file
			checksums = lockFile.Plugins[lockIdx].Checksums
		}
		res, err := chain.Resolve(ctx, name, max, checksums)
		if err != nil {
			return nil, hcl.Diagnostics{{
				Severity: hcl.DiagError,
				Summary:  fmt.Sprintf("Failed to resolve plugin '%s@%s'", name, max),
				Detail:   err.Error(),
			}}
		}
		// sort the checksums
		slices.SortFunc(res.Checksums, func(a, b Checksum) int {
			return a.Compare(b)
		})
		locks = append(locks, PluginLock{
			Name:      name,
			Version:   max,
			Checksums: res.Checksums,
		})
		// check if context is cancelled
		if ctx.Err() != nil {
			return nil, hcl.Diagnostics{{
				Severity: hcl.DiagError,
				Summary:  "Cancelled by context",
				Detail:   ctx.Err().Error(),
			}}
		}
	}
	// resolve the rest of plugins based on the strict locked versions
	for _, lock := range lockFile.Plugins {
		// skip plugins that are already resolved
		if _, ok := lookupMap[lock.Name]; ok {
			continue
		}
		// skip plugins that are removed from the version constraints
		if _, ok := check.Removed[lock.Name]; ok {
			continue
		}
		r.logger.Info("Installing plugin", "name", lock.Name.String(), "version", lock.Version.String())
		_, err := chain.Resolve(ctx, lock.Name, lock.Version, lock.Checksums)
		if err != nil {
			return nil, hcl.Diagnostics{{
				Severity: hcl.DiagError,
				Summary:  fmt.Sprintf("Failed to resolve plugin '%s@%s'", lock.Name, lock.Version),
				Detail:   err.Error(),
			}}
		}
		locks = append(locks, lock)
		// check if context is cancelled
		if ctx.Err() != nil {
			return nil, hcl.Diagnostics{{
				Severity: hcl.DiagError,
				Summary:  "Cancelled by context",
				Detail:   ctx.Err().Error(),
			}}
		}
	}
	// sort plugin locks by name
	slices.SortFunc(locks, func(a, b PluginLock) int {
		return a.Name.Compare(b.Name)
	})
	return &LockFile{
		Plugins: locks,
	}, nil
}

// Resolve all plugins based on the constraints and returns a map of plugin names to binary paths.
// If the lock file is not satisfied, an error is returned.
func (r *Resolver) Resolve(ctx context.Context, lockFile *LockFile) (map[string]string, hcl.Diagnostics) {
	var diags hcl.Diagnostics
	// check if the lock file is satisfied by version constraints
	check := lockFile.Check(r.constraints)
	for name, lock := range check.Removed {
		// warn about plugins that are removed from the version constraints
		diags = diags.Extend(hcl.Diagnostics{{
			Severity: hcl.DiagWarning,
			Summary:  fmt.Sprintf("Plugin '%s' is not used", name),
			Detail:   fmt.Sprintf("Version '%s' is no longer used. Run install to update lock file", lock),
		}})
	}
	if check.IsInstallRequired() {
		// error out about missing & mismatched plugins
		for name := range check.Missing {
			diags = diags.Extend(hcl.Diagnostics{{
				Severity: hcl.DiagError,
				Summary:  fmt.Sprintf("Plugin '%s' is not locked", name),
				Detail:   "Run install to resolve missing plugins.",
			}})
		}
		for name, constraint := range check.Mismatch {
			pluginIdx := slices.IndexFunc(lockFile.Plugins, func(lock PluginLock) bool {
				return lock.Name == name
			})
			if pluginIdx == -1 {
				continue
			}
			detailFormat := "Version locked at '%s' does not match the new constraint '%s'\nRun install to update lock file."
			diags = diags.Extend(hcl.Diagnostics{{
				Severity: hcl.DiagError,
				Summary:  fmt.Sprintf("Plugin '%s' version mismatch", name),
				Detail:   fmt.Sprintf(detailFormat, lockFile.Plugins[pluginIdx].Version, constraint),
			}})
		}
		return nil, diags
	}
	// chain the sources together
	chain := makeSourceChain(r.sources...)
	// resolve the plugins
	binaryMap := make(map[string]string)
	for _, lock := range lockFile.Plugins {
		// skip plugins that are removed from the version constraints
		if _, ok := check.Removed[lock.Name]; ok {
			continue
		}
		plugin, err := chain.Resolve(ctx, lock.Name, lock.Version, lock.Checksums)
		if err != nil {
			return nil, diags.Extend(hcl.Diagnostics{{
				Severity: hcl.DiagError,
				Summary:  fmt.Sprintf("Failed to resolve plugin '%s@%s'", lock.Name, lock.Version),
				Detail:   err.Error(),
			}})
		}
		binaryMap[lock.Name.String()] = plugin.BinaryPath
		// check if context is cancelled
		if ctx.Err() != nil {
			return nil, diags.Extend(hcl.Diagnostics{{
				Severity: hcl.DiagError,
				Summary:  "Context cancelled",
				Detail:   ctx.Err().Error(),
			}})
		}
	}
	return binaryMap, diags
}
