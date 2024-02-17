package definitions

import (
	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/hcldec"
	"github.com/hashicorp/hcl/v2/hclsyntax"
	"github.com/zclconf/go-cty/cty"

	"github.com/blackstork-io/fabric/parser/evaluation"
	"github.com/blackstork-io/fabric/pkg/diagnostics"
)

// Desugars `title = "foo"` into appropriate `context` invocation.
type titleInvocation struct {
	hcl.Expression
}

// GetBody implements evaluation.Invocation.
func (t *titleInvocation) GetBody() *hclsyntax.Body {
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
func (*titleInvocation) SetBody(*hclsyntax.Body) {}

var _ evaluation.Invocation = (*titleInvocation)(nil)

func (t *titleInvocation) DefRange() hcl.Range {
	return t.Expression.Range()
}

func (t *titleInvocation) MissingItemRange() hcl.Range {
	return t.Expression.Range()
}

func (t *titleInvocation) ParseInvocation(spec hcldec.Spec) (val cty.Value, diags diagnostics.Diag) {
	// Titles can only be rendered once, so there's no reason to put `sync.Once` like in proper blocks
	expr, ok := t.Expression.(hclsyntax.Expression)
	if !ok {
		diags.Append(&hcl.Diagnostic{
			Severity: hcl.DiagError,
			Summary:  "Incorrect title",
			Detail:   "Title must be an expression",
			Subject:  t.DefRange().Ptr(),
		})
		return
	}
	body := t.GetBody()
	body.Attributes = hclsyntax.Attributes{
		"text": &hclsyntax.Attribute{
			Name: "text",
			Expr: expr,
		},
		"format_as": &hclsyntax.Attribute{
			Name: "format_as",
			Expr: &hclsyntax.LiteralValueExpr{
				Val:      cty.StringVal("title"),
				SrcRange: t.Expression.Range(),
			},
		},
	}

	val, diag := hcldec.Decode(body, spec, nil)
	diags.ExtendHcl(diag)
	return
}

func NewTitle(title *hclsyntax.Attribute, resolver ConfigResolver) *ParsedContent {
	pluginName := "text"
	return &ParsedContent{
		PluginName: pluginName,
		Config:     resolver(BlockKindContent, pluginName),
		Invocation: &titleInvocation{
			Expression: title.Expr,
		},
	}
}
