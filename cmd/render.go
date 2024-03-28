package cmd

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/spf13/cobra"

	"github.com/blackstork-io/fabric/parser"
	"github.com/blackstork-io/fabric/parser/definitions"
	"github.com/blackstork-io/fabric/pkg/diagnostics"
)

func Render(ctx context.Context, blocks *parser.DefinedBlocks, pluginCaller *parser.Caller, docName string) (results []string, diags diagnostics.Diag) {
	doc, found := blocks.Documents.Map[docName]
	if !found {
		diags.Add(
			"Document not found",
			fmt.Sprintf(
				"Definition for document named '%s' not found",
				docName,
			),
		)
		return
	}

	pd, diag := blocks.ParseDocument(doc)
	if diags.Extend(diag) {
		return
	}

	results, diag = pd.Render(ctx, pluginCaller)
	if diags.Extend(diag) {
		return
	}
	return
}

func writeResults(dest io.Writer, results []string) (diags diagnostics.Diag) {
	if len(results) == 0 {
		diags.Add("Empty output", "No content was produced")
		return
	}
	w := bufio.NewWriter(dest)

	// bufio.Writer preserves the first encountered error,
	// so we're only cheking it once at flush
	_, _ = w.WriteString(results[0])
	for _, result := range results[1:] {
		_, _ = w.WriteString("\n\n")
		_, _ = w.WriteString(result)
	}
	_ = w.WriteByte('\n')
	err := w.Flush()
	diags.AppendErr(err, "Error while outputing result")
	return
}

// renderCmd represents the render command
var renderCmd = &cobra.Command{
	Use:   "render TARGET",
	Short: "Render the document",
	Long:  `Render the specified document into Markdown and output it either to stdout or to a file`,
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) (err error) {
		target := strings.TrimSpace(args[0])
		const docPrefix = definitions.BlockKindDocument + "."
		switch {
		case strings.HasPrefix(target, docPrefix):
			target = target[len(docPrefix):]
		default:
			return fmt.Errorf("target should have the format '%s<name_of_the_document>'", docPrefix)
		}

		dest := os.Stdout
		if outFile != "" {
			dest, err = os.Create(outFile)
			if err != nil {
				return fmt.Errorf("can't create the out-file: %w", err)
			}
			defer dest.Close()
		}

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
		res, diag := Render(cmd.Context(), eval.Blocks, eval.PluginCaller(), target)
		if diags.Extend(diag) {
			return
		}
		diags.Extend(
			writeResults(dest, res),
		)
		return
	},
}

var outFile string

func init() {
	rootCmd.AddCommand(renderCmd)

	renderCmd.Flags().StringVar(&outFile, "out-file", "", "name of the output file where the rendered document must be saved to. If not set - the Markdown is printed to stdout")

	renderCmd.SetUsageTemplate(UsageTemplate(
		[2]string{"TARGET", "name of the document to be rendered as 'document.<name>'"},
	))
}
