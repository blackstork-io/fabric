package cmd

import (
	"log/slog"
	"os"

	"github.com/spf13/cobra"

	"github.com/blackstork-io/fabric/engine"
	"github.com/blackstork-io/fabric/internal/builtin"
	"github.com/blackstork-io/fabric/pkg/diagnostics"
)

var installUpgrade bool

func init() {
	rootCmd.AddCommand(installCmd)
	installCmd.Flags().BoolVarP(&installUpgrade, "upgrade", "u", false, "Upgrade plugin versions")
}

var installCmd = &cobra.Command{
	Use:   "install",
	Short: "Install plugins",
	Long:  "Install Fabric plugins",
	RunE: func(cmd *cobra.Command, args []string) (err error) {
		ctx := cmd.Context()
		var diags diagnostics.Diag
		eng := engine.New(
			engine.WithLogger(slog.Default()),
			engine.WithTracer(tracer),
			engine.WithBuiltIn(builtin.Plugin(version, slog.Default(), tracer)),
		)
		defer func() {
			diag := eng.Cleanup()
			if diags.Extend(diag) {
				err = diags
			}
			eng.PrintDiagnostics(os.Stderr, diags, cliArgs.colorize)
		}()
		diag := eng.ParseDir(ctx, os.DirFS(cliArgs.sourceDir))
		if diags.Extend(diag) {
			return diags
		}
		diag = eng.LoadPluginResolver(ctx, true)
		if diags.Extend(diag) {
			return diags
		}
		diag = eng.Install(ctx, installUpgrade)
		if diags.Extend(diag) {
			return diags
		}
		return nil
	},
}
