package cmd

import (
	"bufio"
	"fmt"
	"io"
	"os"

	"github.com/spf13/cobra"

	"github.com/blackstork-io/fabric/internal/builtin"
	"github.com/blackstork-io/fabric/parser"
	"github.com/blackstork-io/fabric/pkg/diagnostics"
	"github.com/blackstork-io/fabric/plugin/runner"
)

var outFile string

func render(docName string, w io.Writer) (diags diagnostics.Diag) {
	result := parser.ParseDir(os.DirFS(SourceDir))
	diags = result.Diags
	defer func() { diagnostics.PrintDiags(diags, result.FileMap) }()
	if diags.HasErrors() {
		return
	}
	if len(result.FileMap) == 0 {
		diags.Add(
			"No correct fabric files found",
			fmt.Sprintf("There are no *.fabric files at '%s' or all of them have failed to parse", SourceDir),
		)
	}

	doc, found := result.Blocks.Documents[docName]
	if !found {
		diags.Add(
			"Document not found",
			fmt.Sprintf(
				"Definition for document named '%s' not found in '%s/**.fabric' files",
				docName,
				SourceDir,
			),
		)
	}

	// TODO: read from config
	pluginPath := "./plugins"
	runner, stdDiag := runner.Load(
		runner.WithBuiltIn(
			builtin.Plugin(version),
		),
		runner.WithPluginDir(pluginPath),
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
	diags.Extend(writeResults(w, results))
	return
}

func writeResults(w io.Writer, results []string) (diags diagnostics.Diag) {
	if len(results) == 0 {
		diags.Add("Empty output", "No content was produced")
		return
	}
	var err error
	defer func() {
		diags.AppendErr(err, "Error while outputing result")
	}()
	_, err = w.Write([]byte(results[0]))
	if err != nil {
		return
	}
	for _, result := range results[1:] {
		_, err = w.Write([]byte("\n\n"))
		if err != nil {
			return
		}
		_, err = w.Write([]byte(result))
		if err != nil {
			return
		}
	}
	_, err = w.Write([]byte("\n"))
	return
}

// renderCmd represents the render command
var renderCmd = &cobra.Command{
	Use:   "render TARGET",
	Short: "Render the document",
	Long:  `Render the specified document into Markdown and output it either to stdout or to a file`,
	RunE: func(_ *cobra.Command, args []string) (err error) {
		var out *os.File
		if outFile == "" {
			out = os.Stdout
		} else {
			out, err = os.Create(outFile)
			if err != nil {
				return
			}
			defer out.Close()
		}
		wr := bufio.NewWriter(out)
		defer wr.Flush()
		render(args[0], wr)
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
