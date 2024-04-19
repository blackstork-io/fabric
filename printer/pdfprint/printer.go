package pdfprint

import (
	"bytes"
	"image/color"
	"io"

	pdf "github.com/stephenafamo/goldmark-pdf"
	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/extension"
	"github.com/yuin/goldmark/parser"

	"github.com/blackstork-io/fabric/plugin"
	"github.com/blackstork-io/fabric/printer/mdprint"
)

type Printer struct {
	md *mdprint.Printer
}

func New() *Printer {
	return &Printer{
		md: mdprint.New(),
	}
}

func (p *Printer) Print(w io.Writer, el plugin.Content) error {
	p.removeFrontatter(el)
	buf := bytes.NewBuffer(nil)
	if err := p.md.Print(buf, el); err != nil {
		return err
	}
	md := goldmark.New(
		goldmark.WithExtensions(extension.GFM),
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

func (p *Printer) removeFrontatter(el plugin.Content) bool {
	section, ok := el.(*plugin.ContentSection)
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
