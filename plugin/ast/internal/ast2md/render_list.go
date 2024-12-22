package ast2md

import (
	"bytes"
	"fmt"

	"github.com/blackstork-io/fabric/plugin/ast/nodes"
)

func (r *Renderer) renderList(c *nodes.List) (lines linesSlice) {
	var spacer, marker []byte
	isOrdered := c.Marker == '.' || c.Marker == ')'
	if isOrdered {
		// generate a spacer that is the same length as the longest(last) marker
		spacer = bytes.Repeat(
			space,
			intLen(int(c.Start)+len(c.Items)-1)+2,
		)
	} else {
		spacer = []byte("  ")
		marker = []byte{c.Marker, ' '}
	}
	for lineNo, item := range c.Items {
		itemLines := r.renderNodes(item)
		if isOrdered {
			marker = fmt.Appendf(nil, "%d%c ", int(c.Start)+lineNo, c.Marker)
		}
		itemLines.
			TrimEmptyLines().
			SetMarginBottom(1).
			PrependPrefixes(marker, spacer[:len(marker)])

		lines.Extend(itemLines)
	}
	lines.SetMinMarginBottom(1)
	return
}

// intLen returns the length of the decimal representation of n.
func intLen(n int) (length int) {
	length = 1
	if n < 0 {
		length++
		n /= -10
	} else {
		n /= 10
	}
	for n != 0 {
		n /= 10
		length++
	}
	return
}
