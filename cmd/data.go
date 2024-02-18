package cmd

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"os"
	"regexp"

	"github.com/spf13/cobra"
	"golang.org/x/exp/slices"

	"github.com/blackstork-io/fabric/parser"
	"github.com/blackstork-io/fabric/parser/definitions"
	"github.com/blackstork-io/fabric/pkg/diagnostics"
	"github.com/blackstork-io/fabric/plugin"
)

var dataTgtRe = regexp.MustCompile(`document\.([^.]+)\.data\.([^.]+)\.([^.]+)`)

func Data(ctx context.Context, blocks *parser.DefinedBlocks, caller *parser.Caller, target string) (result plugin.MapData, diags diagnostics.Diag) {
	// docName, pluginName, blockName
	// target: document.<doc-name>.data.<plugin-name>.<data-name>
	tgt := dataTgtRe.FindStringSubmatch(target)
	if tgt == nil {
		diags.Add(
			"Incorrect target",
			"Target should have the format 'document.<doc-name>.data.<plugin-name>.<block-name>'",
		)
		return
	}

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
			fmt.Sprintf("Data block '%s.%s' not found", tgt[2], tgt[3]),
		)
		return
	}
	data := pd.Data[idx]
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

		bw := bufio.NewWriter(os.Stdout)
		defer func() {
			diags.AppendErr(bw.Flush(), "Failed to print the result")
		}()
		enc := json.NewEncoder(bw)
		enc.SetIndent("", "    ")
		diags.AppendErr(
			enc.Encode(res.Any()),
			"Failed to encode the json",
		)
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
