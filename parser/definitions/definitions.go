package definitions

import (
	"github.com/hashicorp/hcl/v2/hclsyntax"
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
	BlockKindDynamic      = "dynamic"
	BlockKindFormat       = "format"

	PluginTypeRef    = "ref"
	AttrRefBase      = "base"
	AttrTitle        = "title"
	AttrDependsOn    = "depends_on"
	AttrLocalVar     = "local_var"
	AttrRequiredVars = "required_vars"
	AttrIsIncluded   = "is_included"
	AttrDynamicItems = "items"
)

type FabricBlock interface {
	GetHCLBlock() *hclsyntax.Block
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
