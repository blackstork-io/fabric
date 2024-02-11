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

func hclBodyToVal(body *hclsyntax.Body) (val cty.Value, diags diagnostics.Diag) {
	// TODO: this makes a full dump of all of the attributes, not abiding by hidden
	// Think about ways to fix, or whether fix is needed
	obj := make(map[string]cty.Value, len(body.Attributes)+len(body.Blocks))
	for name, attr := range body.Attributes {
		attrVal, diag := attr.Expr.Value(nil)
		if diags.ExtendHcl(diag) {
			continue
		}
		obj[name] = attrVal
	}
	for _, block := range body.Blocks {
		blockVal, diag := hclBodyToVal(block.Body)
		if diags.Extend(diag) {
			continue
		}
		obj[block.Type] = blockVal
	}
	val = cty.ObjectVal(obj)
	return
}

// ParseInvocation implements Invocation.
func (b *BlockInvocation) ParseInvocation(spec hcldec.Spec) (cty.Value, diagnostics.Diag) {
	if spec == nil {
		res, err := hclBodyToVal(b.Body)
		return res, err
	}
	res, diag := hcldec.Decode(b.Body, spec, nil)
	return res, diagnostics.Diag(diag)
}

// Range implements Invocation.
func (b *BlockInvocation) Range() hcl.Range {
	return b.Body.Range()
}

var _ Invocation = (*BlockInvocation)(nil)
