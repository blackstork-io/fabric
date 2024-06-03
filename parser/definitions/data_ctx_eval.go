package definitions

import (
	"context"

	"github.com/hashicorp/hcl/v2"
	"github.com/zclconf/go-cty/cty"

	"github.com/blackstork-io/fabric/pkg/diagnostics"
	"github.com/blackstork-io/fabric/pkg/encapsulator"
	"github.com/blackstork-io/fabric/plugin"
)

// DataCtxEvalNeeded is an interface that should be implemented by any cty capsule type
// that needs to be evaluated with data context. Example of such type is jq query or row_vals in table.
type DataCtxEvalNeeded interface {
	Eval(ctx context.Context, dataCtx plugin.MapData) (result cty.Value, diags diagnostics.Diag)
	Range() *hcl.Range
}

var DataCtxEvalNeededType = encapsulator.NewDecoder[DataCtxEvalNeeded]()
