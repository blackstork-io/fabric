package evaluation

import (
	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/hcldec"
	"github.com/hashicorp/hcl/v2/hclsyntax"
	"github.com/zclconf/go-cty/cty"

	"github.com/blackstork-io/fabric/pkg/diagnostics"
)

type BlockInvocation struct {
	*hclsyntax.Body
	DefinitionRange hcl.Range
}

// DefRange implements Invocation.
func (b *BlockInvocation) DefRange() hcl.Range {
	return b.DefinitionRange
}

// Parse implements Invocation.
func (b *BlockInvocation) Parse(spec hcldec.Spec) (cty.Value, diagnostics.Diag) {
	res, diag := hcldec.Decode(b.Body, spec, nil)
	return res, diagnostics.Diag(diag)
}

// Range implements Invocation.
func (b *BlockInvocation) Range() hcl.Range {
	return b.Body.Range()
}

var _ Invocation = (*BlockInvocation)(nil)
