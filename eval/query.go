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
	jqQuery, err := gojq.Parse(q.Value.AsString())
	if err != nil {
		return nil, diagnostics.Diag{{
			Severity: hcl.DiagError,
			Summary:  "Failed to parse the query",
			Detail:   err.Error(),
			Subject:  &q.SrcRange,
		}}
	}

	code, err := gojq.Compile(jqQuery)
	if err != nil {
		return nil, diagnostics.Diag{{
			Severity: hcl.DiagError,
			Summary:  "Failed to compile the query",
			Detail:   err.Error(),
			Subject:  &q.SrcRange,
		}}
	}
	res, hasResult := code.Run(dataCtx.Any()).Next()
	if !hasResult {
		return nil, nil
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
