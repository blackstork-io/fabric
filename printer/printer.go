package printer

import (
	"io"

	"github.com/blackstork-io/fabric/plugin"
)

type Printer interface {
	Print(w io.Writer, el plugin.Content) error
}
