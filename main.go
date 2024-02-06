package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/blackstork-io/fabric/internal/builtin"
	"github.com/blackstork-io/fabric/parser"
	"github.com/blackstork-io/fabric/pkg/diagnostics"
	"github.com/blackstork-io/fabric/plugin/runner"
)

// TODO: replace flag with a better parser (argparse).
var path, pluginPath, docName string

// set using goreleaser by default
var version string

func argParse() (diags diagnostics.Diag) {
	flag.StringVar(&path, "path", "", "a path to a directory with *.hcl files")
	flag.StringVar(&pluginPath, "plugins", "", "a path to a __plugin file__")
	flag.StringVar(&docName, "document", "", "the name of the document to process")
	flag.Parse()
	if path == "" {
		diags.Add("Wrong usage", "path required")
	}
	if docName == "" {
		diags.Add("Wrong usage", "document name required")
	}
	if pluginPath == "" {
		diags.Add("Wrong usage", "plugins required")
	}
	return
}

func newRun() (diags diagnostics.Diag) {
	if diags.Extend(argParse()) {
		return
	}
	result := parser.ParseDir(path)
	diags = result.Diags
	defer func() { diagnostics.PrintDiags(diags, result.FileMap) }()
	if diags.HasErrors() {
		return
	}
	if len(result.FileMap) == 0 {
		diags.Add(
			"No correct fabric files found",
			fmt.Sprintf("There are no *.fabric files at '%s' or all of them have failed to parse", path),
		)
	}

	doc, found := result.Blocks.Documents[docName]
	if !found {
		diags.Add(
			"Document not found",
			fmt.Sprintf(
				"Definition for document named '%s' not found in '%s/**.fabric' files",
				docName,
				path,
			),
		)
	}
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
	if diags.ExtendHcl(stdDiag) {
		return
	}
	defer func() { diags.ExtendHcl(runner.Close()) }()

	caller := parser.NewPluginCaller(runner)

	eval := parser.NewEvaluator(caller, result.Blocks)
	str, diag := eval.EvaluateDocument(doc)
	if diags.Extend(diag) {
		return
	}
	fmt.Printf("Document result:\n%s\n", str)
	return
}

func main() {
	if diags := newRun(); diags.HasErrors() {
		os.Exit(1)
	}
}
