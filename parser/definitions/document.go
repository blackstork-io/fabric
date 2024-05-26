package definitions

import (
	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/hclsyntax"
	"github.com/zclconf/go-cty/cty"

	"github.com/blackstork-io/fabric/pkg/diagnostics"
	"github.com/blackstork-io/fabric/pkg/encapsulator"
)

// Document and section are very similar conceptually.
type Document struct {
	Block *hclsyntax.Block
	Name  string
	Meta  *MetaBlock
}

var _ FabricBlock = (*Document)(nil)

func (d *Document) GetHCLBlock() *hcl.Block {
	return d.Block.AsHCLBlock()
}

var ctyDocumentType = encapsulator.New[Document]("document")

func (d *Document) CtyType() cty.Type {
	return ctyDocumentType.Type()
}

func DefineDocument(block *hclsyntax.Block) (doc *Document, diags diagnostics.Diag) {
	diags.Append(validateBlockName(block, 0, true))
	diags.Append(validateLabelsLength(block, 1, "document_name"))
	if diags.HasErrors() {
		return
	}

	if block.Labels[0] == AttrRefBase {
		diags.Append(&hcl.Diagnostic{
			Severity: hcl.DiagError,
			Summary:  "Invalid document declaration",
			Detail:   "Documents can't be refs, only sections can",
			Subject:  &block.LabelRanges[0],
			Context:  block.DefRange().Ptr(),
		})
	}

	doc = &Document{
		Block: block,
		Name:  block.Labels[0],
	}
	return
}
