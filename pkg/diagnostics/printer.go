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

func PrintDiags(output *os.File, diags Diag, fileMap map[string]*hcl.File, colorize bool) {
	if len(diags) == 0 {
		return
	}

	width, _, err := term.GetSize(int(output.Fd()))
	if err != nil || width <= 0 {
		width = 80
	}

	var wr io.Writer = output
	if colorize {
		wr = colorable.NewColorable(output)
	}
	bufWr := bufio.NewWriter(wr)
	defer func() {
		err := bufWr.Flush()
		if err != nil {
			slog.Error("Failed to flush diagnostics", "err", err)
		}
	}()

	diagWriter := hcl.NewDiagnosticTextWriter(bufWr, fileMap, uint(width), colorize)

	for _, diag := range diags {
		if _, isHidden := hcl.DiagnosticExtra[hiddenErrorIface](diag); isHidden {
			continue
		}
		// TODO: catch traceback diag here and format it properly
		err := diagWriter.WriteDiagnostic(diag)
		if err != nil {
			slog.Error("Failed to write diagnostics", "err", err)
		}
	}
}
