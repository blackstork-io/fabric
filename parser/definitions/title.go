package definitions

import (
	"context"

	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/hclsyntax"
	"github.com/zclconf/go-cty/cty"

	"github.com/blackstork-io/fabric/cmd/fabctx"
	"github.com/blackstork-io/fabric/parser/evaluation"
	"github.com/blackstork-io/fabric/pkg/diagnostics"
	"github.com/blackstork-io/fabric/pkg/utils"
	"github.com/blackstork-io/fabric/plugin/dataspec"
)

// Desugars `title = "foo"` into appropriate `context` invocation.
type titleInvocation struct {
	*hclsyntax.Block
}

// GetBody implements evaluation.Invocation.
func (t *titleInvocation) GetBody() *hclsyntax.Body {
	return t.Body
}

// SetBody implements evaluation.Invocation.
func (t *titleInvocation) SetBody(b *hclsyntax.Body) {
	t.Body = b
}

var _ evaluation.Invocation = (*titleInvocation)(nil)

func (t *titleInvocation) MissingItemRange() hcl.Range {
	return t.Body.EndRange
}

func (t *titleInvocation) ParseInvocation(ctx context.Context, spec *dataspec.RootSpec) (val *dataspec.Block, diags diagnostics.Diag) {
	// Titles can only be rendered once, so there's no reason to put `sync.Once` like in proper blocks
	return dataspec.Decode(t.Block, spec, fabctx.GetEvalContext(ctx))
}

func NewTitle(title *hclsyntax.Attribute, resolver ConfigResolver) *ParsedContent {
	const pluginName = "title"

	value := *title
	value.Name = "value"

	relativeSize := *title
	relativeSize.Name = "relative_size"
	relativeSize.Expr = &hclsyntax.LiteralValueExpr{
		Val:      cty.NumberIntVal(-1),
		SrcRange: title.Expr.Range(),
	}
	return &ParsedContent{
		Plugin: &ParsedPlugin{
			PluginName: pluginName,
			Config:     resolver(BlockKindContent, pluginName),
			Invocation: &titleInvocation{
				Block: &hclsyntax.Block{
					Type:        BlockKindContent,
					TypeRange:   title.NameRange,
					Labels:      []string{pluginName},
					LabelRanges: []hcl.Range{title.NameRange},
					Body: &hclsyntax.Body{
						Attributes: hclsyntax.Attributes{
							"value":         &value,
							"relative_size": &relativeSize,
						},
						SrcRange: title.SrcRange,
						EndRange: utils.RangeEnd(title.Expr.Range()),
					},
					OpenBraceRange:  utils.RangeStart(title.NameRange),
					CloseBraceRange: utils.RangeEnd(title.Expr.Range()),
				},
			},
		},
	}
}
