package print

import (
	"context"
	"errors"
	"fmt"
	"io"

	"github.com/yuin/goldmark/ast"

	"github.com/blackstork-io/fabric/plugin/ast/nodes"
)

// Printer is the interface for printing content.
type Printer interface {
	Print(ctx context.Context, w io.Writer, el *nodes.Node) error
}

var (
	// ErrReplacerSkipChildren could be returned from the ReplaceNodes replacer func
	// to skip children. Will never be returned by ReplaceNodes itself.
	ErrReplacerSkipChildren = errors.New("skip children")
	// ErrReplacerStuck means that replacer failed to make progress
	ErrReplacerStuck = errors.New("replacer stuck")
)

// ReplaceNodes walks the AST starting from the given node and replaces nodes in it.
// If replacer returns nil - the node is deleted
func ReplaceNodes(n ast.Node, replacer func(n ast.Node) (repl ast.Node, err error)) (ast.Node, error) {
	const maxReplacementsWithoutAdvance = 100
	if n == nil {
		return nil, nil
	}
	n, err := replacer(n)
	if n == nil || err != nil {
		if err == ErrReplacerSkipChildren {
			err = nil
		}
		return n, err
	}
	c := n.FirstChild()
	replacementsWithoutAdvance := 0
	for c != nil {
		repl, err := ReplaceNodes(c, replacer)
		if err != nil {
			return n, err
		}
		switch repl {
		case nil:
			next := c.NextSibling()
			n.RemoveChild(n, c)
			c = next
			replacementsWithoutAdvance = 0
		case c:
			c = c.NextSibling()
			replacementsWithoutAdvance = 0
		default:
			if replacementsWithoutAdvance >= maxReplacementsWithoutAdvance {
				return n, fmt.Errorf("%w: node %q", ErrReplacerStuck, repl.Kind())
			}
			replacementsWithoutAdvance++
			n.ReplaceChild(n, c, repl)
			// Intentionally trying to replace the replacement result
			// This allows replacer to not care whether the replacement node
			// or its children need to be further replaced
			c = repl
		}
	}
	return n, nil
}
