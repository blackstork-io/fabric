package diagnostics

import (
	"bufio"
	"io"
	"log/slog"
	"os"

	"github.com/hashicorp/hcl/v2"
	"github.com/mattn/go-colorable"
	"golang.org/x/term"
)

func PrintDiags[D interface{ ~[]*hcl.Diagnostic }](output io.Writer, diags D, fileMap map[string]*hcl.File, colorize bool) {
	printDiags(output, Diag(diags), fileMap, colorize)
}

func printDiags(output io.Writer, diags Diag, fileMap map[string]*hcl.File, colorize bool) {
	if len(diags) == 0 {
		return
	}
	width := 80
	if file, ok := output.(*os.File); ok {
		term_width, _, err := term.GetSize(int(file.Fd()))
		if err == nil && term_width > 0 {
			width = term_width
		}
		if colorize {
			output = colorable.NewColorable(file)
		}
	}

	bufWr := bufio.NewWriter(output)
	defer func() {
		err := bufWr.Flush()
		if err != nil {
			slog.Error("Failed to flush diagnostics", "err", err)
		}
	}()

	diagWriter := hcl.NewDiagnosticTextWriter(bufWr, fileMap, uint(width), colorize)

	for _, diag := range diags {
		if _, isRepeated := hcl.DiagnosticExtra[repeatedError](diag); isRepeated {
			continue
		}
		err := diagWriter.WriteDiagnostic(diag)
		if err != nil {
			slog.Error("Failed to write diagnostics", "err", err)
		}
	}
}
