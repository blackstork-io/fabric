package tree

import (
	"reflect"

	"github.com/hashicorp/hcl/v2/hclsyntax"
	"github.com/zclconf/go-cty/cty"
)

type NodeSigil struct{}

func (t NodeSigil) isNode() NodeSigil {
	return NodeSigil{}
}

type Node interface {
	isNode() NodeSigil
	Namer
}

var NodeType = cty.Capsule("node", reflect.TypeFor[Node]())

type Namer interface {
	FriendlyName() string
}

// Node can implement following interfaces optionally
type IntIndexable interface {
	IndexInt(idx int64) Node
}
type StrIndexable interface {
	IndexStr(idx string) Node
}
type CtyAble interface {
	CtyType() cty.Type
	AsCtyValue() cty.Value
}
type JQAble interface {
	AsJQValue() (cty.Value, error)
}

type FabricBlock interface {
	Node
	CtyAble
	HCLBlock() *hclsyntax.Block
}
