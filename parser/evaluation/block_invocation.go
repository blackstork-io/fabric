package evaluation

import (
	"context"

	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/hclsyntax"

	"github.com/blackstork-io/fabric/cmd/fabctx"
	"github.com/blackstork-io/fabric/pkg/diagnostics"
	"github.com/blackstork-io/fabric/plugin/dataspec"
)

type BlockInvocation struct {
	*hclsyntax.Block
}

// MissingItemRange implements Invocation.
func (b *BlockInvocation) MissingItemRange() hcl.Range {
	return b.Body.MissingItemRange()
}

// GetBody implements Invocation.
func (b *BlockInvocation) GetBody() *hclsyntax.Body {
	return b.Body
}

// SetBody implements Invocation.
func (b *BlockInvocation) SetBody(body *hclsyntax.Body) {
	b.Body = body
}

// ParseInvocation implements Invocation.
func (b *BlockInvocation) ParseInvocation(ctx context.Context, spec *dataspec.RootSpec) (*dataspec.Block, diagnostics.Diag) {
	return dataspec.Decode(b.Block, spec, fabctx.GetEvalContext(ctx))
}

var _ Invocation = (*BlockInvocation)(nil)
