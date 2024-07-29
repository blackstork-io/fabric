package eval

import (
	"context"
	"maps"

	"github.com/zclconf/go-cty/cty"
	"github.com/zclconf/go-cty/cty/convert"

	"github.com/blackstork-io/fabric/parser/definitions"
	"github.com/blackstork-io/fabric/pkg/diagnostics"
	"github.com/blackstork-io/fabric/plugin/dataspec"
	"github.com/blackstork-io/fabric/plugin/plugindata"
)

// Evaluates `variables` and stores the results in `dataCtx` under the key "vars".
func ApplyVars(ctx context.Context, variables *definitions.ParsedVars, dataCtx plugindata.Map) (diags diagnostics.Diag) {
	if variables.Empty() {
		return
	}
	var vars plugindata.Map

	varsData := dataCtx["vars"]
	if varsData == nil {
		vars = plugindata.Map{}
	} else {
		// avoid modifying the original vars
		vars = maps.Clone(varsData.(plugindata.Map))
	}
	dataCtx["vars"] = vars
	var diag diagnostics.Diag
	for _, variable := range variables.Variables {
		vars[variable.Name], diag = evalVar(ctx, dataCtx, variable.Val)
		diags.Extend(diag.Refine(
			diagnostics.DefaultSubject(variable.ValRange),
		))
	}

	return
}

func evalVar(ctx context.Context, dataCtx plugindata.Map, val cty.Value) (data plugindata.Data, diags diagnostics.Diag) {
	val, diags = dataspec.EvaluateDeferred(ctx, dataCtx, val)
	if diags.HasErrors() {
		return
	}
	v, err := convert.Convert(val, plugindata.Encapsulated.CtyType())
	if diags.AppendErr(err, "Failed to convert variable value") {
		return
	}
	dataVal, err := plugindata.Encapsulated.FromCty(v)
	if diags.AppendErr(err, "Failed to convert variable value") {
		return
	}
	if dataVal != nil {
		data = *dataVal
	}
	return
}
