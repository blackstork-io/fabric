package main

import (
	"fmt"
	"slices"
	"sync"

	"github.com/hashicorp/hcl/v2"
)

func traversalForExpr(expr hcl.Expression) (trav hcl.Traversal, diag hcl.Diagnostics) {
	// ignore diags, just checking if the val is null
	val, _ := expr.Value(nil)
	if val.IsNull() {
		// empty ref
		return
	}
	trav, diag = hcl.AbsTraversalForExpr(expr)
	if diag.HasErrors() {
		trav = nil
	}
	return
}

func (d *Decoder) Traverse(trav hcl.Traversal) (tgt any, diag hcl.Diagnostics) {
	travPos := 0
	blockKind, bkTrav, bkDiag := decodeHclTraverser(trav, travPos, "block kind")
	travPos++
	diag = diag.Extend(bkDiag)
	if diag.HasErrors() {
		return
	}

	switch blockKind {
	case ContentBlockName:
		return d.TraverseContentBlocks(d.root.ContentBlocks, trav, travPos)
	case DataBlockName:
		return d.TraverseDataBlocks(d.root.DataBlocks, trav, travPos)
	case DocumentBlockName:
		return d.TraverseDocuments(d.root.Documents, trav, travPos)
	default:
		return nil, diag.Append(&hcl.Diagnostic{
			Severity: hcl.DiagError,
			Summary:  "Unknown block kind",
			Detail:   fmt.Sprintf("Unknown block kind '%s', valid kinds: %s", blockKind, validBlockKinds()),
			Subject:  bkTrav.SourceRange().Ptr(),
		})
	}
}

// TODO: TraverseContentBlocks and TraverseDataBlocks are identical, except for the type of the block
// Try to join the code by using interfaces at a later date.
func (d *Decoder) TraverseContentBlocks(cb []ContentBlock, trav hcl.Traversal, travPos int) (tgt any, diag hcl.Diagnostics) { //nolint: dupl
	blockType, blockTypeTrav, btDiag := decodeHclTraverser(trav, travPos, "content block type")
	travPos++
	diag = diag.Extend(btDiag)
	blockName, blockNameTrav, bnDiag := decodeHclTraverser(trav, travPos, "content block name")
	travPos++
	diag = diag.Extend(bnDiag)
	if diag.HasErrors() {
		return
	}
	// find referenced block
	n := slices.IndexFunc(cb, func(cb ContentBlock) bool {
		return cb.Type == blockType && cb.Name == blockName
	})
	if n == -1 {
		subj := blockTypeTrav.SourceRange()
		subj.End = blockNameTrav.SourceRange().End
		diag = diag.Append(&hcl.Diagnostic{
			Severity: hcl.DiagError,
			Summary:  "Content block not found",
			Detail: fmt.Sprintf(
				"Content block with type '%s' and name '%s' was not found",
				blockType,
				blockName,
			),
			Subject: &subj,
		})
		return
	}
	return d.TraverseContentBlock(&cb[n], trav, travPos)
}

func (d *Decoder) TraverseContentBlock(cb *ContentBlock, trav hcl.Traversal, travPos int) (tgt any, diag hcl.Diagnostics) {
	if travPos == len(trav) {
		// we've traversed to the destination block!
		return cb, nil
	}
	blockKind, bkTrav, bkDiag := decodeHclTraverser(trav, travPos, "block kind")
	travPos++
	diag = diag.Extend(bkDiag)
	if diag.HasErrors() {
		return
	}

	switch blockKind {
	case ContentBlockName:
		if !cb.Decoded {
			subj := trav[0].SourceRange()
			subj.End = trav[travPos-1].SourceRange().End
			return nil, diag.Append(&hcl.Diagnostic{
				Severity: hcl.DiagError,
				Summary:  "Reference to unparsed data",
				Detail: ("Reference passed through the contents of the block that hasn't been parsed. " +
					"Make sure that the reference is located after the block is defined and that the block is correct"),
				Subject: &subj,
			})
		}
		return d.TraverseContentBlocks(cb.NestedContentBlocks, trav, travPos)
	case DataBlockName, DocumentBlockName:
		return nil, diag.Append(&hcl.Diagnostic{
			Severity: hcl.DiagError,
			Summary:  "Invalid block kind",
			Detail:   fmt.Sprintf("Content blocks can contain only 'content' subblocks, '%s' is invalid", blockKind),
			Subject:  bkTrav.SourceRange().Ptr(),
		})
	default:
		return nil, diag.Append(&hcl.Diagnostic{
			Severity: hcl.DiagError,
			Summary:  "Unknown block kind",
			Detail:   fmt.Sprintf("Unknown kind '%s', content blocks can contain only 'content' subblocks", blockKind),
			Subject:  bkTrav.SourceRange().Ptr(),
		})
	}
}

func (d *Decoder) TraverseDataBlocks(db []DataBlock, trav hcl.Traversal, travPos int) (tgt any, diag hcl.Diagnostics) { //nolint: dupl
	blockType, blockTypeTrav, btDiag := decodeHclTraverser(trav, travPos, "data block type")
	travPos++
	diag = diag.Extend(btDiag)
	blockName, blockNameTrav, bnDiag := decodeHclTraverser(trav, travPos, "data block name")
	travPos++
	diag = diag.Extend(bnDiag)
	if diag.HasErrors() {
		return
	}
	// find referenced block
	n := slices.IndexFunc(db, func(db DataBlock) bool {
		return db.Type == blockType && db.Name == blockName
	})
	if n == -1 {
		subj := blockTypeTrav.SourceRange()
		subj.End = blockNameTrav.SourceRange().End
		diag = diag.Append(&hcl.Diagnostic{
			Severity: hcl.DiagError,
			Summary:  "Content block not found",
			Detail: fmt.Sprintf(
				"Content block with type '%s' and name '%s' was not found",
				blockType,
				blockName,
			),
			Subject: &subj,
		})
		return
	}
	return d.TraverseDataBlock(&db[n], trav, travPos)
}

func (d *Decoder) TraverseDataBlock(db *DataBlock, trav hcl.Traversal, travPos int) (tgt any, diag hcl.Diagnostics) {
	if travPos == len(trav) {
		// we've traversed to the destination block!
		return db, nil
	}
	blockKind, bkTrav, bkDiag := decodeHclTraverser(trav, travPos, "block kind")
	travPos++

	switch blockKind {
	case ContentBlockName, DataBlockName, DocumentBlockName:
		diag = diag.Extend(bkDiag)

		subj := trav[0].SourceRange()
		subj.End = trav[travPos-1].SourceRange().End
		diag = diag.Append(&hcl.Diagnostic{
			Severity: hcl.DiagError,
			Summary:  "Invalid block nesting",
			Detail:   "Data blocks can not contain any subblocks",
			Subject:  &subj,
		})
	default:
		diag = diag.Append(&hcl.Diagnostic{
			Severity: hcl.DiagError,
			Summary:  "Unknown block kind",
			Detail:   fmt.Sprintf("Unknown block kind '%s'", blockKind),
			Subject:  bkTrav.SourceRange().Ptr(),
		})
	}
	return nil, diag
}

func (d *Decoder) TraverseDocuments(docs []Document, trav hcl.Traversal, travPos int) (tgt any, diag hcl.Diagnostics) {
	docName, docNameTrav, btDiag := decodeHclTraverser(trav, travPos, "document name")
	travPos++
	diag = diag.Extend(btDiag)
	if diag.HasErrors() {
		return
	}
	// find referenced block
	n := slices.IndexFunc(docs, func(doc Document) bool {
		return doc.Name == docName
	})

	if n == -1 {
		diag = diag.Append(&hcl.Diagnostic{
			Severity: hcl.DiagError,
			Summary:  "Document not found",
			Detail: fmt.Sprintf(
				"Document with name '%s' was not found",
				docName,
			),
			Subject: docNameTrav.SourceRange().Ptr(),
		})
		return
	}
	return d.TraverseDocument(&docs[n], trav, travPos)
}

func (d *Decoder) TraverseDocument(doc *Document, trav hcl.Traversal, travPos int) (tgt any, diag hcl.Diagnostics) {
	if travPos == len(trav) {
		// we've traversed to the destination document!
		// (currently it's an invalid target for a ref, but this validation is on caller to do)
		return doc, nil
	}
	blockKind, bkTrav, bkDiag := decodeHclTraverser(trav, travPos, "block kind")
	travPos++
	diag = diag.Extend(bkDiag)
	if diag.HasErrors() {
		return
	}

	switch blockKind {
	case ContentBlockName:
		return d.TraverseContentBlocks(doc.ContentBlocks, trav, travPos)
	case DataBlockName:
		return d.TraverseDataBlocks(doc.DataBlocks, trav, travPos)
	case DocumentBlockName:
		return nil, diag.Append(&hcl.Diagnostic{
			Severity: hcl.DiagError,
			Summary:  "Invalid block kind",
			Detail:   fmt.Sprintf("Documents can contain only 'content' and 'data' subblocks, not '%s'", blockKind),
			Subject:  bkTrav.SourceRange().Ptr(),
		})
	default:
		return nil, diag.Append(&hcl.Diagnostic{
			Severity: hcl.DiagError,
			Summary:  "Unknown block kind",
			Detail:   fmt.Sprintf("Unknown block kind '%s', documents can contain only 'content' and 'data' subblocks", blockKind),
			Subject:  bkTrav.SourceRange().Ptr(),
		})
	}
}

// Utils

var validBlockKinds = sync.OnceValue(func() string {
	return JoinSurround(", ", "'", ContentBlockName, DataBlockName, DocumentBlockName)
})

func decodeHclTraverser(trav hcl.Traversal, travPos int, what string) (name string, traverser hcl.Traverser, diag hcl.Diagnostics) {
	if travPos >= len(trav) {
		diag = diag.Append(&hcl.Diagnostic{
			Severity: hcl.DiagError,
			Summary:  fmt.Sprintf("Missing %s", what),
			Detail:   fmt.Sprintf("Required %s path element wasn't specified", what),
			Subject:  trav[len(trav)-1].SourceRange().Ptr(),
		})
		return
	}
	traverser = trav[travPos]

	switch typedTraverser := traverser.(type) {
	case hcl.TraverseRoot:
		name = typedTraverser.Name
	case hcl.TraverseAttr:
		name = typedTraverser.Name
	default:
		diag = diag.Append(&hcl.Diagnostic{
			Severity: hcl.DiagError,
			Summary:  "Invalid reference",
			Detail:   "The ref attribute can not contain this operation",
			Subject:  typedTraverser.SourceRange().Ptr(),
		})
	}
	return
}
