package ast

import (
	"log/slog"

	"github.com/yuin/goldmark/ast"
	east "github.com/yuin/goldmark/extension/ast"
	"github.com/yuin/goldmark/text"

	"github.com/blackstork-io/fabric/pkg/utils"
	"github.com/blackstork-io/fabric/plugin/ast/nodes"
)

// Goldmark2AST converts a goldmark AST to fabric AST.
func Goldmark2AST(root ast.Node, source []byte) (node *nodes.Node) {
	return (&encoder{
		source: source,
	}).encode(root)
}

type encoder struct {
	source []byte
}

func (e *encoder) encodeSegments(seg *text.Segments) []byte {
	var res []byte
	for i := range seg.Len() {
		res = append(res, e.encodeSegment(seg.At(i))...)
	}
	return res
}

func (e *encoder) encodeSegment(seg text.Segment) []byte {
	if seg.Start > seg.Stop || seg.Start < 0 || seg.Stop > len(e.source) {
		// Invalid segment.
		return nil
	}
	// Note: zero-length segments (seg.Start == seg.Stop) are valid,
	// and even could be meaningful if padding is > 0.
	// Render them using seg.Value.
	return seg.Value(e.source)
}

func (e *encoder) encodeChildren(n ast.Node) []*nodes.Node {
	if n.ChildCount() == 0 {
		return nil
	}
	children := make([]*nodes.Node, 0, n.ChildCount())
	for c := n.FirstChild(); c != nil; c = c.NextSibling() {
		children = append(children, e.encode(c))
	}
	return children
}

func (e *encoder) encode(n ast.Node) *nodes.Node {
	node := &nodes.Node{}

	switch n := n.(type) {
	case *ast.Document:
		node.Content = &nodes.Document{}
		node.AppendChildren(e.encodeChildren(n)...)
	case *ast.TextBlock:
		node.Content = &nodes.Paragraph{
			IsTextBlock: true,
		}
		node.AppendChildren(e.encodeChildren(n)...)
	case *ast.Paragraph:
		node.Content = &nodes.Paragraph{
			IsTextBlock: false,
		}
		node.AppendChildren(e.encodeChildren(n)...)
	case *ast.Heading:
		node.Content = &nodes.Heading{
			Level: n.Level,
		}
		node.AppendChildren(e.encodeChildren(n)...)
	case *ast.ThematicBreak:
		node.Content = &nodes.ThematicBreak{}
	case *ast.CodeBlock:
		node.Content = &nodes.CodeBlock{
			Language: nil,
			Code:     e.encodeSegments(n.Lines()),
		}
	case *ast.FencedCodeBlock:
		content := &nodes.CodeBlock{
			Code: e.encodeSegments(n.Lines()),
		}
		if n.Info == nil {
			content.Language = []byte{}
		} else {
			content.Language = e.encodeSegment(n.Info.Segment)
		}
		node.Content = content
	case *ast.Blockquote:
		node.Content = &nodes.Blockquote{}
		node.AppendChildren(e.encodeChildren(n)...)
	case *ast.List:
		content := &nodes.List{
			// https://spec.commonmark.org/0.31.2/#ordered-list-marker
			Start:  uint32(utils.Clamp(0, n.Start, 1_000_000_000 - n.ChildCount())),
			Marker: n.Marker,
			Items:  make([][]*nodes.Node, 0, n.ChildCount()),
		}
		for item := n.FirstChild(); item != nil; item = item.NextSibling() {
			if item, ok := item.(*ast.ListItem); ok {
				itemContent := make([]*nodes.Node, 0, item.ChildCount())
				for i := item.FirstChild(); i != nil; i = i.NextSibling() {
					itemContent = append(itemContent, e.encode(i))
				}
				content.Items = append(content.Items, itemContent)
			} else {
				// should not happen in a well-formed AST
				content.Items = append(content.Items, []*nodes.Node{e.encode(item)})
			}
		}
		node.Content = content
	case *ast.HTMLBlock:
		content := &nodes.HTMLBlock{
			HTML: e.encodeSegments(n.Lines()),
		}
		if !n.ClosureLine.IsEmpty() {
			content.HTML = append(content.HTML, e.encodeSegment(n.ClosureLine)...)
		}
		node.Content = content
	case *ast.Text:
		text := e.encodeSegment(n.Segment)
		if n.SoftLineBreak() {
			text = append(text, '\n')
		}

		node.Content = &nodes.Text{
			Text:          text,
			HardLineBreak: n.HardLineBreak(),
		}
	case *ast.String:
		node.Content = &nodes.Text{
			Text: n.Value,
		}
	case *ast.CodeSpan:
		content := &nodes.CodeSpan{}
		for c := n.FirstChild(); c != nil; c = c.NextSibling() {
			switch c := c.(type) {
			case *ast.Text:
				content.Code = append(content.Code, e.encodeSegment(c.Segment)...)
			case *ast.String:
				content.Code = append(content.Code, c.Value...)
			default:
				// should not happen in a well-formed AST
				slog.Warn("unexpected child of code span", "kind", c.Kind().String())
			}
		}
		node.Content = content
	case *ast.Emphasis:
		node.Content = &nodes.Emphasis{
			Level: n.Level,
		}
		node.AppendChildren(e.encodeChildren(n)...)
	case *ast.Link:
		node.Content = &nodes.Link{
			Destination: n.Destination,
			Title:       n.Title,
		}
		node.AppendChildren(e.encodeChildren(n)...)
	case *ast.Image:
		node.Content = &nodes.Image{
			Source: n.Destination,
			Alt:    n.Title,
		}
		node.AppendChildren(e.encodeChildren(n)...)
	case *ast.AutoLink:
		content := &nodes.AutoLink{
			Value: n.URL(e.source),
		}

		node.Content = content
	case *ast.RawHTML:
		node.Content = &nodes.HTMLInline{
			HTML: e.encodeSegments(n.Segments),
		}
	case *east.Table:
		content := &nodes.Table{
			Alignments: make([]nodes.Alignment, 0, len(n.Alignments)),
		}
		for _, a := range n.Alignments {
			var alignment nodes.Alignment
			switch a {
			case east.AlignLeft:
				alignment = nodes.AlignmentLeft
			case east.AlignRight:
				alignment = nodes.AlignmentRight
			case east.AlignCenter:
				alignment = nodes.AlignmentCenter
			case east.AlignNone:
				alignment = nodes.AlignmentNone
			default:
				slog.Warn("unexpected table alignment", "alignment", a)
				alignment = nodes.AlignmentNone
			}
			content.Alignments = append(content.Alignments, alignment)
		}
		content.Cells = make([][][]*nodes.Node, 0, n.ChildCount())
		for c := n.FirstChild(); c != nil; c = c.NextSibling() {
			if c.Kind() != east.KindTableRow && c.Kind() != east.KindTableHeader {
				slog.Warn("unexpected child of table", "kind", c.Kind().String())
				continue
			}
			row := make([][]*nodes.Node, 0, c.ChildCount())
			for cell := c.FirstChild(); cell != nil; cell = cell.NextSibling() {
				if cell.Kind() != east.KindTableCell {
					slog.Warn("unexpected child of table row", "kind", cell.Kind().String())
					continue
				}
				cellContent := make([]*nodes.Node, 0, cell.ChildCount())
				for i := cell.FirstChild(); i != nil; i = i.NextSibling() {
					cellContent = append(cellContent, e.encode(i))
				}
				row = append(row, cellContent)
			}
			content.Cells = append(content.Cells, row)
		}
		node.Content = content
		return node
	case *east.Strikethrough:
		node.Content = &nodes.Strikethrough{}
		node.AppendChildren(e.encodeChildren(n)...)
	case *east.TaskCheckBox:
		node.Content = &nodes.TaskCheckbox{
			Checked: n.IsChecked,
		}
	}
	return node
}
