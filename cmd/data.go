package cmd

import (
	"encoding/json"
	"log/slog"
	"os"

	"github.com/TylerBrock/colorjson"
	"github.com/spf13/cobra"

	"github.com/blackstork-io/fabric/engine"
	"github.com/blackstork-io/fabric/internal/builtin"
	"github.com/blackstork-io/fabric/pkg/diagnostics"
)

func init() {
	rootCmd.AddCommand(dataCmd)
	dataCmd.SetUsageTemplate(UsageTemplate(
		[2]string{"TARGET", "a path to the data block to be executed. Data block must be inside of a document, so the path would look lile 'document.<doc-name>.data.<plugin-name>.<data-name>'"},
	))
}

var dataCmd = &cobra.Command{
	Use:   "data TARGET",
	Short: "Execute a single data block",
	Long:  `Execute the data block and print out prettified JSON to stdout`,
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) (err error) {
		ctx := cmd.Context()
		var diags diagnostics.Diag
		eng := engine.New(
			engine.WithLogger(slog.Default()),
			engine.WithTracer(tracer),
			engine.WithBuiltIn(builtin.Plugin(version, slog.Default(), tracer)),
		)
		defer func() {
			err = exitCommand(eng, cmd, diags)
		}()
		diag := eng.ParseDir(ctx, os.DirFS(cliArgs.sourceDir))
		if diags.Extend(diag) {
			return
		}
		if diags.Extend(eng.LoadPluginResolver(ctx, false)) {
			return
		}
		if diags.Extend(eng.LoadPluginRunner(ctx)) {
			return
		}
		res, diag := eng.FetchData(ctx, args[0])
		if diags.Extend(diag) {
			return
		}
		val := res.Any()
		var ser []byte
		if cliArgs.colorize {
			fmt := colorjson.NewFormatter()
			fmt.Indent = 4
			ser, err = fmt.Marshal(val)
		} else {
			ser, err = json.MarshalIndent(val, "", "    ")
		}
		if diags.AppendErr(err, "Failed to serialize data output to json") {
			return
		}
		_, err = os.Stdout.Write(ser)

		diags.AppendErr(err, "Failed to output json data")
		return
	},
}
