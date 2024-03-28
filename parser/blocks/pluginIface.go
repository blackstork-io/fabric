package blocks

import "github.com/blackstork-io/fabric/parser/blocks/internal/tree"

type Plugin interface {
	Kind() string
	// Plugin name, as given in definition (may be "ref")
	DefName() string
	// Plugin name, resolved ("ref" replaced by the appropriate plugin name)
	Name() string
	BlockName() string

	commonPlugin() *plugin

	tree.CtyAble
	tree.Node
	FabricBlock
}

// func BlockName() string
// func DefName() string
// func DefRange() hcl.Range
// func IsRef() bool
// func Kind() string
// func Parse(ctx interface{GetEvalContext() *hcl.EvalContext}) (diags diagnostics.Diag)
