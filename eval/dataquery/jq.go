package dataquery

import (
	"context"

	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/ext/customdecode"
	"github.com/itchyny/gojq"
	"github.com/zclconf/go-cty/cty"
	"github.com/zclconf/go-cty/cty/function"

	"github.com/blackstork-io/fabric/parser/definitions"
	"github.com/blackstork-io/fabric/pkg/diagnostics"
	"github.com/blackstork-io/fabric/pkg/encapsulator"
	"github.com/blackstork-io/fabric/pkg/utils"
	"github.com/blackstork-io/fabric/plugin"
)

func compileJq(expr hcl.Expression, evalCtx *hcl.EvalContext) (val cty.Value, diags diagnostics.Diag) {
	queryCtx := utils.EvalContextByVar(evalCtx, definitions.QueryKey)
	if queryCtx == nil {
		panic("query_jq must be called within QueryContext")
	}
	queries := definitions.QueriesType.MustFromCty(queryCtx.Variables[definitions.QueryKey])
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
	query := queryVal.AsString()
	jqQuery, err := gojq.Parse(query)
	if err != nil {
		diags.Append(&hcl.Diagnostic{
			Severity: hcl.DiagError,
			Summary:  "Failed to parse the query",
			Detail:   err.Error(),
			Subject:  expr.Range().Ptr(),
			Extra: diagnostics.GoJQError{
				Err:   err,
				Query: query,
			},
		})
		return
	}

	code, err := gojq.Compile(jqQuery)
	if err != nil {
		diags.Append(&hcl.Diagnostic{
			Severity: hcl.DiagError,
			Summary:  "Failed to compile the query",
			Detail:   err.Error(),
			Subject:  expr.Range().Ptr(),
			Extra: diagnostics.GoJQError{
				Err:   err,
				Query: query,
			},
		})
		return
	}
	idx := queries.Append(&JqQuery{
		query:    query,
		code:     code,
		srcRange: expr.Range().Ptr(),
	})
	val = JqQueryType.ToCty(&idx)
	return
}

var JqQueryType *encapsulator.Encapsulator[int]

func init() {
	JqQueryType = encapsulator.NewWithOps[int]("jq query", &cty.CapsuleOps{
		ExtensionData: func(key interface{}) interface{} {
			switch key {
			case customdecode.CustomExpressionDecoder:
				return customdecode.CustomExpressionDecoderFunc(
					func(expr hcl.Expression, ctx *hcl.EvalContext) (cty.Value, hcl.Diagnostics) {
						val, diags := compileJq(expr, ctx)
						return val, hcl.Diagnostics(diags)
					},
				)
			default:
				return nil
			}
		},
	})
}

type JqQuery struct {
	query    string
	code     *gojq.Code
	srcRange *hcl.Range
}

func (q *JqQuery) Range() *hcl.Range {
	return q.srcRange
}

func (q *JqQuery) Eval(ctx context.Context, dataCtx plugin.MapData) (result plugin.Data, diags diagnostics.Diag) {
	res, hasResult := q.code.RunWithContext(ctx, dataCtx.Any()).Next()
	if !hasResult {
		return nil, nil
	}
	err, ok := res.(error)
	if ok {
		return nil, diagnostics.Diag{{
			Severity: hcl.DiagError,
			Summary:  "Failed to run the query",
			Detail:   err.Error(),
			Subject:  q.srcRange,
			Extra: diagnostics.GoJQError{
				Err:   err,
				Query: q.query,
			},
		}}
	}
	result, err = plugin.ParseDataAny(res)
	if err != nil {
		return nil, diagnostics.Diag{{
			Severity: hcl.DiagError,
			Summary:  "Incorrect query result type",
			Detail:   err.Error(),
			Subject:  q.srcRange,
		}}
	}
	return
}

const funcName = "query_jq"

// Adds "query_jq" function to the eval context
func JqEvalContext(base *hcl.EvalContext) (evalCtx *hcl.EvalContext) {
	// try finding exisitng jq eval context in base
	if ctx := utils.EvalContextByFunc(base, funcName); ctx != nil {
		return base
	}
	// Check for query collection, if found - store query to eval there
	if ctx := utils.EvalContextByVar(base, definitions.QueryKey); ctx == nil {
		panic("JqEvalContext must be called within QueryContext")
	}

	evalCtx = base.NewChild()
	evalCtx.Functions = map[string]function.Function{
		funcName: function.New(&function.Spec{
			Params: []function.Parameter{
				{
					Name:        "query",
					Description: "The jq query string",
					Type:        JqQueryType.Type(),
				},
			},
			Type: function.StaticReturnType(definitions.QueryResultPlaceholderType.Type()),
			Impl: func(args []cty.Value, retType cty.Type) (cty.Value, error) {
				// The parsing is done by the customdecode on JqType
				// Here we just type-erasing the query placeholder
				return definitions.QueryResultPlaceholderType.ToCty(
					JqQueryType.MustFromCty(args[0]),
				), nil
			},
		}),
	}
	return
}
