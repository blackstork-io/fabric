package engine

import (
	"context"

	"github.com/hashicorp/hcl/v2/hclsyntax"
	"github.com/zclconf/go-cty/cty"

	"github.com/blackstork-io/fabric/eval"
	"github.com/blackstork-io/fabric/pkg/diagnostics"
	"github.com/blackstork-io/fabric/plugin/dataspec"
)

func (e *Engine) GetStdoutPublisher(ctx context.Context, format string) (act *eval.PluginPublishAction, diags diagnostics.Diag) {
	pub, found := e.PluginRunner().Publisher("stdout")

	if !found {
		diags.Add("Missing publisher", "stdout publisher not found")
		return
	}
	argsBlock, diag := dataspec.DecodeAndEvalBlock(ctx, &hclsyntax.Block{
		Body: &hclsyntax.Body{
			Attributes: hclsyntax.Attributes{
				"format": &hclsyntax.Attribute{
					Name: "format",
					Expr: &hclsyntax.LiteralValueExpr{
						Val: cty.StringVal(format),
					},
				},
			},
		},
	}, pub.Args, nil)
	if diags.Extend(diag) {
		return
	}
	return &eval.PluginPublishAction{
		Publisher: pub,
		PluginAction: &eval.PluginAction{
			Args: argsBlock,
		},
	}, diags
}
