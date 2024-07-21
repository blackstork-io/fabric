package evaluation

import (
	"context"

	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/hclsyntax"

	"github.com/blackstork-io/fabric/pkg/diagnostics"
	"github.com/blackstork-io/fabric/plugin/dataspec"
)

// To act as a plugin configuration struct must implement this interface.
type Configuration interface {
	ParseConfig(ctx context.Context, spec *dataspec.RootSpec) (*dataspec.Block, diagnostics.Diag)
	Range() hcl.Range
	Exists() bool
}

// To act as a plugin invocation (body of the plugin call block)
// struct must implement this interface.
type Invocation interface {
	GetBody() *hclsyntax.Body
	SetBody(body *hclsyntax.Body)
	ParseInvocation(ctx context.Context, spec *dataspec.RootSpec) (*dataspec.Block, diagnostics.Diag)
	Range() hcl.Range
	DefRange() hcl.Range
	MissingItemRange() hcl.Range
}
