package definitions

import (
	"reflect"
	"strings"

	"github.com/hashicorp/hcl/v2"
	"github.com/zclconf/go-cty/cty"
)

const (
	BlockKindDocument     = "document"
	BlockKindConfig       = "config"
	BlockKindContent      = "content"
	BlockKindData         = "data"
	BlockKindMeta         = "meta"
	BlockKindSection      = "section"
	BlockKindGlobalConfig = "fabric"

	PluginTypeRef = "ref"
	AttrRefBase   = "base"
	AttrTitle     = "title"
)

type FabricBlock interface {
	// tree.Node
	GetHCLBlock() *hcl.Block
	CtyType() cty.Type
}

func ToCtyValue(b FabricBlock) cty.Value {
	return cty.CapsuleVal(b.CtyType(), b)
}

func capsuleTypeFor[V any]() cty.Type {
	ty := reflect.TypeOf((*V)(nil)).Elem()
	return cty.Capsule(
		strings.ToLower(ty.Name()),
		ty,
	)
}

// Identifies a plugin block
type Key struct {
	PluginKind string
	PluginName string
	BlockName  string
}
