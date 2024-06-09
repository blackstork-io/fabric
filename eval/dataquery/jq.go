package dataquery

import (
	"context"
	"log/slog"

	"github.com/hashicorp/hcl/v2"
	"github.com/itchyny/gojq"
	"github.com/zclconf/go-cty/cty"
	"github.com/zclconf/go-cty/cty/function"

	"github.com/blackstork-io/fabric/pkg/diagnostics"
	"github.com/blackstork-io/fabric/pkg/encapsulator"
	"github.com/blackstork-io/fabric/pkg/utils"
	"github.com/blackstork-io/fabric/plugin"
)

type JqQuery struct {
	*JqQueryDefinition
	Evaluated bool
	Result    plugin.Data
}

var _ plugin.CustomEval = (*JqQuery)(nil)

type JqQueryDefinition struct {
	query     string
	srcRange  *hcl.Range
	parseOnce func() (*gojq.Code, diagnostics.Diag)
}

var JqQueryType = encapsulator.NewCodec("jq query", &encapsulator.CapsuleOps[JqQuery]{
	CustomExpressionDecoder: func(expr hcl.Expression, evalCtx *hcl.EvalContext) (val *JqQuery, diags diagnostics.Diag) {
		slog.Error("CustomExpressionDecoder called")
		queryVal, diag := expr.Value(evalCtx)
		if diags.Extend(diag) {
			return
		}
		if queryVal.IsNull() || !queryVal.Type().Equals(cty.String) {
			diags.Append(&hcl.Diagnostic{
				Severity:    hcl.DiagError,
				Summary:     "Invalid argument type",
				Detail:      "A string is required",
				Subject:     expr.Range().Ptr(),
				Expression:  expr,
				EvalContext: evalCtx,
			})
			return
		}

		val = &JqQuery{
			JqQueryDefinition: &JqQueryDefinition{
				query:    queryVal.AsString(),
				srcRange: expr.Range().Ptr(),
			},
		}
		val.parseOnce = utils.OnceVal(val.parse)
		return
	},
	ConversionFrom: func(src cty.Type) func(*JqQuery, cty.Path) (cty.Value, error) {
		if src.Equals(plugin.EncapsulatedData.CtyType()) {
			return func(jq *JqQuery, p cty.Path) (cty.Value, error) {
				if !jq.Evaluated {
					return cty.NilVal, p.NewErrorf("Attempted to encode non-evaluated JqQuery")
				}
				return plugin.EncapsulatedData.ValToCty(jq.Result), nil
			}
		}
		return nil
	},
})

const funcName = "query_jq"

// Adds "query_jq" function to the eval context
func JqEvalContext(base *hcl.EvalContext) (evalCtx *hcl.EvalContext) {
	// try finding existing jq eval context in base
	if ctx := utils.EvalContextByFunc(base, funcName); ctx != nil {
		return base
	}

	evalCtx = base.NewChild()
	evalCtx.Functions = map[string]function.Function{
		funcName: function.New(&function.Spec{
			Params: []function.Parameter{
				{
					Name:        "query",
					Description: "The jq query string",
					Type:        JqQueryType.CtyType(),
				},
			},
			Type: function.StaticReturnType(JqQueryType.CtyType()),
			Impl: func(args []cty.Value, retType cty.Type) (cty.Value, error) {
				// The parsing is done by the customdecode on JqType
				return args[0], nil
			},
		}),
	}
	return
}

func (q *JqQueryDefinition) parse() (code *gojq.Code, diags diagnostics.Diag) {
	if q == nil {
		diags.Append(&hcl.Diagnostic{
			Severity: hcl.DiagError,
			Summary:  "Failed to parse the query",
			Detail:   "Query wasn't defined",
		})
		return
	}
	jqQuery, err := gojq.Parse(q.query)
	if err != nil {
		diags.Append(&hcl.Diagnostic{
			Severity: hcl.DiagError,
			Summary:  "Failed to parse the query",
			Detail:   err.Error(),
			Subject:  q.srcRange,
			Extra: diagnostics.GoJQError{
				Err:   err,
				Query: q.query,
			},
		})
		return
	}

	code, err = gojq.Compile(jqQuery)
	if err != nil {
		diags.Append(&hcl.Diagnostic{
			Severity: hcl.DiagError,
			Summary:  "Failed to compile the query",
			Detail:   err.Error(),
			Subject:  q.srcRange,
			Extra: diagnostics.GoJQError{
				Err:   err,
				Query: q.query,
			},
		})
	}
	return
}

func (q *JqQueryDefinition) CustomEval(ctx context.Context, dataCtx plugin.MapData) (result cty.Value, diags diagnostics.Diag) {
	var newQ *JqQuery
	newQ, diags = q.Eval(ctx, dataCtx)
	if diags.HasErrors() {
		return
	}
	result = JqQueryType.ToCty(newQ)
	return
}

func (q *JqQueryDefinition) Eval(ctx context.Context, dataCtx plugin.MapData) (result *JqQuery, diags diagnostics.Diag) {
	if q == nil {
		diags.Append(&hcl.Diagnostic{
			Severity: hcl.DiagError,
			Summary:  "Failed to eval the query",
			Detail:   "Query wasn't defined",
		})
		return
	}
	code, diags := q.parseOnce()
	if diags.HasErrors() {
		return
	}

	res, hasResult := code.RunWithContext(ctx, dataCtx.Any()).Next()
	err, ok := res.(error)
	if ok {
		diags.Append(&hcl.Diagnostic{
			Severity: hcl.DiagError,
			Summary:  "Failed to run the query",
			Detail:   err.Error(),
			Subject:  q.srcRange,
			Extra: diagnostics.GoJQError{
				Err:   err,
				Query: q.query,
			},
		})
		return
	}
	if !hasResult {
		res = nil
	}
	data, err := plugin.ParseDataAny(res)
	if err != nil {
		diags.Append(&hcl.Diagnostic{
			Severity: hcl.DiagError,
			Summary:  "Incorrect query result type",
			Detail:   err.Error(),
			Subject:  q.srcRange,
		})
		return
	}
	result = &JqQuery{
		JqQueryDefinition: q,
		Result:            data,
		Evaluated:         true,
	}
	return
}
