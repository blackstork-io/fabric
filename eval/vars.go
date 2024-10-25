package eval

import (
	"context"
	"maps"

	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/hclsyntax"

	"github.com/blackstork-io/fabric/parser/definitions"
	"github.com/blackstork-io/fabric/pkg/diagnostics"
	"github.com/blackstork-io/fabric/plugin/dataspec"
	"github.com/blackstork-io/fabric/plugin/plugindata"
)

// getVarsCopy creates a new vars key in the data context if it doesn't exist,
// or clones the existing one.
func getVarsCopy(dataCtx plugindata.Map) (vars plugindata.Map) {
	varsData := dataCtx["vars"]
	if varsData == nil {
		vars = plugindata.Map{}
	} else {
		// avoid modifying the original vars
		vars = maps.Clone(varsData.(plugindata.Map))
	}
	dataCtx["vars"] = vars
	return vars
}

// Evaluates `variables` and stores the results in `dataCtx` under the key "vars".
func ApplyVars(ctx context.Context, variables *definitions.ParsedVars, dataCtx plugindata.Map) (diags diagnostics.Diag) {
	if variables.Empty() {
		return
	}
	vars := getVarsCopy(dataCtx)
	var diag diagnostics.Diag
	for _, variable := range variables.Variables {
		vars[variable.Name], diag = evalVar(ctx, dataCtx, variable)
		diags.Extend(diag.Refine(
			diagnostics.DefaultSubject(variable.ValueRange),
		))
	}

	return
}

func evalVar(ctx context.Context, dataCtx plugindata.Map, attr *dataspec.Attr) (data plugindata.Data, diags diagnostics.Diag) {
	val, diags := dataspec.EvalAttr(ctx, attr, dataCtx)
	if diags.HasErrors() {
		return
	}
	dataVal, err := plugindata.Encapsulated.FromCty(val)
	if diags.AppendErr(err, "Failed to convert variable value") {
		return
	}
	if dataVal != nil {
		data = *dataVal
	}
	return
}

func verifyRequiredVars(docDataCtx plugindata.Map, requiredVars []string, block *hclsyntax.Block) (diag diagnostics.Diag) {
	vars, varsPresent := docDataCtx["vars"].(plugindata.Map)
	for _, reqVar := range requiredVars {
		if !varsPresent || vars[reqVar] == nil {
			return diagnostics.FromHcl(&hcl.Diagnostic{
				Severity: hcl.DiagError,
				Summary:  "Missing required variable",
				Detail:   "block requires '" + reqVar + "' var which is not set.",
				Subject:  block.Range().Ptr(),
			})
		}
	}
	return nil
}
