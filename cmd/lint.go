package cmd

import (
	"context"
	"io/fs"
	"os"

	"github.com/spf13/cobra"

	"github.com/blackstork-io/fabric/parser/evaluation"
	"github.com/blackstork-io/fabric/pkg/diagnostics"
	"github.com/blackstork-io/fabric/pkg/fabctx"
	"github.com/blackstork-io/fabric/plugin"
)

type noopPluginCaller struct{}

// CallContent implements evaluation.PluginCaller.
func (n *noopPluginCaller) CallContent(ctx context.Context, name string, config evaluation.Configuration, invocation evaluation.Invocation, context plugin.MapData) (result *plugin.Content, diag diagnostics.Diag) {
	return nil, nil
}

func (n *noopPluginCaller) ContentInvocationOrder(ctx context.Context, name string) (order plugin.InvocationOrder, diag diagnostics.Diag) {
	return plugin.InvocationOrderUnspecified, nil
}

// CallData implements evaluation.PluginCaller.
func (n *noopPluginCaller) CallData(ctx context.Context, name string, config evaluation.Configuration, invocation evaluation.Invocation) (result plugin.Data, diag diagnostics.Diag) {
	return plugin.MapData{}, nil
}

var _ evaluation.PluginCaller = (*noopPluginCaller)(nil)

func Lint(ctx context.Context, eval *Evaluator, sourceDir fs.FS, fullLint bool) (diags diagnostics.Diag) {
	diags = eval.ParseFabricFiles(sourceDir)
	if diags.HasErrors() {
		return
	}

	var caller evaluation.PluginCaller
	if fullLint {
		if diags.Extend(eval.LoadPluginResolver(false)) {
			return
		}
		if diags.Extend(eval.LoadPluginRunner(ctx)) {
			return
		}
		caller = eval.PluginCaller()
	} else {
		caller = &noopPluginCaller{}
	}

	for _, doc := range eval.Blocks.Documents {
		pd, diag := eval.Blocks.ParseDocument(doc)
		if diags.Extend(diag) {
			continue
		}

		_, diag = pd.Render(ctx, caller)
		diags.Extend(diag)
	}
	return
}

var lintCmd = &cobra.Command{
	Use:   "lint",
	Short: "Evaluate *.fabric files for syntax mistakes",
	Long:  `Doesn't call plugins, only checks the *.fabric templates for correctness`,
	RunE: func(cmd *cobra.Command, _ []string) (err error) {
		ctx := fabctx.Get(cmd.Context())
		ctx = fabctx.WithLinting(ctx)

		var diags diagnostics.Diag

		eval := NewEvaluator()
		defer func() {
			err = eval.Cleanup(diags)
		}()

		diags = Lint(ctx, eval, os.DirFS(cliArgs.sourceDir), fullLint)

		return
	},
}
var fullLint bool

func init() {
	lintCmd.Flags().BoolVar(&fullLint, "full", false, "Lint plugin bodies (requires plugins to be installed)")

	rootCmd.AddCommand(lintCmd)
}
