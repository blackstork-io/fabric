package definitions

import (
	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/hclsyntax"
	"github.com/zclconf/go-cty/cty"

	"github.com/blackstork-io/fabric/parser/evaluation"
	"github.com/blackstork-io/fabric/pkg/utils"
)

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
			Invocation: &evaluation.BlockInvocation{
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
