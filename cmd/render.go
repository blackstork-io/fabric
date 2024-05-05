package cmd

import (
	"context"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/spf13/cobra"

	"github.com/blackstork-io/fabric/parser"
	"github.com/blackstork-io/fabric/parser/definitions"
	"github.com/blackstork-io/fabric/pkg/diagnostics"
	"github.com/blackstork-io/fabric/plugin"
	"github.com/blackstork-io/fabric/printer"
	"github.com/blackstork-io/fabric/printer/htmlprint"
	"github.com/blackstork-io/fabric/printer/mdprint"
)

func Render(ctx context.Context, blocks *parser.DefinedBlocks, pluginCaller *parser.Caller, docName string, w io.Writer) (diags diagnostics.Diag) {
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

	pd, diag := blocks.ParseDocument(doc)
	if diags.Extend(diag) {
		return
	}
	content, _, diag := pd.Render(ctx, pluginCaller)
	if diags.Extend(diag) {
		return
	}
	var print printer.Printer
	switch format {
	case "md":
		print = mdprint.New()
	case "html":
		print = htmlprint.New()
	default:
		diags.Add("Unsupported format", fmt.Sprintf("Format '%s' is not supported for stdout", format))
		return
	}
	err := print.Print(w, content)
	if err != nil {
		diags.Add("Error while rendering", err.Error())
	}
	return
}

func Publish(ctx context.Context, blocks *parser.DefinedBlocks, pluginCaller *parser.Caller, docName string) (diags diagnostics.Diag) {
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

	pd, diag := blocks.ParseDocument(doc)
	if diags.Extend(diag) {
		return
	}
	var defFormat plugin.OutputFormat
	switch format {
	case "md":
		defFormat = plugin.OutputFormatMD
	case "html":
		defFormat = plugin.OutputFormatHTML
	case "pdf":
		defFormat = plugin.OutputFormatPDF
	default:
		diags.Add("Unsupported format", fmt.Sprintf("Format '%s' is not supported for publishing", format))
		return
	}
	diags.Extend(pd.Publish(ctx, pluginCaller, defFormat))
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
		if publish {
			diags.Extend(Publish(cmd.Context(), eval.Blocks, eval.PluginCaller(), target))
		} else {
			diags.Extend(Render(cmd.Context(), eval.Blocks, eval.PluginCaller(), target, os.Stdout))
		}
		return
	},
}

var (
	publish bool
	format  string
)

func init() {
	rootCmd.AddCommand(renderCmd)

	renderCmd.Flags().BoolVar(&publish, "publish", false, "publish the rendered document")
	renderCmd.Flags().StringVar(&format, "format", "md", "default output format of the document (md, html or pdf)")
	renderCmd.SetUsageTemplate(UsageTemplate(
		[2]string{"TARGET", "name of the document to be rendered as 'document.<name>'"},
	))
}
