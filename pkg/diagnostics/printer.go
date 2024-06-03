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

func PrintDiags(output io.Writer, diags []*hcl.Diagnostic, fileMap map[string]*hcl.File, colorize bool) {
	if len(diags) == 0 {
		return
	}
	width := 80
	if file, ok := output.(*os.File); ok {
		termWidth, _, err := term.GetSize(int(file.Fd()))
		if err == nil && termWidth > 0 {
			width = termWidth
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
		if _, isRepeated := DiagnosticExtra[repeatedError](diag); isRepeated {
			continue
		}
		if gojqErr, ok := DiagnosticExtra[GoJQError](diag); ok {
			gojqErr.improveDiagnostic(diag, fileMap)
		}
		if traceback, ok := DiagnosticExtra[TracebackExtra](diag); ok {
			traceback.improveDiagnostic(diag)
		}
		if path, ok := DiagnosticExtra[PathExtra](diag); ok {
			path.improveDiagnostic(diag)
		}
		err := diagWriter.WriteDiagnostic(diag)
		if err != nil {
			slog.Error("Failed to write diagnostics", "err", err)
		}
	}
}
