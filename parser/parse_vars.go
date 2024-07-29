package parser

import (
	"context"
	"slices"

	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/hclsyntax"

	"github.com/blackstork-io/fabric/cmd/fabctx"
	"github.com/blackstork-io/fabric/parser/definitions"
	"github.com/blackstork-io/fabric/pkg/diagnostics"
	"github.com/blackstork-io/fabric/plugin/dataspec"
	"github.com/blackstork-io/fabric/plugin/dataspec/deferred"
)

func ParseVars(ctx context.Context, block *hclsyntax.Block, localVar *hclsyntax.Attribute) (parsed *definitions.ParsedVars, diags diagnostics.Diag) {
	if block == nil && localVar == nil {
		parsed = &definitions.ParsedVars{}
		return
	}
	if block != nil && localVar != nil {
		localVarInVars := block.Body.Attributes[definitions.LocalVarName]
		if localVarInVars != nil {
			diags.Append(&hcl.Diagnostic{
				Severity: hcl.DiagError,
				Summary:  "Local var redefinition",
				Detail:   "Local var is defined both in vars block and as a separate argument",
				Subject:  localVar.Range().Ptr(),
				Context:  block.Body.Range().Ptr(),
			})
		} else {
			diags.Append(&hcl.Diagnostic{
				Severity: hcl.DiagWarning,
				Summary:  "Local var specified together with vars block",
				Detail: "It is recommended to use either vars block or local var, not both. " +
					"You can define a variable `local` in the vars block to achieve the same effect.",
				Subject: localVar.Range().Ptr(),
			})
		}
	}
	var varCount int
	if block != nil {
		for _, subBlock := range block.Body.Blocks {
			diags.Append(&hcl.Diagnostic{
				Severity: hcl.DiagWarning,
				Summary:  "Unsupported nesting",
				Detail:   `Vars block does not support nested blocks, did you mean to use nested maps?`,
				Subject:  subBlock.Range().Ptr(),
			})
		}
		varCount = len(block.Body.Attributes)
	}

	if localVar != nil {
		varCount++
	}
	ctx = deferred.WithQueryFuncs(ctx)
	evalCtx := fabctx.GetEvalContext(ctx)
	vars := make([]*definitions.Variable, 0, varCount)

	if block != nil {
		for _, attr := range block.Body.Attributes {
			val, diag := dataspec.DecodeAttr(nil, attr, evalCtx)
			if diags.Extend(diag) {
				continue
			}
			vars = append(vars, &definitions.Variable{
				Name:      attr.Name,
				NameRange: attr.NameRange,
				Val:       val,
				ValRange:  attr.Expr.Range(),
			})
		}
	}
	// ordered by definition
	slices.SortFunc(vars, func(a, b *definitions.Variable) int {
		return a.NameRange.Start.Byte - b.NameRange.Start.Byte
	})
	if localVar != nil {
		// ordered last
		val, diag := dataspec.DecodeAttr(nil, localVar, evalCtx)
		if !diags.Extend(diag) {
			vars = append(vars, &definitions.Variable{
				Name:      definitions.LocalVarName,
				NameRange: localVar.NameRange,
				Val:       val,
				ValRange:  localVar.Expr.Range(),
			})
		}
	}
	byName := make(map[string]int, len(vars))
	for i, v := range vars {
		byName[v.Name] = i
	}
	parsed = &definitions.ParsedVars{
		Variables: vars,
		ByName:    byName,
	}
	return
}
