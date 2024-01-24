package definitions

import (
	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/hcldec"
	"github.com/hashicorp/hcl/v2/hclsyntax"
	"github.com/zclconf/go-cty/cty"
	"github.com/zclconf/go-cty/cty/convert"

	"github.com/blackstork-io/fabric/parser/evaluation"
	"github.com/blackstork-io/fabric/pkg/diagnostics"
)

// Desugars `title = "foo"` into appropriate `context` invocation.
type titleInvocation hclsyntax.Attribute

// GetBody implements evaluation.Invocation.
func (t *titleInvocation) GetBody() *hclsyntax.Body {
	return &hclsyntax.Body{
		SrcRange: t.SrcRange,
		EndRange: hcl.Range{
			Filename: t.SrcRange.Filename,
			Start:    t.SrcRange.End,
			End:      t.SrcRange.End,
		},
	}
}

// SetBody implements evaluation.Invocation.
func (*titleInvocation) SetBody(*hclsyntax.Body) {
	return
}

var _ evaluation.Invocation = (*titleInvocation)(nil)

func (t *titleInvocation) DefRange() hcl.Range {
	return t.SrcRange
}

func (t *titleInvocation) MissingItemRange() hcl.Range {
	return t.SrcRange
}

// Range implements InvocationObject.
func (t *titleInvocation) Range() hcl.Range {
	return t.SrcRange
}

func (t *titleInvocation) ParseInvocation(spec hcldec.Spec) (val cty.Value, diags diagnostics.Diag) {
	// Titles can only be rendered once, so there's no reason to put `sync.Once` like in proper blocks

	titleVal, diag := t.Expr.Value(nil)
	if diags.ExtendHcl(diag) {
		return
	}

	titleStrVal, err := convert.Convert(titleVal, cty.String)
	if err != nil {
		diags.Append(&hcl.Diagnostic{
			Severity: hcl.DiagError,
			Summary:  "Failed to turn title into a string",
			Detail:   err.Error(),
			Subject:  t.Expr.Range().Ptr(),
		})
		return
	}
	// cty.MapVal()?
	val = cty.ObjectVal(map[string]cty.Value{
		"text":      titleStrVal,
		"format_as": cty.StringVal("title"),
	})
	return
}

func NewTitle(title *hclsyntax.Attribute) *titleInvocation {
	return (*titleInvocation)(title)
}
