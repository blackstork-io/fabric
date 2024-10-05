package nodes

import (
	"github.com/yuin/goldmark/ast"
	"google.golang.org/protobuf/types/known/anypb"
)

var (
	ContentNodeKind  = ast.NewNodeKind("FabricContentNode")
	CustomBlockKind  = ast.NewNodeKind("FabricCustomBlock")
	CustomInlineKind = ast.NewNodeKind("FabricCustomInline")
)

func NewCustomNode(isInline bool, data *anypb.Any) ast.Node {
	if isInline {
		return &CustomInline{
			Data: data,
		}
	} else {
		return &CustomBlock{
			Data: data,
		}
	}
}

type CustomInline struct {
	ast.BaseInline
	Data *anypb.Any
}

var _ ast.Node = &CustomInline{}

// Kind implements ast.Node.
func (o *CustomInline) Kind() ast.NodeKind {
	return CustomInlineKind
}

// Dump implements ast.Node.
func (o *CustomInline) Dump(source []byte, level int) {
	var kv map[string]string
	if o == nil || o.Data == nil {
		kv = map[string]string{
			"other": "nil",
		}
	} else {
		kv = map[string]string{
			"other.TypeUrl": o.Data.GetTypeUrl(),
		}
	}
	ast.DumpHelper(o, source, level, kv, nil)
}

type CustomBlock struct {
	ast.BaseBlock
	Data *anypb.Any
}

var _ ast.Node = &CustomBlock{}

// Kind implements ast.Node.
func (o *CustomBlock) Kind() ast.NodeKind {
	return CustomBlockKind
}

// Dump implements ast.Node.
func (o *CustomBlock) Dump(source []byte, level int) {
	var kv map[string]string
	if o == nil || o.Data == nil {
		kv = map[string]string{
			"other": "nil",
		}
	} else {
		kv = map[string]string{
			"other.TypeUrl": o.Data.GetTypeUrl(),
		}
	}
	ast.DumpHelper(o, source, level, kv, nil)
}
