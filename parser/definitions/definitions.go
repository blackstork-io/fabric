package definitions

import (
	"github.com/hashicorp/hcl/v2"
	"github.com/zclconf/go-cty/cty"
)

const (
	BlockKindDocument     = "document"
	BlockKindConfig       = "config"
	BlockKindContent      = "content"
	BlockKindPublish      = "publish"
	BlockKindData         = "data"
	BlockKindMeta         = "meta"
	BlockKindVars         = "vars"
	BlockKindSection      = "section"
	BlockKindGlobalConfig = "fabric"

	PluginTypeRef = "ref"
	AttrRefBase   = "base"
	AttrTitle     = "title"
)

type FabricBlock interface {
	GetHCLBlock() *hcl.Block
	CtyType() cty.Type
}

func ToCtyValue(b FabricBlock) cty.Value {
	return cty.CapsuleVal(b.CtyType(), b)
}

// Identifies a plugin block
type Key struct {
	PluginKind string
	PluginName string
	BlockName  string
}
