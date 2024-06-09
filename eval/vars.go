package eval

import (
	"context"
	"maps"

	"github.com/blackstork-io/fabric/parser/definitions"
	"github.com/blackstork-io/fabric/pkg/diagnostics"
	"github.com/blackstork-io/fabric/plugin"
	"github.com/blackstork-io/fabric/plugin/plugincty"
)

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
		val, diag := plugin.CustomEvalTransform(ctx, dataCtx, variable.Val)
		diag.DefaultSubject(&variable.ValRange)
		if diags.Extend(diag) {
			vars[variable.Name] = nil
			continue
		}
		dataVal, diag := plugincty.Encode(val)
		diag.DefaultSubject(&variable.ValRange)
		diags.Extend(diag)
		vars[variable.Name] = dataVal
	}

	return
}
