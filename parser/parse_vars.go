package parser

import (
	"slices"

	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/hclsyntax"

	"github.com/blackstork-io/fabric/eval/dataquery"
	"github.com/blackstork-io/fabric/parser/definitions"
	"github.com/blackstork-io/fabric/parser/evaluation"
	"github.com/blackstork-io/fabric/pkg/diagnostics"
)

func ParseVars(block *hclsyntax.Block) (parsed *definitions.ParsedVars, diags diagnostics.Diag) {
	for _, subBlock := range block.Body.Blocks {
		diags.Append(&hcl.Diagnostic{
			Severity: hcl.DiagWarning,
			Summary:  "Unsupported nesting",
			Detail:   `Vars block does not support nested blocks, did you mean to use nested maps?`,
			Subject:  subBlock.Range().Ptr(),
		})
	}

	evalCtx := dataquery.JqEvalContext(evaluation.EvalContext())
	vars := make([]*definitions.Variable, 0, len(block.Body.Attributes))

	for _, attr := range block.Body.Attributes {
		val, diag := attr.Expr.Value(evalCtx)
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
	// ordered by definition
	slices.SortFunc(vars, func(a, b *definitions.Variable) int {
		return a.NameRange.Start.Byte - b.NameRange.Start.Byte
	})
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
