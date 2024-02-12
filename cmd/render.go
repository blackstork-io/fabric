package cmd

import (
	"bufio"
	"fmt"
	"io"
	"io/fs"
	"os"
	"strings"

	"github.com/hashicorp/hcl/v2"
	"github.com/spf13/cobra"

	"github.com/blackstork-io/fabric/internal/builtin"
	"github.com/blackstork-io/fabric/parser"
	"github.com/blackstork-io/fabric/parser/definitions"
	"github.com/blackstork-io/fabric/pkg/diagnostics"
	"github.com/blackstork-io/fabric/plugin/runner"
)

func Render(pluginsDir string, sourceDir fs.FS, docName string) (results []string, fileMap map[string]*hcl.File, diags diagnostics.Diag) {
	blocks, fileMap, diags := parser.ParseDir(sourceDir)
	if diags.HasErrors() {
		return
	}
	if len(fileMap) == 0 {
		diags.Add(
			"No correct fabric files found",
			fmt.Sprintf("There are no *.fabric files at '%s' or all of them have failed to parse", cliArgs.sourceDir),
		)
		return
	}
	doc, found := blocks.Documents[docName]
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

	if pluginsDir == "" && blocks.GlobalConfig != nil && blocks.GlobalConfig.PluginRegistry != nil {
		// use pluginsDir from config, unless overriden by cli arg
		pluginsDir = blocks.GlobalConfig.PluginRegistry.MirrorDir
	}

	var pluginVersions runner.VersionMap
	if blocks.GlobalConfig != nil {
		pluginVersions = blocks.GlobalConfig.PluginVersions
	}

	runner, stdDiag := runner.Load(
		runner.WithBuiltIn(
			builtin.Plugin(version),
		),
		runner.WithPluginDir(pluginsDir),
		runner.WithPluginVersions(runner.VersionMap(pluginVersions)),
	)
	if diags.ExtendHcl(stdDiag) {
		return
	}
	defer func() { diags.ExtendHcl(runner.Close()) }()

	eval := parser.NewEvaluator(
		parser.NewPluginCaller(runner),
		blocks,
	)
	results, diag := eval.EvaluateDocument(doc)
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
		res, fileMap, diags := Render(cliArgs.pluginsDir, os.DirFS(cliArgs.sourceDir), target)
		if !diags.HasErrors() {
			diags.Extend(writeResults(dest, res))
		}
		diagnostics.PrintDiags(os.Stderr, diags, fileMap, cliArgs.colorize)
		if diags.HasErrors() {
			// Errors have been already displayed
			rootCmd.SilenceErrors = true
			rootCmd.SilenceUsage = true
			return diags
		}
		return nil
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
