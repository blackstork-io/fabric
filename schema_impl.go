package main

import (
	"fmt"

	"github.com/hashicorp/hcl/v2"
)

func updateCommon(b Block, ref Block) {
	// Assigning data from referenced block
	refAttrs := *ref.GetAttrs()
	bAttrs := *b.GetAttrs()
	for k, v := range refAttrs {
		if _, exists := bAttrs[k]; exists {
			continue
		}
		bAttrs[k] = v
	}
	*b.GetType() = *ref.GetType()

	// TODO: ref with meta-blocks: field by field or all together?
	if *b.GetMeta() == nil {
		*b.GetMeta() = *ref.GetMeta()
	}
}

// Implementing Block for ContentBlock

var _ BlockExtra = (*ContentBlockExtra)(nil)

func (br *ContentBlockExtra) GetRef() hcl.Expression {
	return br.Ref
}

func (br *ContentBlockExtra) GetUnparsed() hcl.Body {
	return br.Unparsed
}

var _ Block = (*ContentBlock)(nil)

func (b *ContentBlock) GetAttrs() *hcl.Attributes {
	return &b.Attrs
}

func (b *ContentBlock) GetDecoded() *bool {
	return &b.Decoded
}

func (b *ContentBlock) GetMeta() **MetaBlock {
	return &b.Meta
}

func (b *ContentBlock) GetName() string {
	return b.Name
}

func (b *ContentBlock) GetUnparsed() hcl.Body {
	return b.Unparsed
}

func (b *ContentBlock) GetType() *string {
	return &b.Type
}

func (b *ContentBlock) GetBlockKind() string {
	return BK_CONTENT
}

func (b *ContentBlock) NewBlockExtra() BlockExtra {
	return &ContentBlockExtra{}
}

func (b *ContentBlock) DecodeNestedBlocks(d *Decoder, br BlockExtra) (diag hcl.Diagnostics) {
	extra := br.((*ContentBlockExtra))
	b.NestedContentBlocks = extra.ContentBlocks
	for i := range b.NestedContentBlocks {
		// errors in nested content blocks do not prevent us from parsing the current one
		diag = diag.Extend(d.DecodeBlock(&b.NestedContentBlocks[i]))
	}
	return
}

func (b *ContentBlock) UpdateFromRef(refTgt any, ref hcl.Expression) (diag hcl.Diagnostics) {
	tgt, ok := refTgt.(*ContentBlock)
	if !ok || tgt == nil {
		return diag.Append(invalidRefDiag(refTgt, ref, b.GetBlockKind()))
	}

	if !tgt.Decoded {
		return diag.Append(&hcl.Diagnostic{
			Severity: hcl.DiagError,
			Summary:  "Reference to unparsed data",
			Detail: ("Reference points to contents of the block that hasn't been parsed. " +
				"Make sure that the reference is located after the block is defined and that the block has no errors"),
			Subject: ref.Range().Ptr(),
		})
	}

	updateCommon(b, tgt)

	if b.Query == nil {
		b.Query = tgt.Query
	}
	if b.Title == nil {
		b.Title = tgt.Title
	}

	// TODO: do nested content blocks in the "ref" come over too, or just the attrs?
	return
}

// Implementing Block for DataBlock

var _ BlockExtra = (*DataBlockExtra)(nil)

func (br *DataBlockExtra) GetRef() hcl.Expression {
	return br.Ref
}

func (br *DataBlockExtra) GetUnparsed() hcl.Body {
	return br.Extra
}

var _ Block = (*DataBlock)(nil)

func (b *DataBlock) GetAttrs() *hcl.Attributes {
	return &b.Attrs
}

func (b *DataBlock) GetDecoded() *bool {
	return &b.Decoded
}

func (b *DataBlock) GetMeta() **MetaBlock {
	return &b.Meta
}

func (b *DataBlock) GetName() string {
	return b.Name
}

func (b *DataBlock) GetUnparsed() hcl.Body {
	return b.Extra
}

func (b *DataBlock) GetType() *string {
	return &b.Type
}

func (b *DataBlock) GetBlockKind() string {
	return BK_DATA
}

func (b *DataBlock) NewBlockExtra() BlockExtra {
	return &DataBlockExtra{}
}

func (b *DataBlock) DecodeNestedBlocks(d *Decoder, br BlockExtra) (diag hcl.Diagnostics) {
	// DataBlock doesn't have nested blocks
	return
}

func (b *DataBlock) UpdateFromRef(refTgt any, ref hcl.Expression) (diag hcl.Diagnostics) {
	tgt, ok := refTgt.(*DataBlock)
	if !ok || tgt == nil {
		return diag.Append(invalidRefDiag(refTgt, ref, b.GetBlockKind()))
	}
	if !tgt.Decoded {
		return diag.Append(&hcl.Diagnostic{
			Severity: hcl.DiagError,
			Summary:  "Reference to unparsed data",
			Detail: ("Reference points to contents of the block that hasn't been parsed. " +
				"Make sure that the reference is located after the block is defined and that the block has no errors"),
		})
	}

	updateCommon(b, tgt)
	return
}

func invalidRefDiag(tgt any, ref hcl.Expression, curBlock string) (diag *hcl.Diagnostic) {
	diag = &hcl.Diagnostic{
		Severity:   hcl.DiagError,
		Summary:    "Invalid reference destination",
		Subject:    ref.Range().Ptr(),
		Expression: ref,
	}
	switch tgt.(type) {
	case (*ContentBlock):
		diag.Detail = fmt.Sprintf("%s block can not reference content blocks", CapitalizeFirstLetter(curBlock))
	case (*DataBlock):
		diag.Detail = fmt.Sprintf("%s block can not reference data blocks", CapitalizeFirstLetter(curBlock))
	case (*Document):
		diag.Detail = fmt.Sprintf("%s block can not reference documents", CapitalizeFirstLetter(curBlock))
	case nil:
		diag.Detail = "Unknown error while traversing a reference"
	default:
		diag.Detail = fmt.Sprintf("Reference in %s block points to an unsupported block type", curBlock)
	}
	return
}
