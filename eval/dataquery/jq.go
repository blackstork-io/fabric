package dataquery

import (
	"context"

	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/ext/customdecode"
	"github.com/itchyny/gojq"
	"github.com/zclconf/go-cty/cty"
	"github.com/zclconf/go-cty/cty/function"

	"github.com/blackstork-io/fabric/pkg/diagnostics"
	"github.com/blackstork-io/fabric/pkg/encapsulator"
	"github.com/blackstork-io/fabric/pkg/utils"
	"github.com/blackstork-io/fabric/plugin"
)

type JqQuery struct {
	query     string
	srcRange  *hcl.Range
	parseOnce func() (*gojq.Code, diagnostics.Diag)
}

var JqQueryType *encapsulator.Codec[JqQuery]

func init() {
	JqQueryType = encapsulator.NewCodec("jq query", &encapsulator.CapsuleOps[JqQuery]{
		ExtensionData: func(key any) any {
			if key != customdecode.CustomExpressionDecoder {
				return nil
			}
			return customdecode.CustomExpressionDecoderFunc(func(expr hcl.Expression, ctx *hcl.EvalContext) (cty.Value, hcl.Diagnostics) {
				val, diags := defineJqQuery(expr, ctx)
				return val, hcl.Diagnostics(diags)
			})
		},
	})
}

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

func defineJqQuery(expr hcl.Expression, evalCtx *hcl.EvalContext) (val cty.Value, diags diagnostics.Diag) {
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

	jq := &JqQuery{
		query:    queryVal.AsString(),
		srcRange: expr.Range().Ptr(),
	}
	jq.parseOnce = utils.OnceVal(jq.parse)

	val = JqQueryType.ToCty(jq)
	return
}

func (q *JqQuery) parse() (code *gojq.Code, diags diagnostics.Diag) {
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

func (q *JqQuery) Eval(ctx context.Context, dataCtx plugin.MapData) (result cty.Value, diags diagnostics.Diag) {
	code, diags := q.parseOnce()
	if diags.HasErrors() {
		return
	}

	res, hasResult := code.RunWithContext(ctx, dataCtx.Any()).Next()
	if !hasResult {
		res = (map[string]any)(nil)
	}
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
	result = plugin.EncapsulatedData.ValToCty(data)
	return
}

func (q *JqQuery) Range() *hcl.Range {
	return q.srcRange
}
