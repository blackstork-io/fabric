package parser

import (
	"context"
	"fmt"

	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/hclsyntax"

	"github.com/blackstork-io/fabric/parser/definitions"
	"github.com/blackstork-io/fabric/pkg/circularRefDetector"
	"github.com/blackstork-io/fabric/pkg/diagnostics"
	"github.com/blackstork-io/fabric/pkg/utils"
)

func (db *DefinedBlocks) ParseDynamic(ctx context.Context, block *hclsyntax.Block) (parsed *definitions.ParsedDynamic, diags diagnostics.Diag) {
	res := definitions.ParsedDynamic{
		Block: block,
	}

	res.Items, _ = utils.Pop(block.Body.Attributes, definitions.AttrDynamicItems)

	if res.Items == nil {
		diags.Append(&hcl.Diagnostic{
			Severity: hcl.DiagError,
			Summary:  "Dynamic block without items",
			Detail:   fmt.Sprintf("Dynamic block must have an attribute %q", definitions.AttrDynamicItems),
			Subject:  block.DefRange().Ptr(),
		})
	}

	for k, v := range block.Body.Attributes {
		diags.Append(&hcl.Diagnostic{
			Severity: hcl.DiagWarning,
			Summary:  "Unsupported attribute",
			Detail:   fmt.Sprintf("Dynamic block does not support attribute %q, it will be ignored", k),
			Subject:  &v.NameRange,
		})
	}

	validChildren := []string{
		definitions.BlockKindContent,
		definitions.BlockKindSection,
		definitions.BlockKindDynamic,
	}
	validChildrenSet := utils.SliceToSet(validChildren)

	for _, block := range block.Body.Blocks {
		if !utils.Contains(validChildrenSet, block.Type) {
			diags.Append(definitions.NewNestingDiag(
				block.Type,
				block,
				block.Body,
				validChildren,
			))
			continue
		}
		switch block.Type {
		case definitions.BlockKindContent:
			plugin, diag := definitions.DefinePlugin(block, false)
			if diags.Extend(diag) {
				continue
			}
			call, diag := db.ParsePlugin(ctx, plugin)
			if diags.Extend(diag) {
				continue
			}
			res.Content = append(res.Content, &definitions.ParsedContent{
				Plugin: call,
			})
		case definitions.BlockKindSection:
			subSection, diag := definitions.DefineSection(block, false)
			if diags.Extend(diag) {
				continue
			}
			circularRefDetector.Add(block, block.DefRange().Ptr())
			parsedSubSection, diag := db.ParseSection(ctx, subSection)
			circularRefDetector.Remove(block, &diag)
			if diags.Extend(diag) {
				continue
			}
			res.Content = append(res.Content, &definitions.ParsedContent{
				Section: parsedSubSection,
			})
		case definitions.BlockKindDynamic:
			subDynamic, diag := db.ParseDynamic(ctx, block)
			if diags.Extend(diag) {
				continue
			}
			res.Content = append(res.Content, &definitions.ParsedContent{
				Dynamic: subDynamic,
			})
		}
	}
	if len(res.Content) == 0 {
		diags.Append(&hcl.Diagnostic{
			Severity: hcl.DiagWarning,
			Summary:  "Dynamic block without content",
			Detail:   "Dynamic block without any content can be removed, as it has no effect",
			Subject:  block.DefRange().Ptr(),
		})
	}
	parsed = &res
	return
}
