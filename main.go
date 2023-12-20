package main

import (
	"flag"
	"fmt"

	"os"

	"github.com/hashicorp/hcl/v2/gohcl"
)

var path, pluginPath, docName string

type Decoder struct {
	root    *Templates
	plugins *Plugins
}

func argParse() error {
	flag.StringVar(&path, "path", "", "a path to a directory with *.hcl files")
	flag.StringVar(&pluginPath, "plugins", "", "a path to a __plugin file__")
	flag.StringVar(&docName, "document", "", "the name of the document to process")
	flag.Parse()
	if path == "" {
		return fmt.Errorf("path required")
	}
	if docName == "" {
		return fmt.Errorf("document name required")
	}
	if pluginPath == "" {
		return fmt.Errorf("plugins required")
	}
	return nil
}

func run() error {
	err := argParse()
	if err != nil {
		return err
	}
	body, fileMap, diags := fromDisk()
	defer func() { PrintDiags(diags, fileMap) }()
	if diags.HasErrors() {
		return diags
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
	fmt.Println(output)
	return nil
}

func main() {
	if err := run(); err != nil {
		fmt.Fprintf(os.Stderr, "error: %s\n", err)
		os.Exit(1)
	}
}
