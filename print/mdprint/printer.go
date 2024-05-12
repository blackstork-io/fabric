package mdprint

import (
	"bytes"
	"context"
	"io"

	"github.com/blackstork-io/fabric/plugin"
)

// Printer is the interface for printing markdown content.
type Printer struct{}

// New creates a new markdown printer.
func New() Printer {
	return Printer{}
}

// PrintString is a helper function to print markdown content to a string.
func PrintString(el plugin.Content) string {
	buf := bytes.NewBuffer(nil)
	if err := Print(buf, el); err != nil {
		return ""
	}
	return buf.String()
}

// Print is a helper function to print markdown content to a writer.
func Print(w io.Writer, el plugin.Content) error {
	p := New()
	return p.Print(context.Background(), w, el)
}

func (p Printer) Print(ctx context.Context, w io.Writer, el plugin.Content) (err error) {
	return p.printContent(w, el)
}

func (p Printer) printContent(w io.Writer, content plugin.Content) (err error) {
	switch content := content.(type) {
	case *plugin.ContentElement:
		if err = p.printContentElement(w, content); err != nil {
			return
		}
	case *plugin.ContentSection:
		if err = p.printContentSection(w, content); err != nil {
			return
		}
	}
	return nil
}

func (p Printer) printContentElement(w io.Writer, elem *plugin.ContentElement) error {
	_, err := w.Write([]byte(elem.Markdown))
	return err
}

func (p Printer) printContentSection(w io.Writer, sec *plugin.ContentSection) (err error) {
	for i, child := range sec.Children {
		if i > 0 {
			if _, err = w.Write([]byte("\n\n")); err != nil {
				return
			}
		}
		if err = p.printContent(w, child); err != nil {
			return
		}
	}
	return err
}
