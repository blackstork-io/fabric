package definitions

import (
	"context"
	"maps"
	"slices"

	"github.com/hashicorp/hcl/v2"
	"github.com/zclconf/go-cty/cty"

	"github.com/blackstork-io/fabric/pkg/diagnostics"
	"github.com/blackstork-io/fabric/pkg/encapsulator"
	"github.com/blackstork-io/fabric/plugin"
)

type Query interface {
	Eval(ctx context.Context, dataCtx plugin.MapData) (result plugin.Data, diags diagnostics.Diag)
	Range() *hcl.Range
}

var QueryType = encapsulator.NewDecoder[Query]()

type ParsedVars struct {
	// stored in the order of definition
	Variables []*Variable
	ByName    map[string]int
}

type Variable struct {
	Name      string
	NameRange hcl.Range
	Val       cty.Value
	ValRange  hcl.Range
}

func (pv *ParsedVars) Empty() bool {
	return pv == nil || len(pv.Variables) == 0
}

// MergeWithBaseVars handles merging with vars from ref base.
// Shadowing has different rules, and will be handled at the evaluation stage.
func (pv *ParsedVars) MergeWithBaseVars(baseVars *ParsedVars) *ParsedVars {
	if pv.Empty() {
		return baseVars
	}
	if baseVars.Empty() {
		return pv
	}

	vars := slices.Clone(baseVars.Variables)
	byName := maps.Clone(baseVars.ByName)
	for _, v := range pv.Variables {
		if idx, found := byName[v.Name]; found {
			// redefine, but keep the definition order
			vars[idx] = v
		} else {
			byName[v.Name] = len(vars)
			vars = append(vars, v)
		}
	}
	return &ParsedVars{
		Variables: vars,
		ByName:    byName,
	}
}
