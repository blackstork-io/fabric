package eval

import (
	"context"
	"maps"

	"github.com/hashicorp/hcl/v2"
	"github.com/zclconf/go-cty/cty/convert"

	"github.com/blackstork-io/fabric/eval/dataquery"
	"github.com/blackstork-io/fabric/parser/definitions"
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
		val, diag := dataquery.EvaluateDeferred(ctx, dataCtx, variable.Val)
		diag.DefaultSubject(&variable.ValRange)
		if diags.Extend(diag) {
			vars[variable.Name] = nil
			continue
		}
		v, err := convert.Convert(val, plugin.EncapsulatedData.CtyType())
		if err != nil {
			diags.Append(&hcl.Diagnostic{
				Severity: hcl.DiagError,
				Summary:  "Failed to convert variable value",
				Detail:   err.Error(),
				Subject:  &variable.ValRange,
			})
			continue
		}
		dataVal, err := plugin.EncapsulatedData.FromCty(v)
		if err != nil {
			diags.Append(&hcl.Diagnostic{
				Severity: hcl.DiagError,
				Summary:  "Failed to convert variable value",
				Detail:   err.Error(),
				Subject:  &variable.ValRange,
			})
			continue
		}
		if dataVal == nil {
			vars[variable.Name] = nil
		} else {
			vars[variable.Name] = *dataVal
		}
	}

	return
}
