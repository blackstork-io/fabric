package plugin

import (
	"context"
	"log/slog"

	"github.com/zclconf/go-cty/cty"

	"github.com/blackstork-io/fabric/pkg/diagnostics"
	"github.com/blackstork-io/fabric/pkg/encapsulator"
)

// If encapsulated type implements this interface, it will be evaluated with data context.
type CustomEval interface {
	// Should return a new object, cty values should generally be immutable.
	// The returned value must have the same type or it would be ignored.
	CustomEval(ctx context.Context, dataCtx MapData) (cty.Value, diagnostics.Diag)
}

var CustomEvalType = encapsulator.NewDecoder[CustomEval]()

// Walks the (possibly deeply nested) cty.Value and applies the CustomEval if needed.
func CustomEvalTransform(ctx context.Context, dataCtx MapData, val cty.Value) (res cty.Value, diags diagnostics.Diag) {
	res, _ = cty.Transform(val, func(p cty.Path, v cty.Value) (cty.Value, error) {
		if v.IsNull() || !v.IsKnown() || !CustomEvalType.ValDecodable(v) {
			return v, nil
		}
		v, marks := v.Unmark()
		eval := CustomEvalType.MustFromCty(v)
		newV, diag := eval.CustomEval(ctx, dataCtx)
		if newV.Type().Equals(v.Type()) {
			v = newV
		} else if newV != cty.NilVal {
			slog.Error("Type mismatch in CustomEvalTransform", "original type", v.Type().FriendlyName(), "new type", newV.Type().FriendlyName())
		}
		diag.AddExtra(diagnostics.NewPathExtra(p))
		diags.Extend(diag)
		v = v.WithMarks(marks)
		return v, nil
	})
	return
}
