package main

import (
	"fmt"
	"os"

	"github.com/alexflint/go-arg"

	"github.com/blackstork-io/fabric/parser"
	"github.com/blackstork-io/fabric/pkg/diagnostics"
)

// TODO: CD: automatically set to correct version using const overrides in build command
const Version string = "0.0.0-dev"

// TODO: multiple plugin support
type args struct {
	Path     string `arg:"positional,required"    help:"a path to a directory with *.fabric files to evaluate" placeholder:"DIR"`
	Plugins  string `arg:"--pluginDir,required"   help:"a path to a __plugin file__"`
	Document string `arg:"-d,--document,required" help:"the name of the document to render"`
}

func (args) Version() string {
	// TODO: How this product is called once again?)
	return "fabric " + Version
}

// TODO: write some desctiption
func (args) Description() string {
	return "fabric document evaluator"
}

func parseArgs() (args args) {
	p := arg.MustParse(&args)
	stat, err := os.Stat(args.Path)
	switch {
	case err == nil:
		if !stat.IsDir() {
			p.Fail("Supplied path does not refer to a directory")
		}
	case os.IsNotExist(err):
		// err
		p.Fail("Supplied path does not exist")
	default:
		// err
		p.Fail(fmt.Sprintf("Path error: %s", err))
	}
	// TODO: path to plugins
	return
}

func newRun() (diags diagnostics.Diag) {
	args := parseArgs()
	result := parser.ParseDir(args.Path)
	diags = result.Diags
	defer func() { diagnostics.PrintDiags(diags, result.FileMap) }()
	if diags.HasErrors() {
		return
	}
	if len(result.FileMap) == 0 {
		diags.Add(
			"No correct fabric files found",
			fmt.Sprintf("There are no *.fabric files at '%s' or all of them have failed to parse", args.Path),
		)
	}

	doc, found := result.Blocks.Documents[args.Document]
	if !found {
		diags.Add(
			"Document not found",
			fmt.Sprintf(
				"Definition for document named '%s' not found in '%s/**.fabric' files",
				args.Document,
				args.Path,
			),
		)
	}

	eval := parser.NewEvaluator(&parser.MockCaller{}, result.Blocks)
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
