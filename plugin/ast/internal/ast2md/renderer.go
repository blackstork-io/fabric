package ast2md

import (
	"bytes"
	"fmt"
	"regexp"

	"github.com/blackstork-io/fabric/pkg/utils"
	"github.com/blackstork-io/fabric/plugin/ast/nodes"
)

const unescapedRePart = `((?:^|[^\\])(?:\\{2})*)`

var (
	punctuationRe         = regexp.MustCompile(`([\\_~*#\-=+>|&\x60])`) // \x60 is backtick (`)
	angleBracketsEscapeRe = regexp.MustCompile(unescapedRePart + `([<>])`)
	parensEscapeRe        = regexp.MustCompile(unescapedRePart + `([()])`)
	quoteEscapeRe         = regexp.MustCompile(unescapedRePart + `(")`)
	pipeEscapeRe          = regexp.MustCompile(unescapedRePart + `(|)`)

	autolinkRe = regexp.MustCompile(`[\x00-\x1f\x7f <>]`)
)

type Renderer struct {
	// render context
	inTable int
}

func (r *Renderer) renderNodes(nodes []*nodes.Node) (lines linesSlice) {
	for _, c := range nodes {
		lines.Extend(r.Render(c))
	}
	return lines
}

func (r *Renderer) renderChildren(n *nodes.Node) (lines linesSlice) {
	if n == nil {
		return
	}
	return r.renderNodes(n.GetChildren())
}

var (
	space            = []byte(" ")
	colon            = []byte(":")
	dash             = []byte("-")
	pipe             = []byte("|")
	backtick         = []byte("`")
	newLine          = []byte("\n")
	blockquotePrefix = []byte("> ")
	starEmphasis     = []byte("**")
	bracketOpen      = []byte("[")
	bracketClose     = []byte("]")
	parenOpen        = []byte("(")
	parenClose       = []byte(")")
	angleOpen        = []byte("<")
	angleClose       = []byte(">")
	bang             = []byte("!")
	quote            = []byte("\"")
	backslash        = []byte("\\")
	atxHeading       = []byte("###### ")
	setextHeading    = []byte("=-")
	strikethrough    = []byte("~~")
	thematicBreak    = []byte("_____") // other thematic break styles can interfere with list markers
	hardBreak        = []byte("\\\n")
	checkboxFilled   = []byte("[x] ")
	checkboxEmpty    = []byte("[ ] ")
)

func (r *Renderer) Render(n *nodes.Node) (lines linesSlice) {
	switch c := n.Content.(type) {
	case *nodes.Paragraph:
		lines = r.renderChildren(n)
		lines.SetMinMarginTop(2).SetMinMarginBottom(1)
		return
	case *nodes.Text:
		lines.AppendLines(c.Text)
		if c.HardLineBreak {
			lines.AppendHardBreak()
		}
	case *nodes.Emphasis:
		// TODO: markdown's rules for parsing emphasis are way more complex than this
		// figure out how to produce the correct output in more cases
		// (not possible in all cases, but we can try, or replace with html)
		lines = r.renderChildren(n)
		if lines.IsEmpty() {
			return
		}
		lines.Surround(starEmphasis[:c.Level])
	case *nodes.Strikethrough:
		lines = r.renderChildren(n)
		if lines.IsEmpty() {
			return
		}
		lines.Surround(strikethrough)
	case *nodes.Heading:
		lines = r.renderChildren(n)
		lines.RemoveEmptyLines()
		if lines.IsEmpty() {
			return
		}
		level := utils.Clamp(1, c.Level, 6)
		if len(lines) > 1 && level <= 2 {
			// setext-style heading
			lines.ClearPrefixes().SetMarginBottom(1).Append(bytes.Repeat(
				setextHeading[level-1:level],
				lines.MaxLineLength(),
			)).SetMarginTop(2).SetMarginBottom(1)
			return
		}
		lines.JoinLines(space).
			SetMarginTop(2).
			SetMarginBottom(2).
			PrependPrefix(atxHeading[6-level:])

	case *nodes.ThematicBreak:
		lines.Append(thematicBreak).SetMarginBlock(1)
	case *nodes.CodeBlock:
		lines = r.renderCodeBlock(c)
	case *nodes.Blockquote:
		lines = r.renderChildren(n)
		lines.PrependPrefix(blockquotePrefix)
	case *nodes.List:
		lines = r.renderList(c, n.GetChildren())
	case *nodes.ListItem:
		lines = r.renderChildren(n)
	case *nodes.TaskCheckbox:
		if c.Checked {
			lines.Append(checkboxFilled)
		} else {
			lines.Append(checkboxEmpty)
		}
	case *nodes.HTMLBlock:
		lines.AppendLines(c.HTML).
			SetMinMarginBlock(2)
	case *nodes.HTMLInline:
		lines.AppendLines(c.HTML).JoinLines(space)

	case *nodes.CodeSpan:
		lines = r.renderCodeSpan(c)

	case *nodes.Link:
		lines = r.renderLinkOrImage(c, n.GetChildren())

	case *nodes.Image:
		lines = r.renderLinkOrImage(c, n.GetChildren())

	case *nodes.AutoLink:
		lines.Append(
			angleOpen,
			autolinkRe.ReplaceAllFunc(c.Value, func(b []byte) []byte {
				return fmt.Appendf(nil, "%%%02X", b[0])
			}),
			angleClose,
		)

	case *nodes.Table:
		lines = r.renderTable(c, n.GetChildren())
	case *nodes.FabricDocument:
		lines = r.renderChildren(n)
		lines.TrimEmptyLines().SetMarginBottom(1)
	case *nodes.FabricSection:
		lines = r.renderChildren(n)
		lines.TrimEmptyLines().SetMarginBottom(2)
	case *nodes.FabricContent:
		lines = r.renderChildren(n)
	case *nodes.Custom:
		lines.Append(fmt.Appendf(nil, "<!-- custom node %q not rendered -->", c.Data.GetTypeUrl()))
	}
	return
}
