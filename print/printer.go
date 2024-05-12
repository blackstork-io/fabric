package print

import (
	"context"
	"io"

	"github.com/blackstork-io/fabric/plugin"
)

// Printer is the interface for printing content.
type Printer interface {
	Print(ctx context.Context, w io.Writer, el plugin.Content) error
}
