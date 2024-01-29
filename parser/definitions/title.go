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
type TitleInvocation struct {
	hcl.Expression
}

// GetBody implements evaluation.Invocation.
func (t *TitleInvocation) GetBody() *hclsyntax.Body {
	rng := t.Expression.Range()
	return &hclsyntax.Body{
		SrcRange: rng,
		EndRange: hcl.Range{
			Filename: rng.Filename,
			Start:    rng.End,
			End:      rng.End,
		},
	}
}

// SetBody implements evaluation.Invocation.
func (*TitleInvocation) SetBody(*hclsyntax.Body) {
	return
}

var _ evaluation.Invocation = (*TitleInvocation)(nil)

func (t *TitleInvocation) DefRange() hcl.Range {
	return t.Expression.Range()
}

func (t *TitleInvocation) MissingItemRange() hcl.Range {
	return t.Expression.Range()
}

func (t *TitleInvocation) ParseInvocation(spec hcldec.Spec) (val cty.Value, diags diagnostics.Diag) {
	// Titles can only be rendered once, so there's no reason to put `sync.Once` like in proper blocks

	titleVal, diag := t.Expression.Value(nil)
	if diags.ExtendHcl(diag) {
		return
	}

	titleStrVal, err := convert.Convert(titleVal, cty.String)
	if err != nil {
		diags.Append(&hcl.Diagnostic{
			Severity: hcl.DiagError,
			Summary:  "Failed to turn title into a string",
			Detail:   err.Error(),
			Subject:  t.Range().Ptr(),
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

func NewTitle(title hcl.Expression) *TitleInvocation {
	return &TitleInvocation{
		Expression: title,
	}
}
