package main

import (
	"bufio"
	"os"

	"github.com/hashicorp/hcl/v2"
	"golang.org/x/term"

	"github.com/blackstork-io/fabric/pkg/diagnostics"
)

func PrintDiags(diags diagnostics.Diag, fileMap map[string]*hcl.File) {
	if len(diags) == 0 {
		return
	}

	colorize := term.IsTerminal(0)
	width, _, err := term.GetSize(0)
	if err != nil || width <= 0 {
		width = 80
	}
	wr := bufio.NewWriter(os.Stderr)
	_ = hcl.NewDiagnosticTextWriter(wr, fileMap, uint(width), colorize).WriteDiagnostics(hcl.Diagnostics(diags))
	wr.Flush()
}
