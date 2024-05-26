package eval

import (
	"context"
	"log/slog"
	"maps"

	"github.com/zclconf/go-cty/cty"

	"github.com/blackstork-io/fabric/eval/dataquery"
	"github.com/blackstork-io/fabric/parser/definitions"
	"github.com/blackstork-io/fabric/parser/evaluation"
	"github.com/blackstork-io/fabric/pkg/diagnostics"
	"github.com/blackstork-io/fabric/plugin"
)

func ctyToData(queriesObj *definitions.Queries, val cty.Value, path []any) plugin.Data {
	if val.IsNull() {
		return nil
	}
	ty := val.Type()

	switch {
	case ty.IsPrimitiveType():
		switch ty {
		case cty.String:
			return plugin.StringData(val.AsString())
		case cty.Number:
			n, _ := val.AsBigFloat().Float64()
			return plugin.NumberData(n)
		case cty.Bool:
			return plugin.BoolData(val.True())
		}
	case ty.IsMapType() || ty.IsObjectType():
		result := make(plugin.MapData, val.LengthInt())
		path = append(path, nil)
		for it := val.ElementIterator(); it.Next(); {
			k, v := it.Element()
			key := k.AsString()
			path[len(path)-1] = key
			if id, err := definitions.QueryResultPlaceholderType.FromCty(v); err == nil {
				result[key] = plugin.Data(nil)
				queriesObj.ResultDest(*id, path)
			} else {
				result[key] = ctyToData(queriesObj, v, path)
			}
		}
		return result
	case ty.IsListType() || ty.IsTupleType():
		result := make(plugin.ListData, 0, val.LengthInt())
		path = append(path, nil)
		for it := val.ElementIterator(); it.Next(); {
			_, v := it.Element()
			path[len(path)-1] = len(result)
			if id, err := definitions.QueryResultPlaceholderType.FromCty(v); err == nil {
				result = append(result, plugin.Data(nil))
				queriesObj.ResultDest(*id, path)
			} else {
				result = append(result, ctyToData(queriesObj, v, path))
			}
		}
		return result
	}
	slog.Warn("Unknown type while converting cty to data", "type", ty.FriendlyName())
	return nil
}

// Evaluates `variables` and stores the results in `dataCtx` under the key "vars".
func ApplyVars(ctx context.Context, variables definitions.ParsedVars, dataCtx plugin.MapData) (diags diagnostics.Diag) {
	if len(variables) == 0 {
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

	evalCtx, queriesObj := dataquery.QueryContext(evaluation.EvalContext())
	evalCtx = dataquery.JqEvalContext(evalCtx)

	varsAsObj := make(map[string]cty.Value, len(variables))

	for _, v := range variables {
		val, stdDiag := v.Expr.Value(evalCtx)
		if diags.Extend(stdDiag) {
			continue
		}
		varsAsObj[v.Name] = val
	}

	varsVal := cty.ObjectVal(varsAsObj)
	newVarsMap := ctyToData(queriesObj, varsVal, nil).(plugin.MapData)
	for k, v := range newVarsMap {
		vars[k] = v
	}

	for _, query := range queriesObj.Take() {
		if len(query.ResultPath) == 0 {
			// query is probaly under non-executed condition
			continue
		}
		res, diag := query.Query.Eval(ctx, dataCtx)
		if diags.Extend(diag) {
			continue
		}
		_, diag = plugin.NewDataPath(query.ResultPath).SetRootName(".vars").Set(vars, res)
		diag.DefaultSubject(query.Query.Range())
		diags.Extend(diag)
	}

	return
}
