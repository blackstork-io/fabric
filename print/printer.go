package print

import (
	"context"
	"io"

	"github.com/yuin/goldmark/ast"

	"github.com/blackstork-io/fabric/plugin"
	"github.com/blackstork-io/fabric/plugin/ast/astsrc"
)

// Printer is the interface for printing content.
type Printer interface {
	Print(ctx context.Context, w io.Writer, el plugin.Content) error
}

// ReplaceNodes walks the AST starting from the given node and replaces nodes in it.
// If replacer returns nil - the node is deleted
func ReplaceNodes(n ast.Node, replacer func(n ast.Node) (repl ast.Node, skipChildren bool)) ast.Node {
	if n == nil {
		return nil
	}
	n, skipChildren := replacer(n)
	if n == nil || skipChildren {
		return n
	}
	c := n.FirstChild()
	for c != nil {
		repl := ReplaceNodes(c, replacer)
		switch repl {
		case nil:
			next := c.NextSibling()
			n.RemoveChild(n, c)
			c = next
		case c:
			c = c.NextSibling()
		default:
			n.ReplaceChild(n, c, repl)
			c = repl
		}
	}
	return n
}

// ReplaceNodesInContent runs ReplaceNodes on plugin.ContentAST nodes
func ReplaceNodesInContent(el plugin.Content, replacer func(src *astsrc.ASTSource, n ast.Node) (repl ast.Node, skipChildren bool)) {
	switch el := el.(type) {
	case *plugin.ContentSection:
		for _, child := range el.Children {
			ReplaceNodesInContent(child, replacer)
		}
	case *plugin.ContentElement:
		if !el.IsAst() {
			return
		}
		src, node := el.AsNode()
		ReplaceNodes(node, func(n ast.Node) (repl ast.Node, skipChildren bool) {
			return replacer(src, n)
		})
		node.Dump(src.AsBytes(), 0)
	}
}
