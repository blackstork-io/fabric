package cmd

import (
	"log/slog"
	"os"

	"github.com/spf13/cobra"

	"github.com/blackstork-io/fabric/engine"
	"github.com/blackstork-io/fabric/internal/builtin"
	"github.com/blackstork-io/fabric/pkg/diagnostics"
)

var fullLint bool

func init() {
	lintCmd.Flags().BoolVar(&fullLint, "full", false, "Lint plugin bodies (requires plugins to be installed)")
	rootCmd.AddCommand(lintCmd)
}

var lintCmd = &cobra.Command{
	Use:   "lint",
	Short: "Evaluate *.fabric files for syntax mistakes",
	Long:  `Doesn't call plugins, only checks the *.fabric templates for correctness`,
	RunE: func(cmd *cobra.Command, _ []string) (err error) {
		ctx := cmd.Context()
		var diags diagnostics.Diag
		eng := engine.New(
			engine.WithLogger(slog.Default()),
			engine.WithTracer(tracer),
			engine.WithBuiltIn(builtin.Plugin(version, slog.Default(), tracer)),
		)
		defer func() {
			diags.Extend(eng.Cleanup())
			if diags.HasErrors() {
				err = diags
				cmd.SilenceErrors = true
				cmd.SilenceUsage = true
			}
			eng.PrintDiagnostics(os.Stderr, diags, cliArgs.colorize)
		}()
		diag := eng.ParseDir(ctx, os.DirFS(cliArgs.sourceDir))
		if diags.Extend(diag) {
			return diags
		}
		if fullLint {
			if diags.Extend(eng.LoadPluginResolver(ctx, false)) {
				return
			}
			if diags.Extend(eng.LoadPluginRunner(ctx)) {
				return
			}
		}
		diag = eng.Lint(ctx, fullLint)
		diags.Extend(diag)
		return nil
	},
}
