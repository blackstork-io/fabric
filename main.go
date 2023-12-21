package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/gohcl"

	"github.com/blackstork-io/fabric/pkg/diagnostics"
)

// TODO: replace flag with a better parser (argparse).
var path, pluginPath, docName string

type Decoder struct {
	root    *Templates
	plugins *Plugins
}

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

func run() (diags diagnostics.Diag) {
	var fileMap map[string]*hcl.File
	defer func() { PrintDiags(diags, fileMap) }()

	diag := argParse()
	if diags.Extend(diag) {
		return
	}

	body, fileMap, diag := fromDisk()
	if diags.Extend(diag) {
		return
	}

	plugins, pluginDiag := NewPlugins(pluginPath)

	if diags.Extend(pluginDiag) {
		return diags
	}

	defer plugins.Kill()
	d := Decoder{
		root:    &Templates{},
		plugins: plugins,
	}
	if diags.ExtendHcl(gohcl.DecodeBody(body, nil, d.root)) {
		return diags
	}

	if diags.ExtendHcl(d.Decode()) {
		return diags
	}

	output, diag := d.Evaluate(docName)
	if diag.HasErrors() {
		return diags
	}
	fmt.Println(output) //nolint: forbidigo
	return nil
}

func main() {
	if diags := run(); diags.HasErrors() {
		os.Exit(1)
	}
}
