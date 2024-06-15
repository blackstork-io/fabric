package parser

import (
	"context"
	"fmt"

	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/gohcl"
	"github.com/hashicorp/hcl/v2/hclsyntax"

	"github.com/blackstork-io/fabric/parser/definitions"
	"github.com/blackstork-io/fabric/pkg/circularRefDetector"
	"github.com/blackstork-io/fabric/pkg/diagnostics"
	"github.com/blackstork-io/fabric/pkg/utils"
)

// Evaluates a defined plugin.
func (db *DefinedBlocks) ParseSection(ctx context.Context, section *definitions.Section) (res *definitions.ParsedSection, diags diagnostics.Diag) {
	if circularRefDetector.Check(section) {
		// This produces a bit of an incorrect error and shouldn't trigger in normal operation
		// but I re-check for the circular refs here out of abundance of caution:
		// deadlocks or infinite loops may occur, and are hard to debug
		diags.Append(&hcl.Diagnostic{
			Severity: hcl.DiagError,
			Summary:  "Circular reference detected",
			Detail:   "Looped back to this block through reference chain:",
			Subject:  section.Block.DefRange().Ptr(),
			Extra:    diagnostics.NewTracebackExtra(),
		})
		return
	}
	section.Once.Do(func() {
		res, diags = db.parseSection(ctx, section)
		if diags.HasErrors() {
			return
		}
		section.ParseResult = res
		section.Parsed = true
	})
	if !section.Parsed {
		if diags == nil {
			diags.Append(diagnostics.RepeatedError)
		}
		return
	}
	res = section.ParseResult
	return
}

func (db *DefinedBlocks) parseSection(ctx context.Context, section *definitions.Section) (parsed *definitions.ParsedSection, diags diagnostics.Diag) {
	res := definitions.ParsedSection{}
	if title := section.Block.Body.Attributes["title"]; title != nil {
		res.Title = definitions.NewTitle(title, db.DefaultConfig)
	}

	var origMeta *hcl.Range
	var varsBlock *hclsyntax.Block
	var refBase hclsyntax.Expression

	var validChildren []string
	if section.IsRef() {
		base := section.Block.Body.Attributes["base"]
		if base == nil {
			diags.Append(&hcl.Diagnostic{
				Severity: hcl.DiagError,
				Summary:  "Missing 'base' attribute",
				Detail:   "Ref blocks must contain a 'base' attribute",
				Subject:  section.Block.Body.MissingItemRange().Ptr(),
			})
			return
		}
		refBase = base.Expr
		validChildren = []string{
			definitions.BlockKindMeta,
			definitions.BlockKindVars,
		}
	} else {
		validChildren = []string{
			definitions.BlockKindContent,
			definitions.BlockKindMeta,
			definitions.BlockKindSection,
			definitions.BlockKindVars,
		}
	}
	validChildrenSet := utils.SliceToSet(validChildren)

	for _, block := range section.Block.Body.Blocks {
		if !utils.Contains(validChildrenSet, block.Type) {
			diags.Append(definitions.NewNestingDiag(
				section.Block.Type,
				block,
				section.Block.Body,
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
		case definitions.BlockKindMeta:
			if origMeta != nil {
				diags.Append(&hcl.Diagnostic{
					Severity: hcl.DiagError,
					Summary:  "Meta block redefinition",
					Detail: fmt.Sprintf(
						"%s block allows at most one meta block, original meta block was defined at %s:%d",
						section.Block.Type, origMeta.Filename, origMeta.Start.Line,
					),
					Subject: block.DefRange().Ptr(),
					Context: section.Block.Body.Range().Ptr(),
				})
				continue
			}
			var meta definitions.MetaBlock
			if diags.Extend(gohcl.DecodeBody(block.Body, nil, &meta)) {
				continue
			}
			res.Meta = &meta
			origMeta = block.DefRange().Ptr()
		case definitions.BlockKindVars:
			if varsBlock != nil {
				diags.Append(&hcl.Diagnostic{
					Severity: hcl.DiagError,
					Summary:  "Vars block redefinition",
					Detail: fmt.Sprintf(
						"%s block allows at most one vars block, original vars block was defined at %s:%d",
						section.Block.Type, varsBlock.DefRange().Filename, varsBlock.DefRange().Start.Line,
					),
					Subject: block.DefRange().Ptr(),
					Context: section.Block.Body.Range().Ptr(),
				})
				continue
			}
			varsBlock = block

		case definitions.BlockKindSection:
			subSection, diag := definitions.DefineSection(block, false)
			if diags.Extend(diag) {
				continue
			}
			circularRefDetector.Add(section, block.DefRange().Ptr())
			parsedSubSection, diag := db.ParseSection(ctx, subSection)
			circularRefDetector.Remove(section, &diag)
			if diags.Extend(diag) {
				continue
			}
			res.Content = append(res.Content, &definitions.ParsedContent{
				Section: parsedSubSection,
			})
		}
	}

	var diag diagnostics.Diag
	res.Vars, diag = ParseVars(
		ctx,
		varsBlock,
		section.Block.Body.Attributes[definitions.AttrLocalVar],
	)
	diags.Extend(diag)

	if refBase == nil {
		parsed = &res
		return
	}
	// Parse ref
	baseSection, diag := Resolve[*definitions.Section](db, refBase)
	if diags.Extend(diag) {
		return
	}
	circularRefDetector.Add(section, refBase.Range().Ptr())
	defer circularRefDetector.Remove(section, &diags)
	if circularRefDetector.Check(baseSection) {
		diags.Append(&hcl.Diagnostic{
			Severity: hcl.DiagError,
			Summary:  "Circular reference detected",
			Detail:   "Looped back to this block through reference chain:",
			Subject:  section.Block.DefRange().Ptr(),
			Extra:    diagnostics.NewTracebackExtra(),
		})
		return
	}
	baseEval, diag := db.ParseSection(ctx, baseSection)
	if diags.Extend(diag) {
		return
	}

	// update from base:
	if res.Title == nil {
		res.Title = baseEval.Title
	}
	if res.Meta == nil {
		res.Meta = baseEval.Meta
	}
	res.Vars = res.Vars.MergeWithBaseVars(baseEval.Vars)
	res.Content = append(res.Content, baseEval.Content...)

	parsed = &res
	return
}
