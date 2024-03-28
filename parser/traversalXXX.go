package parser

import (
	"fmt"

	"github.com/hashicorp/hcl/v2"
	"github.com/zclconf/go-cty/cty"

	"github.com/blackstork-io/fabric/parser/definitions"
	"github.com/blackstork-io/fabric/pkg/diagnostics"
)

func Resolve[B definitions.FabricBlock](db *DefinedBlocks, expr hcl.Expression) (res B, diags diagnostics.Diag) {
	resAny, diags := db.resolve(expr, res.CtyType())
	if diags.HasErrors() {
		return
	}
	res = resAny.(B) //nolint:forcetypeassert // This type assertion is done via cty in db.resolve
	return
}

func (db *DefinedBlocks) resolve(expr hcl.Expression, expectedType cty.Type) (res any, diags diagnostics.Diag) {
	val, diag := expr.Value(&hcl.EvalContext{
		Variables: db.AsValueMap(),
	})
	if diags.ExtendHcl(diag) {
		return
	}
	ty := val.Type()
	if !ty.Equals(expectedType) {
		diags.Append(&hcl.Diagnostic{
			Severity: hcl.DiagError,
			Summary:  "Incorrect reference",
			Detail: fmt.Sprintf(
				"Expected reference to '%s' got reference to '%s'",
				expectedType.FriendlyName(),
				ty.FriendlyName(),
			),
			Subject: expr.Range().Ptr(),
		})
		return
	}
	res = val.EncapsulatedValue()
	return
}
