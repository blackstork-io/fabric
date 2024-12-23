package parser

import (
	"context"
	"fmt"

	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/gohcl"
	"github.com/hashicorp/hcl/v2/hclsyntax"

	"github.com/blackstork-io/fabric/parser/definitions"
	"github.com/blackstork-io/fabric/pkg/diagnostics"
)

func (db *DefinedBlocks) ParseDocument(ctx context.Context, d *definitions.Document) (doc *definitions.ParsedDocument, diags diagnostics.Diag) {
	doc = &definitions.ParsedDocument{}
	doc.Source = d

	if title := d.Block.Body.Attributes[definitions.AttrTitle]; title != nil {
		titleContent, diag := db.ParseTitle(ctx, title)
		if !diag.Extend(diags) {
			doc.Content = append(doc.Content, titleContent)
		}
	}

	var origMeta *hcl.Range
	var varsBlock *hclsyntax.Block

	for _, block := range d.Block.Body.Blocks {
		switch block.Type {
		case definitions.BlockKindContent, definitions.BlockKindData, definitions.BlockKindPublish:
			plugin, diag := definitions.DefinePlugin(block, false)
			if diags.Extend(diag) {
				continue
			}
			call, diag := db.ParsePlugin(ctx, plugin)
			if diags.Extend(diag) {
				continue
			}
			switch block.Type {
			case definitions.BlockKindContent:
				doc.Content = append(doc.Content, &definitions.ParsedContent{
					Plugin: call,
				})
			case definitions.BlockKindData:
				doc.Data = append(doc.Data, call)
			case definitions.BlockKindPublish:
				doc.Publish = append(doc.Publish, call)
			default:
				panic("must be exhaustive")
			}
		case definitions.BlockKindVars:
			if varsBlock != nil {
				diags.Append(&hcl.Diagnostic{
					Severity: hcl.DiagError,
					Summary:  "Vars block redefinition",
					Detail: fmt.Sprintf(
						"%s block allows at most one vars block, original vars block was defined at %s:%d",
						d.Block.Type, varsBlock.DefRange().Filename, varsBlock.DefRange().Start.Line,
					),
					Subject: block.DefRange().Ptr(),
					Context: d.Block.Body.Range().Ptr(),
				})
				continue
			}
			varsBlock = block
		case definitions.BlockKindMeta:
			if origMeta != nil {
				diags.Append(&hcl.Diagnostic{
					Severity: hcl.DiagError,
					Summary:  "Meta block redefinition",
					Detail: fmt.Sprintf(
						"%s block allows at most one meta block, original meta block was defined at %s:%d",
						d.Block.Type, origMeta.Filename, origMeta.Start.Line,
					),
					Subject: block.DefRange().Ptr(),
					Context: d.Block.Body.Range().Ptr(),
				})
				continue
			}
			var meta definitions.MetaBlock
			if diags.Extend(gohcl.DecodeBody(block.Body, nil, &meta)) {
				continue
			}
			doc.Meta = &meta
			origMeta = block.DefRange().Ptr()
		case definitions.BlockKindSection:
			section, diag := definitions.DefineSection(block, false)
			if diags.Extend(diag) {
				continue
			}
			parsedSection, diag := db.ParseSection(ctx, section)
			if diags.Extend(diag) {
				continue
			}
			doc.Content = append(doc.Content, &definitions.ParsedContent{
				Section: parsedSection,
			})
		case definitions.BlockKindDynamic:
			dynamic, diag := db.ParseDynamic(ctx, block)
			if diags.Extend(diag) {
				continue
			}
			doc.Content = append(doc.Content, &definitions.ParsedContent{
				Dynamic: dynamic,
			})

		default:
			diags.Append(definitions.NewNestingDiag(
				d.Block.Type,
				block,
				d.Block.Body,
				[]string{
					definitions.BlockKindContent,
					definitions.BlockKindData,
					definitions.BlockKindMeta,
					definitions.BlockKindVars,
					definitions.BlockKindSection,
					definitions.BlockKindPublish,
					definitions.BlockKindDynamic,
				},
			))
			continue
		}
	}

	var diag diagnostics.Diag
	doc.Vars, diag = ParseVars(ctx, varsBlock, d.Block.Body.Attributes[definitions.AttrLocalVar])
	diags.Extend(diag)

	if requiredVarsAttr := d.Block.Body.Attributes[definitions.AttrRequiredVars]; requiredVarsAttr != nil {
		diag := gohcl.DecodeExpression(requiredVarsAttr.Expr, nil, &doc.RequiredVars)
		diags.Extend(diag)
	}
	return
}
