package cmd

import (
	"fmt"
	"log/slog"
	"os"
	"strings"

	"github.com/spf13/cobra"
	"go.opentelemetry.io/otel/attribute"

	"github.com/blackstork-io/fabric/engine"
	"github.com/blackstork-io/fabric/internal/builtin"
	"github.com/blackstork-io/fabric/parser/definitions"
	"github.com/blackstork-io/fabric/pkg/diagnostics"
	"github.com/blackstork-io/fabric/plugin"
	"github.com/blackstork-io/fabric/print"
	"github.com/blackstork-io/fabric/print/htmlprint"
	"github.com/blackstork-io/fabric/print/mdprint"
)

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

var renderCmd = &cobra.Command{
	Use:   "render TARGET",
	Short: "Render the document",
	Long:  `Render the specified document and either publish it or output it to stdout.`,
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
		ctx := cmd.Context()
		logger := slog.Default()

		var diags diagnostics.Diag
		eng := engine.New(
			engine.WithLogger(logger),
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
		diag = eng.LoadPluginResolver(ctx, false)
		if diags.Extend(diag) {
			return diags
		}
		diag = eng.LoadPluginRunner(ctx)
		if diags.Extend(diag) {
			return diags
		}
		var content plugin.Content
		if publish {
			content, diag = eng.RenderAndPublishContent(ctx, target)
		} else {
			_, content, _, diag = eng.RenderContent(ctx, target)
		}
		if diags.Extend(diag) {
			return diags
		}

		// If publish requested, no need to print out to stdout
		if publish {
			return nil
		}

		logger.InfoContext(ctx, "Printing to stdout", "format", format)

		var printer print.Printer
		switch format {
		case "md":
			printer = mdprint.New()
		case "html":
			printer = htmlprint.New()
		default:
			diags.Add("Unsupported format", fmt.Sprintf("Format '%s' is not supported for stdout", format))
			return
		}
		printer = print.WithLogging(printer, slog.Default(), slog.String("format", format))
		printer = print.WithTracing(printer, tracer, attribute.String("format", format))
		err = printer.Print(ctx, os.Stdout, content)
		if err != nil {
			diags.AppendErr(err, "Error while printing")
		}

		// Making sure the stdout printout has a linebreak at the end
		fmt.Printf("\n")

		return nil
	},
}
