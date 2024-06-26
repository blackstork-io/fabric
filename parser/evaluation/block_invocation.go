package evaluation

import (
	"context"

	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/hclsyntax"
	"github.com/zclconf/go-cty/cty"

	"github.com/blackstork-io/fabric/cmd/fabctx"
	"github.com/blackstork-io/fabric/pkg/diagnostics"
	"github.com/blackstork-io/fabric/plugin/dataspec"
)

type BlockInvocation struct {
	*hclsyntax.Body
	DefinitionRange hcl.Range
}

// GetBody implements Invocation.
func (b *BlockInvocation) GetBody() *hclsyntax.Body {
	return b.Body
}

// SetBody implements Invocation.
func (b *BlockInvocation) SetBody(body *hclsyntax.Body) {
	b.Body = body
}

// DefRange implements Invocation.
func (b *BlockInvocation) DefRange() hcl.Range {
	return b.DefinitionRange
}

// ParseInvocation implements Invocation.
func (b *BlockInvocation) ParseInvocation(ctx context.Context, spec dataspec.RootSpec) (cty.Value, diagnostics.Diag) {
	return dataspec.Decode(b.Body, spec, fabctx.GetEvalContext(ctx))
}

// Range implements Invocation.
func (b *BlockInvocation) Range() hcl.Range {
	return b.Body.Range()
}

var _ Invocation = (*BlockInvocation)(nil)
