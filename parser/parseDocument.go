package parser

import (
	"fmt"

	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/gohcl"

	"github.com/blackstork-io/fabric/parser/definitions"
	"github.com/blackstork-io/fabric/pkg/diagnostics"
)

func (db *DefinedBlocks) ParseDocument(d *definitions.Document) (doc *definitions.ParsedDocument, diags diagnostics.Diag) {
	doc = &definitions.ParsedDocument{}
	if title := d.Block.Body.Attributes["title"]; title != nil {
		doc.Content = append(doc.Content, definitions.NewTitle(title, db.DefaultConfig))
	}

	var origMeta *hcl.Range

	for _, block := range d.Block.Body.Blocks {
		switch block.Type {
		case definitions.BlockKindContent, definitions.BlockKindData:
			plugin, diag := definitions.DefinePlugin(block, false)
			if diags.Extend(diag) {
				continue
			}
			call, diag := db.ParsePlugin(plugin)
			if diags.Extend(diag) {
				continue
			}
			switch block.Type {
			case definitions.BlockKindContent:
				doc.Content = append(doc.Content, (*definitions.ParsedContent)(call))
			case definitions.BlockKindData:
				doc.Data = append(doc.Data, (*definitions.ParsedData)(call))
			default:
				panic("must be exhaustive")
			}

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
			if diags.ExtendHcl(gohcl.DecodeBody(block.Body, nil, &meta)) {
				continue
			}
			doc.Meta = &meta
			origMeta = block.DefRange().Ptr()
		case definitions.BlockKindSection:
			section, diag := definitions.DefineSection(block, false)
			if diags.Extend(diag) {
				continue
			}
			parsedSection, diag := db.ParseSection(section)
			if diags.Extend(diag) {
				continue
			}
			doc.Content = append(doc.Content, parsedSection)
		default:
			diags.Append(definitions.NewNestingDiag(
				d.Block.Type,
				block,
				d.Block.Body,
				[]string{
					definitions.BlockKindContent,
					definitions.BlockKindData,
					definitions.BlockKindMeta,
					definitions.BlockKindSection,
				},
			))
			continue
		}
	}
	return
}
