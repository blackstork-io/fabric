package cmd

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/spf13/cobra"

	"github.com/blackstork-io/fabric/internal/builtin"
	"github.com/blackstork-io/fabric/parser"
	"github.com/blackstork-io/fabric/parser/definitions"
	"github.com/blackstork-io/fabric/pkg/diagnostics"
	"github.com/blackstork-io/fabric/plugin/runner"
)

var outFile string

func render(dest io.Writer, docName string) {
	result := parser.ParseDir(os.DirFS(cliArgs.sourceDir))
	diags := result.Diags
	defer func() {
		diagnostics.PrintDiags(os.Stderr, diags, result.FileMap, cliArgs.colorize)
	}()
	if diags.HasErrors() {
		return
	}
	if len(result.FileMap) == 0 {
		diags.Add(
			"No correct fabric files found",
			fmt.Sprintf("There are no *.fabric files at '%s' or all of them have failed to parse", cliArgs.sourceDir),
		)
		return
	}

	doc, found := result.Blocks.Documents[docName]
	if !found {
		diags.Add(
			"Document not found",
			fmt.Sprintf(
				"Definition for document named '%s' not found in '%s/**.fabric' files",
				docName,
				cliArgs.sourceDir,
			),
		)
		return
	}

	// TODO: read pluginsDir from config #5
	var pluginsDir string
	if cliArgs.pluginsDir != "" {
		pluginsDir = cliArgs.pluginsDir
	}
	runner, stdDiag := runner.Load(
		runner.WithBuiltIn(
			builtin.Plugin(version),
		),
		runner.WithPluginDir(pluginsDir),
		// TODO: get versions from the fabric configuration file.
		// atm, it's hardcoded to use all plugins with the same version as the CLI.
		runner.WithPluginVersions(runner.VersionMap{
			"blackstork/elasticsearch": version,
			"blackstork/github":        version,
			"blackstork/graphql":       version,
			"blackstork/openai":        version,
			"blackstork/opencti":       version,
			"blackstork/postgresql":    version,
			"blackstork/sqlite":        version,
			"blackstork/terraform":     version,
		}),
	)
	if diags.Extend(diagnostics.Diag(stdDiag)) {
		return
	}
	defer func() { diags.Extend(diagnostics.Diag(runner.Close())) }()

	caller := parser.NewPluginCaller(runner)

	eval := parser.NewEvaluator(caller, result.Blocks)
	results, diag := eval.EvaluateDocument(doc)
	if diags.Extend(diag) {
		return
	}
	diags.Extend(writeResults(dest, results))
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
	RunE: func(_ *cobra.Command, args []string) (err error) {
		target := strings.TrimSpace(args[0])
		const docPrefix = definitions.BlockKindDocument + "."
		switch {
		case strings.HasPrefix(target, docPrefix):
			target = target[len(docPrefix):]
		default:
			return fmt.Errorf("target should have the format '%s<name_of_the_document>'", docPrefix)
		}

		var dest *os.File
		if outFile == "" {
			dest = os.Stdout
		} else {
			dest, err = os.Create(outFile)
			if err != nil {
				return fmt.Errorf("can't create the out-file: %w", err)
			}
			defer dest.Close()
		}
		render(dest, target)
		return nil
	},
	Args: cobra.ExactArgs(1),
}

func init() {
	rootCmd.AddCommand(renderCmd)

	renderCmd.Flags().StringVar(&outFile, "out-file", "", "name of the output file where the rendered document must be saved to. If not set - the Markdown is printed to stdout")

	renderCmd.SetUsageTemplate(UsageTemplate(
		[2]string{"TARGET", "name of the document to be rendered as 'document.<name>'"},
	))
}
