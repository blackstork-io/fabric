package ast

import (
	"bytes"

	"github.com/blackstork-io/fabric/plugin/ast/internal/ast2md"
	"github.com/blackstork-io/fabric/plugin/ast/nodes"
)

// AST2Md converts a fabric AST to a markdown document.
func AST2Md(root *nodes.Node) []byte {
	lines := (&ast2md.Renderer{}).Render(root)
	lines.TrimEmptyLines().SetMarginBottom(1)
	var buf bytes.Buffer
	lines.WriteTo(&buf)
	return buf.Bytes()
}
