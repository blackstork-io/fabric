// Simple formatter for testing purposes
// 1) Reads markdown from stdin
// 2) Parses it with goldmark
// 3) Converts it to fabric AST
// 4) Renders fabric AST to markdown
// 5) Writes the result to stdout
package main

import (
	"io"
	"os"

	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/extension"
	"github.com/yuin/goldmark/text"

	"github.com/blackstork-io/fabric/plugin/ast"
	"github.com/blackstork-io/fabric/plugin/ast/internal/ast2md"
	"github.com/blackstork-io/fabric/plugin/ast/nodes"
)

func main() {
	data, err := io.ReadAll(os.Stdin)
	if err != nil {
		panic(err)
	}

	root := goldmark.New(goldmark.WithExtensions(
		extension.Table,
		extension.Strikethrough,
		extension.TaskList,
	)).Parser().Parse(
		text.NewReader(data),
	)

	fabRoot := nodes.NewNode(&nodes.FabricDocument{})
	fabRoot.SetChildren(ast.Goldmark2AST(root, data))
	lines := (&ast2md.Renderer{}).Render(fabRoot)
	lines.WriteTo(os.Stdout)
}
