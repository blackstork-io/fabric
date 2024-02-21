package cmd

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"regexp"

	"github.com/TylerBrock/colorjson"
	"github.com/spf13/cobra"
	"golang.org/x/exp/slices"

	"github.com/blackstork-io/fabric/parser"
	"github.com/blackstork-io/fabric/parser/definitions"
	"github.com/blackstork-io/fabric/pkg/diagnostics"
	"github.com/blackstork-io/fabric/plugin"
)

var dataTgtRe = regexp.MustCompile(`(?:document\.([^.]+)\.data\.([^.]+)\.([^.\n]+))|(?:data\.([^.]+)\.([^.]+))`)

func Data(ctx context.Context, blocks *parser.DefinedBlocks, caller *parser.Caller, target string) (result plugin.Data, diags diagnostics.Diag) {
	// docName, pluginName, blockName
	// target: document.<doc-name>.data.<plugin-name>.<data-name>
	tgt := dataTgtRe.FindStringSubmatch(target)
	if tgt == nil {
		diags.Add(
			"Incorrect target",
			"Target should have the format 'document.<doc-name>.data.<plugin-name>.<block-name>' or 'data.<plugin-name>.<block-name>'",
		)
		return
	}

	var data *definitions.ParsedData

	if tgt[1] != "" {
		// document.<doc-name>.data.<plugin-name>.<block-name>
		doc, found := blocks.Documents[tgt[1]]
		if !found {
			diags.Add(
				"Document not found",
				fmt.Sprintf(
					"Definition for document named '%s' not found",
					tgt[1],
				),
			)
			return
		}

		pd, diag := blocks.ParseDocument(doc)
		if diags.Extend(diag) {
			return
		}

		idx := slices.IndexFunc(pd.Data, func(data *definitions.ParsedData) bool {
			return data.PluginName == tgt[2] && data.BlockName == tgt[3]
		})
		if idx == -1 {
			diags.Add(
				"Data block not found",
				fmt.Sprintf("Data block '%s.%s' not found in document '%s'", tgt[2], tgt[3], tgt[1]),
			)
			return
		}
		data = pd.Data[idx]
	} else {
		// data.<plugin-name>.<block-name>
		defPlugin, found := blocks.Plugins[definitions.Key{
			PluginKind: definitions.BlockKindData,
			PluginName: tgt[4],
			BlockName:  tgt[5],
		}]
		if !found {
			diags.Add(
				"Data block not found",
				fmt.Sprintf("Data block '%s.%s' not found in global scope", tgt[4], tgt[5]),
			)
			return
		}
		res, diag := blocks.ParsePlugin(defPlugin)
		if diags.Extend(diag) {
			return
		}
		data = (*definitions.ParsedData)(res)
	}
	res, diag := caller.CallData(ctx, data.PluginName, data.Config, data.Invocation)
	if diags.Extend(diag) {
		return
	}
	return res, diags
}

// dataCmd represents the data command
var dataCmd = &cobra.Command{
	Use:   "data TARGET",
	Short: "Execute a single data block",
	Long:  `Execute the data block and print out prettified JSON to stdout`,
	RunE: func(cmd *cobra.Command, args []string) (err error) {
		var diags diagnostics.Diag
		eval := NewEvaluator(cliArgs.pluginsDir)
		defer func() {
			err = eval.Cleanup(diags)
		}()
		diags = eval.ParseFabricFiles(os.DirFS(cliArgs.sourceDir))
		if diags.HasErrors() {
			return
		}
		if diags.Extend(eval.LoadRunner()) {
			return
		}

		res, diag := Data(cmd.Context(), eval.Blocks, eval.PluginCaller(), args[0])
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
	Args: cobra.ExactArgs(1),
}

func init() {
	rootCmd.AddCommand(dataCmd)

	dataCmd.SetUsageTemplate(UsageTemplate(
		[2]string{"TARGET", "a path to the data block to be executed. Data block must be inside of a document, so the path would look lile 'document.<doc-name>.data.<plugin-name>.<data-name>'"},
	))
}
