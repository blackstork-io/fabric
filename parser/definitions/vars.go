package definitions

import (
	"maps"
	"slices"

	"github.com/blackstork-io/fabric/plugin/dataspec"
)

const LocalVarName = "local"

type ParsedVars struct {
	// stored in the order of definition
	Variables []*dataspec.Attr
	ByName    map[string]int
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

// AppendVar append a variable to the parsed vars struct (last in evaluation order).
func (pv *ParsedVars) AppendVar(variable *dataspec.Attr) {
	idx := len(pv.Variables)
	pv.Variables = append(pv.Variables, variable)
	if idx == 0 {
		pv.ByName = make(map[string]int)
	}
	pv.ByName[variable.Name] = idx
}
