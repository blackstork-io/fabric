package eval

import (
	"context"

	"github.com/hashicorp/hcl/v2"
	"github.com/itchyny/gojq"
	"github.com/zclconf/go-cty/cty"

	"github.com/blackstork-io/fabric/pkg/diagnostics"
	"github.com/blackstork-io/fabric/plugin"
)

type Query struct {
	Value    cty.Value
	SrcRange hcl.Range
}

func (q *Query) EvalQuery(ctx context.Context, dataCtx plugin.MapData) (plugin.Data, diagnostics.Diag) {
	query := q.Value.AsString()
	jqQuery, err := gojq.Parse(query)
	if err != nil {
		return nil, diagnostics.Diag{{
			Severity: hcl.DiagError,
			Summary:  "Failed to parse the query",
			Detail:   err.Error(),
			Subject:  &q.SrcRange,
			Extra: diagnostics.GoJQError{
				Err:   err,
				Query: query,
			},
		}}
	}

	code, err := gojq.Compile(jqQuery)
	if err != nil {
		return nil, diagnostics.Diag{{
			Severity: hcl.DiagError,
			Summary:  "Failed to compile the query",
			Detail:   err.Error(),
			Subject:  &q.SrcRange,
			Extra: diagnostics.GoJQError{
				Err:   err,
				Query: query,
			},
		}}
	}
	res, hasResult := code.Run(dataCtx.Any()).Next()
	if !hasResult {
		return nil, nil
	}
	var ok bool
	if err, ok = res.(error); ok {
		return nil, diagnostics.Diag{{
			Severity: hcl.DiagError,
			Summary:  "Failed to run the query",
			Detail:   err.Error(),
			Subject:  &q.SrcRange,
			Extra: diagnostics.GoJQError{
				Err:   err,
				Query: query,
			},
		}}
	}
	result, err := plugin.ParseDataAny(res)
	if err != nil {
		return nil, diagnostics.Diag{{
			Severity: hcl.DiagError,
			Summary:  "Incorrect query result type",
			Detail:   err.Error(),
			Subject:  &q.SrcRange,
		}}
	}
	return result, nil
}
