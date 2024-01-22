package definitions

import (
	"sync"

	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/hclsyntax"

	"github.com/blackstork-io/fabric/pkg/diagnostics"
)

// Document and section are very similar conceptually
type DocumentOrSection struct {
	Block *hclsyntax.Block
	Once  sync.Once
	Meta  MetaBlock
}

func (d *DocumentOrSection) IsDocument() bool {
	return d.Block.Type == BlockKindDocument
}

var _ FabricBlock = (*DocumentOrSection)(nil)

func (d *DocumentOrSection) GetHCLBlock() *hcl.Block {
	return d.Block.AsHCLBlock()
}

func (d *DocumentOrSection) Name() string {
	return d.Block.Labels[0]
}

func DefineSectionOrDocument(block *hclsyntax.Block, atTopLevel bool) (doc *DocumentOrSection, diags diagnostics.Diag) {
	nameRequired := atTopLevel || block.Type == BlockKindDocument

	if nameRequired {
		diags.Append(validateBlockName(block, 0, true))
		diags.Append(validateLabelsLength(block, 1, "block_name"))
	} else {
		diags.Append(validateBlockName(block, 0, false))
		diags.Append(validateLabelsLength(block, 1, "<block_name>"))
	}

	if diags.HasErrors() {
		return
	}
	doc = &DocumentOrSection{
		Block: block,
	}
	return
}
