package cmd

import (
	"io/fs"
	"os"

	"github.com/hashicorp/hcl/v2"

	"github.com/blackstork-io/fabric/internal/builtin"
	"github.com/blackstork-io/fabric/parser"
	"github.com/blackstork-io/fabric/pkg/diagnostics"
	"github.com/blackstork-io/fabric/plugin/runner"
)

type Evaluator struct {
	PluginsDir string
	Blocks     *parser.DefinedBlocks
	Runner     *runner.Runner
	FileMap    map[string]*hcl.File
}

func NewEvaluator(pluginsDir string) *Evaluator {
	return &Evaluator{
		PluginsDir: pluginsDir,
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
	if e.PluginsDir == "" && e.Blocks.GlobalConfig != nil && e.Blocks.GlobalConfig.PluginRegistry != nil {
		// use pluginsDir from config, unless overridden by cli arg
		e.PluginsDir = e.Blocks.GlobalConfig.PluginRegistry.MirrorDir
	}
	return
}

func (e *Evaluator) LoadRunner() diagnostics.Diag {
	var pluginVersions runner.VersionMap
	if e.Blocks.GlobalConfig != nil {
		pluginVersions = e.Blocks.GlobalConfig.PluginVersions
	}
	var stdDiag hcl.Diagnostics

	e.Runner, stdDiag = runner.Load(
		runner.WithBuiltIn(
			builtin.Plugin(version),
		),
		runner.WithPluginDir(e.PluginsDir),
		runner.WithPluginVersions(pluginVersions),
	)
	return diagnostics.Diag(stdDiag)
}

func (e *Evaluator) PluginCaller() *parser.Caller {
	return parser.NewPluginCaller(e.Runner)
}
