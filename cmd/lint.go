package cmd

import (
	"context"
	"os"

	"github.com/spf13/cobra"

	"github.com/blackstork-io/fabric/parser"
	"github.com/blackstork-io/fabric/parser/lint"
	"github.com/blackstork-io/fabric/pkg/diagnostics"
)

func Lint(ctx context.Context, blocks *parser.DefinedBlocks, pluginCaller *parser.Caller) (diags diagnostics.Diag) {
	ctx = lint.MakeLintContext(ctx)
	for _, doc := range blocks.Documents {
		pd, diag := blocks.ParseDocument(doc)
		if diags.Extend(diag) {
			continue
		}

		_, diag = pd.Render(ctx, pluginCaller)
		diags.Extend(diag)
	}
	return
}

var lintCmd = &cobra.Command{
	Use:   "lint",
	Short: "Evaluate *.fabric files and report mistakes",
	Long:  `Doesn't call plugins, only checks the *.fabric templates for correctness`,
	RunE: func(cmd *cobra.Command, _ []string) (err error) {
		var diags diagnostics.Diag
		eval := NewEvaluator()
		defer func() {
			err = eval.Cleanup(diags)
		}()
		diags = eval.ParseFabricFiles(os.DirFS(cliArgs.sourceDir))
		if diags.HasErrors() {
			return
		}
		if diags.Extend(eval.LoadPluginResolver(false)) {
			return
		}
		if diags.Extend(eval.LoadPluginRunner(cmd.Context())) {
			return
		}
		diags.Extend(Lint(cmd.Context(), eval.Blocks, eval.PluginCaller()))
		return
	},
	Args: cobra.ExactArgs(1),
}

func init() {
	rootCmd.AddCommand(lintCmd)
}
