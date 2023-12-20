package main

import (
	"fmt"

	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/gohcl"
	"github.com/hashicorp/hcl/v2/hclsyntax"
)

func (d *Decoder) Decode() (diag hcl.Diagnostics) {
	for i := range d.root.ContentBlocks {
		diag = diag.Extend(d.DecodeBlock(&d.root.ContentBlocks[i]))
	}
	for i := range d.root.DataBlocks {
		diag = diag.Extend(d.DecodeBlock(&d.root.DataBlocks[i]))
	}
	for i := range d.root.Documents {
		diag = diag.Extend(d.DecodeDocumnet(&d.root.Documents[i]))
	}
	return
}

func (d *Decoder) DecodeDocumnet(doc *Document) (diag hcl.Diagnostics) {
	for i := range doc.ContentBlocks {
		diag = diag.Extend(d.DecodeBlock(&doc.ContentBlocks[i]))
	}
	for i := range doc.DataBlocks {
		diag = diag.Extend(d.DecodeBlock(&doc.DataBlocks[i]))
	}
	return
}

func (d *Decoder) DecodeBlock(block Block) (diag hcl.Diagnostics) {
	extra := block.NewBlockExtra()

	diag = gohcl.DecodeBody(block.GetUnparsed(), nil, extra)

	if diag.HasErrors() {
		return
	}

	// deferring errors in attrs, they do not prevent us from parsing
	deferredDiags := block.DecodeNestedBlocks(d, extra)
	defer func() {
		diag = deferredDiags.Extend(diag)
	}()

	leftover := extra.GetUnparsed()
	attrs, attrDiags := leftover.JustAttributes()
	if attrDiags.HasErrors() {
		// TODO: messy hcl bug workaround, in some cases might silently ignore user's error
		attrDiags = nil
		body := leftover.(*hclsyntax.Body)
		for _, b := range body.Blocks {
			switch b.Type {
			case "meta", "content":
				continue
			default:
				attrDiags = attrDiags.Append(&hcl.Diagnostic{
					Severity: hcl.DiagError,
					Summary:  fmt.Sprintf("Unexpected %q block", b.Type),
					Detail:   "Blocks are not allowed here.",
					Subject:  &b.TypeRange,
				})
			}
		}
	}

	deferredDiags = deferredDiags.Extend(attrDiags)
	*block.GetAttrs() = attrs

	trav, refDiag := traversalForExpr(extra.GetRef())
	if *block.GetType() != "ref" {
		if len(trav) != 0 || len(refDiag) != 0 {
			diag = diag.Append(&hcl.Diagnostic{
				Severity: hcl.DiagWarning,
				Summary:  "Non-empty ref attribute",
				Detail: fmt.Sprintf(
					"Non-empty ref attribute found in block of type '%s'. It will be ignored. Block must have type 'ref' in order to use references",
					*block.GetType(),
				),
				Subject:    extra.GetRef().Range().Ptr(),
				Expression: extra.GetRef(),
			})
		}
		// validate block type
		plugins := d.plugins.ByKind(block.GetBlockKind())
		if _, found := plugins.plugins[*block.GetType()]; !found {
			return diag.Append(&hcl.Diagnostic{
				Severity: hcl.DiagWarning,
				Summary:  fmt.Sprintf("Unknown %s block type", block.GetBlockKind()),
				Detail: fmt.Sprintf(
					"Unknown content block type '%s', valid block types: %s. Referencing or evaluating this block would cause an error",
					*block.GetType(),
					plugins.Names(),
				),
				// TODO: storing type as string doensn't allow good error here. Switch to Expression?
				Subject: block.GetUnparsed().MissingItemRange().Ptr(),
			})
		}
		*block.GetDecoded() = true
		return
	}
	// handling ref
	diag = diag.Extend(refDiag)
	if diag.HasErrors() {
		return
	}
	if len(trav) == 0 {
		missingRef := &hcl.Diagnostic{
			Severity: hcl.DiagError,
		}
		if extra.GetRef().Range().Empty() {
			missingRef.Summary = "Missing ref"
			missingRef.Detail = fmt.Sprintf("Block '%s %s' is of type 'ref', but the ref field is missing", block.GetBlockKind(), block.GetName())
			missingRef.Subject = block.GetUnparsed().MissingItemRange().Ptr()
		} else {
			missingRef.Summary = "Empty ref"
			missingRef.Detail = fmt.Sprintf("Block '%s %s' is of type 'ref', but the ref field is empty", block.GetBlockKind(), block.GetName())
			missingRef.Subject = extra.GetRef().Range().Ptr()
			missingRef.Expression = extra.GetRef()
		}
		return diag.Append(missingRef)
	}

	refTgt, travDiag := d.Traverse(trav)
	// annotate traverse diags
	for _, d := range travDiag {
		if d.Subject == nil {
			d.Subject = extra.GetRef().Range().Ptr()
		} else if d.Context == nil {
			d.Context = extra.GetRef().Range().Ptr()
		}
		if d.Expression == nil {
			d.Expression = extra.GetRef()
		}
	}
	diag = diag.Extend(travDiag)
	if diag.HasErrors() {
		return
	}
	diag = diag.Extend(
		block.UpdateFromRef(refTgt, extra.GetRef()),
	)
	if diag.HasErrors() {
		return
	}

	*block.GetDecoded() = true
	return
}
