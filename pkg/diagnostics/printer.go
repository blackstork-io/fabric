package diagnostics

import (
	"bufio"
	"os"

	"github.com/hashicorp/hcl/v2"
	"golang.org/x/term"
)

func PrintDiags(diags Diag, fileMap map[string]*hcl.File) {
	if len(diags) == 0 {
		return
	}

	colorize := term.IsTerminal(0)
	width, _, err := term.GetSize(0)
	if err != nil || width <= 0 {
		width = 80
	}
	wr := bufio.NewWriter(os.Stderr)
	diagWriter := hcl.NewDiagnosticTextWriter(wr, fileMap, uint(width), colorize)

	for _, diag := range diags {
		if _, isRepeated := hcl.DiagnosticExtra[repeatedError](diag); isRepeated {
			continue
		}
		diagWriter.WriteDiagnostic(diag)
	}
	wr.Flush()
}
