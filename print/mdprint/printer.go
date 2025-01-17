package mdprint

import (
	"context"
	"io"

	"github.com/blackstork-io/fabric/plugin/ast"
	"github.com/blackstork-io/fabric/plugin/ast/nodes"
)

// Printer is the interface for printing markdown content.
type Printer struct{}

// New creates a new markdown printer.
func New() Printer {
	return Printer{}
}

// PrintString is a helper function to print markdown content to a string.
func PrintString(node *nodes.Node) string {
	return string(ast.AST2Md(node))
}

// Print is a helper function to print markdown content to a writer.
func Print(w io.Writer, node *nodes.Node) error {
	p := New()
	return p.Print(context.Background(), w, node)
}

func (p Printer) Print(ctx context.Context, w io.Writer, el *nodes.Node) (err error) {
	return p.printContent(w, el)
}

func (p Printer) printContent(w io.Writer, node *nodes.Node) (err error) {
	_, err = w.Write(ast.AST2Md(node))
	return
}
