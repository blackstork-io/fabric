package eval

import (
	"context"
	"fmt"
	"maps"

	"github.com/zclconf/go-cty/cty"

	"github.com/blackstork-io/fabric/parser/definitions"
	"github.com/blackstork-io/fabric/pkg/diagnostics"
	"github.com/blackstork-io/fabric/plugin"
)

func ctyToData(ctx context.Context, dataCtx plugin.MapData, val cty.Value) (result plugin.Data, diags diagnostics.Diag) {
	var diag diagnostics.Diag
	if val.IsNull() {
		return
	}
	ty := val.Type()

	switch {
	case ty.IsPrimitiveType():
		switch ty {
		case cty.String:
			result = plugin.StringData(val.AsString())
			return
		case cty.Number:
			n, _ := val.AsBigFloat().Float64()
			result = plugin.NumberData(n)
			return
		case cty.Bool:
			result = plugin.BoolData(val.True())
			return
		}
	case ty.IsMapType() || ty.IsObjectType():
		res := make(plugin.MapData, val.LengthInt())
		for it := val.ElementIterator(); it.Next(); {
			k, v := it.Element()
			res[k.AsString()], diag = ctyToData(ctx, dataCtx, v)
			diags.Extend(diag)
		}
		result = res
		return
	case ty.IsListType() || ty.IsTupleType():
		res := make(plugin.ListData, 0, val.LengthInt())
		for it := val.ElementIterator(); it.Next(); {
			_, v := it.Element()
			data, diag := ctyToData(ctx, dataCtx, v)
			res = append(res, data)
			diags.Extend(diag)
		}
		result = res
		return
	case definitions.QueryType.ValDecodable(val):
		query := definitions.QueryType.MustFromCty(val)
		return query.Eval(ctx, dataCtx)
	}
	diags.Add("Incorrect data type", fmt.Sprintf("%s is not supported here", ty.FriendlyName()))
	return
}

// Evaluates `variables` and stores the results in `dataCtx` under the key "vars".
func ApplyVars(ctx context.Context, variables *definitions.ParsedVars, dataCtx plugin.MapData) (diags diagnostics.Diag) {
	if variables.Empty() {
		return
	}
	var vars plugin.MapData

	varsData := dataCtx["vars"]
	if varsData == nil {
		vars = plugin.MapData{}
	} else {
		// avoid modifying the original vars
		vars = maps.Clone(varsData.(plugin.MapData))
	}
	dataCtx["vars"] = vars

	for _, variable := range variables.Variables {
		val, diag := ctyToData(ctx, dataCtx, variable.Val)
		diag.DefaultSubject(&variable.ValRange)
		diags.Extend(diag)
		vars[variable.Name] = val
	}

	return
}
