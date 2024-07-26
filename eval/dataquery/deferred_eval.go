package dataquery

import (
	"context"
	"reflect"

	"github.com/zclconf/go-cty/cty"
	"github.com/zclconf/go-cty/cty/convert"

	"github.com/blackstork-io/fabric/pkg/diagnostics"
	"github.com/blackstork-io/fabric/pkg/encapsulator"
	"github.com/blackstork-io/fabric/plugin/dataspec"
	"github.com/blackstork-io/fabric/plugin/plugindata"
)

type DeferredEvaluatable interface {
	DeferredEval(ctx context.Context, dataCtx plugindata.Map) (result cty.Value, diags diagnostics.Diag)
}

var deferredEvalReflectType = reflect.TypeFor[DeferredEvaluatable]()

var DeferredEvalType *encapsulator.Codec[deferredEval] = encapsulator.NewCodec("Deferred Evaluation", &encapsulator.CapsuleOps[deferredEval]{
	ConversionTo: func(dst cty.Type) func(cty.Value, cty.Path) (*deferredEval, error) {
		if !reflect.PointerTo(dst.EncapsulatedType()).Implements(deferredEvalReflectType) {
			return nil
		}
		return func(v cty.Value, p cty.Path) (*deferredEval, error) {
			if v.IsNull() {
				return nil, p.NewErrorf("can't convert null to DeferredEval")
			}
			if !v.IsKnown() {
				return nil, p.NewErrorf("can't convert unknown value to DeferredEval")
			}
			if v.EncapsulatedValue() == nil {
				return nil, p.NewErrorf("can't convert nil value to DeferredEval")
			}
			return &deferredEval{
				res: v,
			}, nil
		}
	},
	ConversionFrom: func(src cty.Type) func(*deferredEval, cty.Path) (cty.Value, error) {
		return func(de *deferredEval, p cty.Path) (cty.Value, error) {
			switch {
			case de.status == 0:
				return cty.NilVal, p.NewErrorf("DeferredEval has not been evaluated")
			case de.status < 0:
				return cty.NullVal(src), nil
			default:
				res, err := convert.Convert(de.res, src)
				err = p.NewError(err)
				return res, err
			}
		}
	},
})

type deferredEval struct {
	res    cty.Value
	status int
}

func (d *deferredEval) EvalWrappedVal(ctx context.Context, dataCtx plugindata.Map) (result cty.Value, diags diagnostics.Diag) {
	res := d
	switch {
	case d.status > 0:
	case d.status < 0:
		diags.Append(diagnostics.RepeatedError)
	default:
		res = &deferredEval{}
		res.res, diags = d.res.EncapsulatedValue().(DeferredEvaluatable).DeferredEval(ctx, dataCtx)
		if diags.HasErrors() {
			res.status = -1
		} else {
			res.status = +1
		}
	}
	result = DeferredEvalType.ToCty(res)
	return
}

// Walks the (possibly deeply nested) cty.Value and applies the CustomEval if needed.
func EvaluateDeferred(ctx context.Context, dataCtx plugindata.Map, val cty.Value) (res cty.Value, diags diagnostics.Diag) {
	res, _ = cty.Transform(val, func(p cty.Path, v cty.Value) (cty.Value, error) {
		if v.IsNull() || !v.IsKnown() || !DeferredEvalType.ValCtyTypeEqual(v) {
			return v, nil
		}
		v, marks := v.Unmark()
		eval := DeferredEvalType.MustFromCty(v)
		v, diag := eval.EvalWrappedVal(ctx, dataCtx)
		diag.AddExtra(diagnostics.NewPathExtra(p))
		diags.Extend(diag)
		v = v.WithMarks(marks)
		return v, nil
	})
	return
}

func EvaluateDeferredBlock(ctx context.Context, dataCtx plugindata.Map, block *dataspec.Block) (diags diagnostics.Diag) {
	if block == nil {
		return
	}
	var diag diagnostics.Diag
	for _, val := range block.Attrs {
		val.Value, diag = EvaluateDeferred(ctx, dataCtx, val.Value)
		diags.Extend(diag)
	}
	for _, block := range block.Blocks {
		diags.Extend(EvaluateDeferredBlock(ctx, dataCtx, block))
	}
	return
}
