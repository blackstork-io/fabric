package dataquery

import (
	"github.com/hashicorp/hcl/v2"
	"github.com/zclconf/go-cty/cty"

	"github.com/blackstork-io/fabric/parser/definitions"
)

func QueryContext(base *hcl.EvalContext) (evalCtx *hcl.EvalContext, queries *definitions.Queries) {
	queries = definitions.NewQueries()
	evalCtx = base.NewChild()
	evalCtx.Variables = map[string]cty.Value{
		definitions.QueryKey: definitions.QueriesType.ToCty(queries),
	}
	return
}
