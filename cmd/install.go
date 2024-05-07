package cmd

import (
	"os"

	"github.com/spf13/cobra"

	"github.com/blackstork-io/fabric/pkg/diagnostics"
	"github.com/blackstork-io/fabric/plugin/resolver"
)

var installUpgrade bool

var installCmd = &cobra.Command{
	Use:   "install",
	Short: "Install plugins",
	Long:  "Install Fabric plugins",
	RunE: func(cmd *cobra.Command, args []string) (err error) {
		var diags diagnostics.Diag
		eval := NewEvaluator()
		defer func() {
			err = eval.Cleanup(diags)
		}()
		diags = eval.ParseFabricFiles(os.DirFS(cliArgs.sourceDir))
		if diags.HasErrors() {
			return
		}
		if diags.Extend(eval.LoadPluginResolver(true)) {
			return
		}
		lockFile, stdDiags := eval.Resolver.Install(cmd.Context(), eval.LockFile, installUpgrade)
		if diags.Extend(stdDiags) {
			return
		}
		return resolver.SaveLockFileTo(defaultLockFile, lockFile)
	},
}

func init() {
	rootCmd.AddCommand(installCmd)
	installCmd.Flags().BoolVarP(&installUpgrade, "upgrade", "u", false, "Upgrade plugin versions")
}
