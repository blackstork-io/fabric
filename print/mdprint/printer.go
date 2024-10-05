package mdprint

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"log/slog"

	markdown "github.com/blackstork-io/goldmark-markdown"
	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/ast"

	"github.com/blackstork-io/fabric/plugin"
	"github.com/blackstork-io/fabric/plugin/ast/astsrc"
	"github.com/blackstork-io/fabric/plugin/ast/nodes"
	"github.com/blackstork-io/fabric/print"
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
		if content.IsAst() {
			src, node := content.AsNode()
			err = p.printContentElement(w, src, node)
		} else {
			_, err = w.Write(content.AsMarkdownSrc())
		}
	case *plugin.ContentSection:
		err = p.printContentSection(w, content)
	}
	return
}

func (p Printer) printContentElement(w io.Writer, source *astsrc.ASTSource, node *nodes.FabricContentNode) error {
	n, err := print.ReplaceNodes(node, func(n ast.Node) (repl ast.Node, err error) {
		switch nT := n.(type) {
		case *nodes.CustomBlock:
			slog.Info("OtherBlock found in AST, replacing with message segment")
			p := ast.NewCodeBlock()
			p.AppendChild(p, ast.NewRawTextSegment(
				source.Appendf("<node of type %q is not supported by the pdf renderer>", nT.Data.GetTypeUrl()),
			))
			return p, nil
		case *nodes.CustomInline:
			slog.Info("OtherInline found in AST, replacing with message segment")
			p := ast.NewCodeSpan()
			p.AppendChild(p, ast.NewRawTextSegment(
				source.Appendf("<node of type %q is not supported by the pdf renderer>", nT.Data.GetTypeUrl()),
			))
			return p, nil
		}
		return n, nil
	})
	if err != nil {
		return fmt.Errorf("replacement failed: %w", err)
	}
	err = goldmark.New(
		plugin.BaseMarkdownOptions,
		goldmark.WithExtensions(
			markdown.NewRenderer(
				markdown.WithIgnoredNodes(
					nodes.ContentNodeKind,
				),
			),
		),
	).Renderer().Render(w, source.AsBytes(), n)
	if err != nil {
		return fmt.Errorf("failed to render markdown from AST: %w", err)
	}
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
