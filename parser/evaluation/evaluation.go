package evaluation

import (
	"context"

	"github.com/hashicorp/hcl/v2"

	"github.com/blackstork-io/fabric/pkg/diagnostics"
	"github.com/blackstork-io/fabric/plugin/dataspec"
)

// To act as a plugin configuration struct must implement this interface.
type Configuration interface {
	ParseConfig(ctx context.Context, spec *dataspec.RootSpec) (*dataspec.Block, diagnostics.Diag)
	Range() hcl.Range
	Exists() bool
}
