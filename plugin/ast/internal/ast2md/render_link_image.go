package ast2md

import (
	"bytes"

	"github.com/blackstork-io/fabric/plugin/ast/nodes"
)

func (r *Renderer) renderLinkOrImage(n nodes.LinkOrImage, children []*nodes.Node) (lines linesSlice) {
	// [text](dest "title")
	// ![text](dest "title")
	linkText := r.renderNodes(children)

	dest := bytes.Trim(n.Url(), "<> \t")
	angleBrackets := len(dest) == 0 || bytes.ContainsFunc(dest, func(r rune) bool {
		return r == ' ' || r <= 0x1f || r == 0x7f
	})

	if angleBrackets {
		dest = bytes.ReplaceAll(dest, newLine, space)
		dest = angleBracketsEscapeRe.ReplaceAll(dest, []byte("$1\\$2"))
	} else {
		dest = parensEscapeRe.ReplaceAll(dest, []byte("$1\\$2"))
	}

	if _, ok := n.(*nodes.Image); ok {
		lines.Append(bang)
	}

	lines.Append(bracketOpen).
		Extend(*linkText.JoinLines(space)).
		Append(bracketClose, parenOpen)

	if angleBrackets {
		lines.Append(angleOpen)
	}
	lines.Append(dest)
	if angleBrackets {
		lines.Append(angleClose)
	}

	title := bytes.ReplaceAll(n.TitleOrAlt(), newLine, space)
	title = quoteEscapeRe.ReplaceAll(title, []byte("$1\\$2"))
	if len(title) > 0 {
		lines.Append(space, quote, title, quote)
	}
	lines.Append(parenClose)
	return
}
