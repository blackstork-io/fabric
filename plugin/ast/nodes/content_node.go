package nodes

import "github.com/yuin/goldmark/ast"

type FabricContentNode struct {
	ast.BaseBlock
	Meta *ContentMeta
}

func ToFabricContentNode(node ast.Node) (meta *FabricContentNode) {
	switch n := node.(type) {
	case *FabricContentNode:
		meta = n
	case *ast.Document:
		meta = &FabricContentNode{}
		child := n.FirstChild()
		for child != nil {
			c := child
			child = child.NextSibling()
			meta.AppendChild(meta, c)
		}
	case nil:
		// meta is nil
	default:
		meta = &FabricContentNode{}
		meta.AppendChild(meta, n)
	}
	return
}

// Dump implements ast.Node.
func (m *FabricContentNode) Dump(source []byte, level int) {
	var kv map[string]string
	if m == nil || m.Meta == nil {
		kv = map[string]string{
			"meta": "nil",
		}
	} else {
		kv = map[string]string{
			"meta.provider": m.Meta.Provider,
			"meta.plugin":   m.Meta.Plugin,
			"meta.version":  m.Meta.Version,
		}
	}
	ast.DumpHelper(m, source, level, kv, nil)
}

// Kind implements ast.Node.
func (m *FabricContentNode) Kind() ast.NodeKind {
	return ContentNodeKind
}

var _ ast.Node = &FabricContentNode{}
