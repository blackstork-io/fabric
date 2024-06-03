package dataquery

import (
	"context"

	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/ext/customdecode"
	"github.com/zclconf/go-cty/cty"

	"github.com/blackstork-io/fabric/pkg/ctyencoder"
	"github.com/blackstork-io/fabric/pkg/diagnostics"
	"github.com/blackstork-io/fabric/pkg/encapsulator"
	"github.com/blackstork-io/fabric/plugin"
)

var DeferredPluginDataType *encapsulator.Codec[DeferredPluginData]

type DeferredPluginData struct {
	val      cty.Value
	srcRange *hcl.Range
}

func (d *DeferredPluginData) Eval(ctx context.Context, dataCtx plugin.MapData) (result cty.Value, diags diagnostics.Diag) {
	var data plugin.Data
	data, diags = ctyencoder.ToPluginData(ctx, dataCtx, d.val)
	if diags.HasErrors() {
		return
	}
	return plugin.EncapsulatedData.ValToCty(data), diags
}

func init() {
	DeferredPluginDataType = encapsulator.NewCodec("arbitrary json-like data", &encapsulator.CapsuleOps[DeferredPluginData]{
		ExtensionData: func(key any) any {
			if key != customdecode.CustomExpressionDecoder {
				return nil
			}
			return customdecode.CustomExpressionDecoderFunc(func(expr hcl.Expression, ctx *hcl.EvalContext) (cty.Value, hcl.Diagnostics) {
				// Allow using jq queries here
				val, diags := expr.Value(JqEvalContext(ctx))
				if diags.HasErrors() {
					return cty.NilVal, diags
				}
				value := DeferredPluginDataType.ToCty(&DeferredPluginData{
					val:      val,
					srcRange: expr.Range().Ptr(),
				})
				return value, diags
			})
		},
	})
}
