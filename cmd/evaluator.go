package cmd

import (
	"context"
	"fmt"
	"io/fs"
	"log/slog"
	"os"
	"path/filepath"

	"github.com/hashicorp/hcl/v2"

	"github.com/blackstork-io/fabric/internal/builtin"
	"github.com/blackstork-io/fabric/parser"
	"github.com/blackstork-io/fabric/parser/definitions"
	"github.com/blackstork-io/fabric/pkg/diagnostics"
	"github.com/blackstork-io/fabric/plugin/resolver"
	"github.com/blackstork-io/fabric/plugin/runner"
)

const (
	defaultLockFile = ".fabric-lock.json"
)

type Evaluator struct {
	Config   *definitions.GlobalConfig
	Blocks   *parser.DefinedBlocks
	Runner   *runner.Runner
	LockFile *resolver.LockFile
	Resolver *resolver.Resolver
	FileMap  map[string]*hcl.File
}

func NewEvaluator() *Evaluator {
	return &Evaluator{
		Config: &definitions.GlobalConfig{
			PluginRegistry: &definitions.PluginRegistry{
				BaseURL:   "http://localhost:8080",
				MirrorDir: "",
			},
			CacheDir: ".fabric",
		},
	}
}

func (e *Evaluator) Cleanup(diags diagnostics.Diag) error {
	if e.Runner != nil {
		diags.ExtendHcl(e.Runner.Close())
	}
	diagnostics.PrintDiags(os.Stderr, diags, e.FileMap, cliArgs.colorize)
	// Errors have been already displayed
	if diags.HasErrors() {
		rootCmd.SilenceErrors = true
		rootCmd.SilenceUsage = true
		return diags
	}
	return nil
}

func (e *Evaluator) ParseFabricFiles(sourceDir fs.FS) (diags diagnostics.Diag) {
	e.Blocks, e.FileMap, diags = parser.ParseDir(sourceDir)
	if diags.HasErrors() {
		return
	}
	if e.Blocks.GlobalConfig != nil {
		e.Config.Merge(e.Blocks.GlobalConfig)
	}
	return
}

func (e *Evaluator) LoadPluginRunner() diagnostics.Diag {
	var diag diagnostics.Diag
	binaryMap, diags := e.Resolver.Resolve(context.Background(), e.LockFile)
	if diag.ExtendHcl(diags) {
		return diag
	}
	e.Runner, diags = runner.Load(binaryMap, builtin.Plugin(version), slog.Default())
	diag.ExtendHcl(diags)
	return diag
}

func (e *Evaluator) LoadPluginResolver(includeRemote bool) diagnostics.Diag {
	pluginDir := filepath.Join(e.Config.CacheDir, "plugins")
	sources := []resolver.Source{
		resolver.LocalSource{
			Path: pluginDir,
		},
	}
	if e.Config.PluginRegistry != nil {
		if e.Config.PluginRegistry.MirrorDir != "" {
			sources = append(sources, resolver.LocalSource{
				Path: e.Config.PluginRegistry.MirrorDir,
			})
		}
		if includeRemote && e.Config.PluginRegistry.BaseURL != "" {
			sources = append(sources, resolver.RemoteSource{
				BaseURL:     e.Config.PluginRegistry.BaseURL,
				DownloadDir: pluginDir,
				UserAgent:   fmt.Sprintf("fabric/%s", version),
			})
		}
	}
	var err error
	e.LockFile, err = resolver.ReadLockFileFrom(defaultLockFile)
	if err != nil {
		return diagnostics.Diag{{
			Severity: hcl.DiagError,
			Summary:  "Failed to read lock file",
			Detail:   err.Error(),
		}}
	}
	var diags hcl.Diagnostics
	e.Resolver, diags = resolver.NewResolver(e.Config.PluginVersions,
		resolver.WithLogger(slog.Default()),
		resolver.WithSources(sources...),
	)
	return diagnostics.Diag(diags)
}

func (e *Evaluator) PluginCaller() *parser.Caller {
	return parser.NewPluginCaller(e.Runner)
}
