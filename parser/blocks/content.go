package blocks

import (
	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/hclsyntax"
	"github.com/zclconf/go-cty/cty"

	"github.com/blackstork-io/fabric/parser/blocks/internal/tree"
	"github.com/blackstork-io/fabric/pkg/diagnostics"
)

// type Query struct {
// 	code *gojq.Code
// 	rng  *hcl.Range
// }

type ContentPlugin struct {
	tree.NodeSigil
	plugin
	// query *Query
}

func (*ContentPlugin) FriendlyName() string {
	return "content block"
}

var contentCtyType = capsuleTypeFor[ContentPlugin]()

func (p *ContentPlugin) CtyType() cty.Type {
	return contentCtyType
}

func (p *ContentPlugin) AsCtyValue() cty.Value {
	return cty.CapsuleVal(contentCtyType, p)
}

func DefineContentPlugin(block *hclsyntax.Block, atTopLevel bool) (res *ContentPlugin, diags diagnostics.Diag) {
	plugin, diags := definePlugin(block, atTopLevel)
	if diags.HasErrors() {
		return
	}
	res = &ContentPlugin{
		plugin: plugin,
	}
	return
}

// func (p *ContentPlugin) parseQuery(ctx interface {
// 	GetEvalContext() *hcl.EvalContext
// 	GetDefinedBlocks() *DefinedBlocks
// 	Traverser(expr hclsyntax.Expression) (tree.Node, diagnostics.Diag)
// }, queryAttr *hclsyntax.Attribute,
// ) (diags diagnostics.Diag) {
// 	val, dgs := queryAttr.Expr.Value(ctx.GetEvalContext())
// 	if diags.ExtendHcl(dgs) {
// 		return
// 	}
// 	strVal, err := convert.Convert(val, cty.String)
// 	if err != nil {
// 		diags.Append(&hcl.Diagnostic{
// 			Severity: hcl.DiagError,
// 			Summary:  "Error in query",
// 			Detail:   "Can't convert to string: " + err.Error(),
// 			Subject:  queryAttr.Expr.Range().Ptr(),
// 		})
// 		return
// 	}
// 	jqQuery, err := gojq.Parse(strVal.AsString())
// 	if err != nil {
// 		diags.Append(&hcl.Diagnostic{
// 			Severity: hcl.DiagError,
// 			Summary:  "Error in query",
// 			Detail:   "Failed to parse the query: " + err.Error(),
// 			Subject:  queryAttr.Expr.Range().Ptr(),
// 		})
// 		return
// 	}

// 	code, err := gojq.Compile(jqQuery)
// 	if err != nil {
// 		diags.Append(&hcl.Diagnostic{
// 			Severity: hcl.DiagError,
// 			Summary:  "Error in query",
// 			Detail:   "Failed to compile the query: " + err.Error(),
// 			Subject:  queryAttr.Expr.Range().Ptr(),
// 		})
// 		return
// 	}
// 	p.query = &Query{
// 		code: code,
// 		rng:  queryAttr.Expr.Range().Ptr(),
// 	}
// 	return
// }

func (p *ContentPlugin) Parse(ctx interface {
	GetEvalContext() *hcl.EvalContext
	GetDefinedBlocks() *DefinedBlocks
	Traverser(expr hclsyntax.Expression) (tree.Node, diagnostics.Diag)
},
) (diags diagnostics.Diag) {
	if diags.Extend(p.plugin.fillParsedBody(ctx)) {
		return
	}
	// query, _ := utils.Pop(p.ParsedBody.Attributes, AttrQuery)
	// if query != nil {
	// 	if diags.Extend(p.parseQuery(ctx, query)) {
	// 		return
	// 	}
	// }

	diags = p.plugin.parseSpecial(ctx)
	if diags.HasErrors() {
		return
	}
	// if p.plugin.Base != nil {
	// 	if p.query == nil {
	// 		p.query = p.plugin.Base.(*ContentPlugin).query
	// 	}
	// }
	return
}
