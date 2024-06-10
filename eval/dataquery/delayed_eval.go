package dataquery

import (
	"context"
	"log/slog"

	"github.com/hashicorp/hcl/v2"
	"github.com/zclconf/go-cty/cty"

	"github.com/blackstork-io/fabric/pkg/diagnostics"
	"github.com/blackstork-io/fabric/pkg/encapsulator"
	"github.com/blackstork-io/fabric/plugin"
	"github.com/blackstork-io/fabric/plugin/plugincty"
)

// DelayedEvalType is  type that transports plugin.Data objects inside of the arguments/configuration of the plugin
// Capable of evaluating queries inside of self, evaluated just before sending the data
// to the plugin (`vars` and `query` have already processed the data)
var DelayedEvalType = encapsulator.NewCodec("data", &encapsulator.CapsuleOps[DelayedEval]{
	CustomExpressionDecoder: func(expr hcl.Expression, ctx *hcl.EvalContext) (value *DelayedEval, diags diagnostics.Diag) {
		// Allow using jq queries here
		val, diag := expr.Value(JqEvalContext(ctx))
		diags = diagnostics.Diag(diag)
		if diags.HasErrors() {
			return
		}
		value = &DelayedEval{
			val:      val,
			srcRange: expr.Range().Ptr(),
		}
		return
	},
	ConversionTo: func(dst cty.Type) func(cty.Value, cty.Path) (*DelayedEval, error) {
		if dst.Equals(plugin.EncapsulatedData.CtyType()) {
			return func(v cty.Value, p cty.Path) (*DelayedEval, error) {
				return &DelayedEval{
					evaluated: true,
					data:      *plugin.EncapsulatedData.MustFromCty(v),
				}, nil
			}
		}
		return nil
	},
	ConversionFrom: func(src cty.Type) func(*DelayedEval, cty.Path) (cty.Value, error) {
		if src.Equals(plugin.EncapsulatedData.CtyType()) {
			return func(de *DelayedEval, p cty.Path) (cty.Value, error) {
				if !de.evaluated {
					return cty.NullVal(plugin.EncapsulatedData.CtyType()), p.NewErrorf("Attempted to encode non-evaluated DelayedEval object")
				}
				return plugin.EncapsulatedData.ToCty(&de.data), nil
			}
		}
		return nil
	},
})

type DelayedEval struct {
	val       cty.Value
	srcRange  *hcl.Range
	evaluated bool
	data      plugin.Data
}

var _ plugin.CustomEval = (*DelayedEval)(nil)

func (d *DelayedEval) Result() plugin.Data {
	if d == nil || !d.evaluated {
		slog.Error("DelayedEval: Result: d.evaluated is false")
		return nil
	}
	return d.data
}

func (d *DelayedEval) CustomEval(ctx context.Context, dataCtx plugin.MapData) (result cty.Value, diags diagnostics.Diag) {
	if d.val == cty.NilVal {
		slog.Error("DelayedEval: CustomEval: d.val is nil")
		return
	}
	var val cty.Value
	val, diags = plugin.CustomEvalTransform(ctx, dataCtx, d.val)
	if diags.HasErrors() {
		return
	}
	data, diag := plugincty.Encode(val)
	if diags.Extend(diag) {
		return
	}
	result = DelayedEvalType.ToCty(&DelayedEval{
		val:       val,
		srcRange:  d.srcRange,
		evaluated: true,
		data:      data,
	})
	return
}
