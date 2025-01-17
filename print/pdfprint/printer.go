package pdfprint

import (
	"bytes"
	"context"
	"fmt"
	"image/color"
	"io"
	"log/slog"

	pdf "github.com/stephenafamo/goldmark-pdf"
	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/ast"
	"github.com/yuin/goldmark/parser"

	"github.com/blackstork-io/fabric/plugin"
	"github.com/blackstork-io/fabric/plugin/ast/astsrc"
	"github.com/blackstork-io/fabric/print"
	"github.com/blackstork-io/fabric/print/mdprint"
)

// Printer is the interface for printing pdf content.
type Printer struct {
	md mdprint.Printer
}

// New creates a new pdf printer.
func New() Printer {
	return Printer{mdprint.New()}
}

// Print is a helper function to print pdf content to a writer.
func Print(w io.Writer, el plugin.Content) error {
	p := New()
	return p.Print(context.Background(), w, el)
}

func (p Printer) Print(ctx context.Context, w io.Writer, el plugin.Content) (err error) {
	p.removeFrontmatter(el)
	err = print.ReplaceNodesInContent(el, func(src *astsrc.ASTSource, n ast.Node) (repl ast.Node, err error) {
		switch n := n.(type) {
		case *ast.HTMLBlock:
			slog.Info("HTML block found in AST, replacing with message segment")
			p := ast.NewCodeBlock()
			p.AppendChild(p, ast.NewRawTextSegment(
				src.Appendf("<node of type %q is not supported by the pdf renderer>", n.Kind()),
			))
			return p, nil
		case *ast.RawHTML:
			slog.Info("Raw HTML found in AST, replacing with message segment")
			p := ast.NewCodeSpan()
			p.AppendChild(p, ast.NewRawTextSegment(
				src.Appendf("<node of type %q is not supported by the pdf renderer>", n.Kind()),
			))
			return p, nil
		}
		return n, nil
	})
	if err != nil {
		return fmt.Errorf("replacement failed: %w", err)
	}

	buf := &bytes.Buffer{}
	if err := p.md.Print(ctx, buf, el); err != nil {
		return err
	}

	md := goldmark.New(
		plugin.BaseMarkdownOptions,
		goldmark.WithParserOptions(
			parser.WithAutoHeadingID(),
		),
		goldmark.WithRenderer(
			pdf.New(
				pdf.WithLinkColor(color.RGBA{
					R: 30,
					G: 30,
					B: 255,
					A: 255,
				}),
				pdf.WithHeadingFont(pdf.GetTextFont("Open Sans", pdf.FontRoboto)),
				pdf.WithBodyFont(pdf.GetTextFont("Open Sans", pdf.FontRoboto)),
				pdf.WithCodeFont(pdf.GetCodeFont("Open Sans", pdf.FontRoboto)),
			),
		),
	)
	return md.Convert(buf.Bytes(), w)
}

func (p Printer) removeFrontmatter(el plugin.Content) bool {
	section, ok := el.(*plugin.ContentSectionOrDoc)
	if !ok {
		return false
	}
	for i, child := range section.Children {
		el, ok := child.(*plugin.ContentElement)
		if !ok {
			continue
		}
		meta := el.Meta()
		if meta != nil && meta.Plugin == "blackstork/builtin" && meta.Provider == "frontmatter" {
			section.Children = append(section.Children[:i], section.Children[i+1:]...)
			return true
		}
	}
	return false
}
