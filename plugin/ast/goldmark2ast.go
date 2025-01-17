package ast

import (
	"log/slog"

	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/ast"
	"github.com/yuin/goldmark/extension"
	east "github.com/yuin/goldmark/extension/ast"
	"github.com/yuin/goldmark/text"

	"github.com/blackstork-io/fabric/pkg/utils"
	"github.com/blackstork-io/fabric/plugin/ast/nodes"
)

var baseMarkdownOptions = goldmark.WithExtensions(
	extension.Table,
	extension.Strikethrough,
	extension.TaskList,
)

// Markdown2AST converts a markdown string to fabric AST.
func Markdown2AST(source []byte) []*nodes.Node {
	markdown := []byte(source)
	root := goldmark.New(baseMarkdownOptions).Parser().Parse(
		text.NewReader(markdown),
	)
	return Goldmark2AST(root, markdown)
}

// Goldmark2AST converts a goldmark AST to fabric AST.
func Goldmark2AST(root ast.Node, source []byte) []*nodes.Node {
	e := &encoder{
		source: source,
	}
	if root.Kind() == ast.KindDocument {
		return e.encodeChildren(root)
	} else {
		return []*nodes.Node{e.encode(root)}
	}
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
		slog.Warn("Attempted to encode a Document node. This should not happen.")
		node.Content = &nodes.FabricDocument{}
	case *ast.TextBlock:
		node.Content = &nodes.Paragraph{
			IsTextBlock: true,
		}
	case *ast.Paragraph:
		node.Content = &nodes.Paragraph{
			IsTextBlock: false,
		}
	case *ast.Heading:
		node.Content = &nodes.Heading{
			Level: n.Level,
		}
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
	case *ast.List:
		node.Content = &nodes.List{
			// https://spec.commonmark.org/0.31.2/#ordered-list-marker
			Start:  uint32(utils.Clamp(0, n.Start, 1_000_000_000-n.ChildCount())),
			Marker: n.Marker,
		}
	case *ast.ListItem:
		node.Content = &nodes.ListItem{}
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
	case *ast.Link:
		node.Content = &nodes.Link{
			Destination: n.Destination,
			Title:       n.Title,
		}
	case *ast.Image:
		node.Content = &nodes.Image{
			Source: n.Destination,
			Alt:    n.Title,
		}
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
		node.Content = &nodes.Table{
			Alignments: utils.FnMap(n.Alignments, func(a east.Alignment) nodes.Alignment {
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
				return alignment
			}),
		}
	case *east.TableRow:
		node.Content = &nodes.TableRow{}
	case *east.TableCell:
		node.Content = &nodes.TableCell{}

	case *east.Strikethrough:
		node.Content = &nodes.Strikethrough{}
	case *east.TaskCheckBox:
		node.Content = &nodes.TaskCheckbox{
			Checked: n.IsChecked,
		}
	}
	node.AppendChildren(e.encodeChildren(n)...)
	return node
}
