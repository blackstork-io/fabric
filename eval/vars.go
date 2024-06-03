package eval

import (
	"context"
	"maps"

	"github.com/blackstork-io/fabric/parser/definitions"
	"github.com/blackstork-io/fabric/pkg/ctyencoder"
	"github.com/blackstork-io/fabric/pkg/diagnostics"
	"github.com/blackstork-io/fabric/plugin"
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
		val, diag := ctyencoder.ToPluginData(ctx, dataCtx, variable.Val)
		diag.DefaultSubject(&variable.ValRange)
		diags.Extend(diag)
		vars[variable.Name] = val
	}

	return
}
