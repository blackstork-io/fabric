package dataquery

import (
	"context"

	"github.com/hashicorp/hcl/v2"
	"github.com/itchyny/gojq"
	"github.com/zclconf/go-cty/cty"
	"github.com/zclconf/go-cty/cty/convert"
	"github.com/zclconf/go-cty/cty/function"

	"github.com/blackstork-io/fabric/pkg/diagnostics"
	"github.com/blackstork-io/fabric/pkg/encapsulator"
	"github.com/blackstork-io/fabric/pkg/utils"
	"github.com/blackstork-io/fabric/plugin/plugindata"
)

type JqQuery struct {
	query     string
	srcRange  *hcl.Range
	parseOnce func() (*gojq.Code, diagnostics.Diag)
}

var JqQueryType = encapsulator.NewCodec("jq query", &encapsulator.CapsuleOps[JqQuery]{
	CustomExpressionDecoder: func(expr hcl.Expression, evalCtx *hcl.EvalContext) (val *JqQuery, diags diagnostics.Diag) {
		queryVal, diag := expr.Value(evalCtx)
		if diags.Extend(diag) {
			return
		}
		if queryVal.IsNull() || !queryVal.IsKnown() || !queryVal.Type().Equals(cty.String) {
			diags.Append(&hcl.Diagnostic{
				Severity:    hcl.DiagError,
				Summary:     "Invalid argument",
				Detail:      "A string is required",
				Subject:     expr.Range().Ptr(),
				Expression:  expr,
				EvalContext: evalCtx,
			})
			return
		}

		val = &JqQuery{
			query:    queryVal.AsString(),
			srcRange: expr.Range().Ptr(),
		}
		val.parseOnce = utils.OnceVal(val.parse)
		return
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
			Type: function.StaticReturnType(DeferredEvalType.CtyType()),
			Impl: func(args []cty.Value, retType cty.Type) (cty.Value, error) {
				return convert.Convert(args[0], retType)
			},
		}),
	}
	return
}

func (q *JqQuery) parse() (code *gojq.Code, diags diagnostics.Diag) {
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

func (q *JqQuery) DeferredEval(ctx context.Context, dataCtx plugindata.Map) (result cty.Value, diags diagnostics.Diag) {
	var data plugindata.Data
	data, diags = q.Eval(ctx, dataCtx)
	if diags.HasErrors() {
		return
	}
	result = plugindata.EncapsulatedData.ToCty(&data)
	return
}

func (q *JqQuery) Eval(ctx context.Context, dataCtx plugindata.Map) (result plugindata.Data, diags diagnostics.Diag) {
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
	result, err = plugindata.ParseAny(res)
	if err != nil {
		diags.Append(&hcl.Diagnostic{
			Severity: hcl.DiagError,
			Summary:  "Incorrect query result type",
			Detail:   err.Error(),
			Subject:  q.srcRange,
		})
		return
	}
	return
}
