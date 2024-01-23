package definitions

import (
	"github.com/hashicorp/hcl/v2"
)

const (
	BlockKindDocument = "document"
	BlockKindConfig   = "config"
	BlockKindContent  = "content"
	BlockKindData     = "data"
	BlockKindMeta     = "meta"
	BlockKindSection  = "section"

	PluginTypeRef = "ref"
)

type FabricBlock interface {
	GetHCLBlock() *hcl.Block
}

// Identifies a plugin block
type Key struct {
	PluginKind string
	PluginName string
	BlockName  string
}
