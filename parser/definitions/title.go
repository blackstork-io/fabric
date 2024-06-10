package definitions

import (
	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/hclsyntax"
	"github.com/zclconf/go-cty/cty"

	"github.com/blackstork-io/fabric/parser/evaluation"
	"github.com/blackstork-io/fabric/pkg/diagnostics"
	"github.com/blackstork-io/fabric/plugin/dataspec"
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

func (t *titleInvocation) ParseInvocation(spec dataspec.RootSpec) (val cty.Value, diags diagnostics.Diag) {
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
		"value": &hclsyntax.Attribute{
			Name: "value",
			Expr: expr,
		},
		"relative_size": &hclsyntax.Attribute{
			Name: "relative_size",
			Expr: &hclsyntax.LiteralValueExpr{
				Val:      cty.NumberIntVal(-1),
				SrcRange: t.Expression.Range(),
			},
		},
	}

	val, diag := dataspec.Decode(body, spec, evaluation.EvalContext())
	diags.Extend(diag)
	return
}

func NewTitle(title *hclsyntax.Attribute, resolver ConfigResolver) *ParsedContent {
	pluginName := "title"
	return &ParsedContent{
		Plugin: &ParsedPlugin{
			PluginName: pluginName,
			Config:     resolver(BlockKindContent, pluginName),
			Invocation: &titleInvocation{
				Expression: title.Expr,
			},
		},
	}
}
