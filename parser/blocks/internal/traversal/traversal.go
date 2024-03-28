package traversal

import (
	"fmt"

	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/hclsyntax"
	"github.com/zclconf/go-cty/cty"
	"github.com/zclconf/go-cty/cty/gocty"

	"github.com/blackstork-io/fabric/parser/blocks/internal/tree"
	"github.com/blackstork-io/fabric/pkg/diagnostics"
)

func indexNode(n tree.Node, idx cty.Value, rng *hcl.Range) (tree.Node, diagnostics.Diag) {
	keyType := idx.Type()
	switch {
	case keyType.Equals(cty.String):
		i, ok := n.(tree.StrIndexable)
		if !ok {
			return nil, diagnostics.Diag{{
				Severity: hcl.DiagError,
				Summary:  "Index failed",
				Detail:   fmt.Sprintf("Indexing by string is not supported by '%s'", n.FriendlyName()),
				Subject:  rng,
			}}
		}
		n = i.IndexStr(idx.AsString())
		if n == nil {
			return nil, diagnostics.Diag{{
				Severity: hcl.DiagError,
				Summary:  "Index failed",
				Detail:   fmt.Sprintf("'%s' doesn't have an element with index %s", n.FriendlyName(), idx.GoString()),
				Subject:  rng,
			}}
		}
	case keyType.Equals(cty.Number):
		var intIdx int64
		i, ok := n.(tree.IntIndexable)
		if !ok {
			return nil, diagnostics.Diag{{
				Severity: hcl.DiagError,
				Summary:  "Index failed",
				Detail:   fmt.Sprintf("Indexing by a number is not supported by '%s'", n.FriendlyName()),
				Subject:  rng,
			}}
		}
		err := gocty.FromCtyValue(idx, &intIdx)
		if err != nil {
			return nil, diagnostics.Diag{{
				Severity: hcl.DiagError,
				Summary:  "Index failed",
				Detail:   "Incorrect index value: " + err.Error(),
				Subject:  rng,
			}}
		}
		n = i.IndexInt(intIdx)
		if n == nil {
			return nil, diagnostics.Diag{{
				Severity: hcl.DiagError,
				Summary:  "Index failed",
				Detail:   fmt.Sprintf("'%s' doesn't have an element with index %s", n.FriendlyName(), idx.GoString()),
				Subject:  rng,
			}}
		}
	default:
		return nil, diagnostics.Diag{{
			Severity: hcl.DiagError,
			Summary:  "Index failed",
			Detail:   fmt.Sprintf("cant use '%s' as an index", keyType.FriendlyName()),
			Subject:  rng,
		}}
	}
	return n, nil
}

type traversalEvaluator struct {
	ctx   *hcl.EvalContext
	diags diagnostics.Diag
	root  tree.Node
}

func tryExprToNode(expr hclsyntax.Expression) tree.Node {
	litExpr, ok := expr.(*hclsyntax.LiteralValueExpr)
	if !ok {
		return nil
	}
	return tryValueToNode(litExpr.Val)
}

func tryValueToNode(val cty.Value) tree.Node {
	if !val.Type().IsCapsuleType() {
		return nil
	}
	node, ok := val.EncapsulatedValue().(tree.Node)
	if !ok {
		return nil
	}
	return node
}

func (te *traversalEvaluator) EvalIndex(n *hclsyntax.IndexExpr) hclsyntax.Expression {
	newN := n
	newCollection := te.ResolveTraversals(n.Collection)
	if newCollection != n.Collection {
		tmp := *n
		newN = &tmp
		newN.Collection = newCollection
	}
	myNode := tryExprToNode(newCollection)
	if myNode == nil {
		return newN
	}

	keyVal, diags := n.Key.Value(te.ctx)
	if te.diags.ExtendHcl(diags) {
		return n
	}

	keyType := keyVal.Type()
	var idxRes tree.Node
	switch {
	case keyType.Equals(cty.String):
		i, ok := myNode.(tree.StrIndexable)
		if !ok {
			te.diags.Append(&hcl.Diagnostic{
				Severity: hcl.DiagError,
				Summary:  "Index failed",
				Detail:   fmt.Sprintf("Indexing by string is not supported by '%s'", myNode.FriendlyName()),
				Subject:  n.Range().Ptr(),
			})
			return n
		}
		idxRes = i.IndexStr(keyVal.AsString())
	case keyType.Equals(cty.Number):
		i, ok := myNode.(tree.IntIndexable)
		if !ok {
			te.diags.Append(&hcl.Diagnostic{
				Severity: hcl.DiagError,
				Summary:  "Index failed",
				Detail:   fmt.Sprintf("Indexing by a number is not supported by '%s'", myNode.FriendlyName()),
				Subject:  n.Range().Ptr(),
			})
			return n
		}

		var intIdx int64
		err := gocty.FromCtyValue(keyVal, &intIdx)
		if err != nil {
			te.diags.Append(&hcl.Diagnostic{
				Severity: hcl.DiagError,
				Summary:  "Index failed",
				Detail:   "Incorrect index value: " + err.Error(),
				Subject:  n.Range().Ptr(),
			})
			return n
		}
		idxRes = i.IndexInt(intIdx)
	default:
		te.diags.Append(&hcl.Diagnostic{
			Severity: hcl.DiagError,
			Summary:  "Failed to index",
			Detail:   fmt.Sprintf("Index of '%s' by '%s' could not be executed", myNode.FriendlyName(), keyType.FriendlyName()),
			Subject:  n.Range().Ptr(),
		})
		return n
	}
	if idxRes == nil {
		te.diags.Append(&hcl.Diagnostic{
			Severity: hcl.DiagError,
			Summary:  "Index failed",
			Detail:   fmt.Sprintf("'%s' doesn't have an element with index %s", myNode.FriendlyName(), keyVal.GoString()),
			Subject:  n.Range().Ptr(),
		})
		return n
	}
	return &hclsyntax.LiteralValueExpr{
		Val:      cty.CapsuleVal(tree.NodeType, idxRes),
		SrcRange: n.Range(),
	}
}

func (te *traversalEvaluator) EvalTraversal(trav hcl.Traversal, n tree.Node) (res hclsyntax.Expression) {
	var diags diagnostics.Diag
	for _, t := range trav {
		switch t := t.(type) {
		case hcl.TraverseRoot, hcl.TraverseAttr:
			var idx string
			switch tt := t.(type) {
			case hcl.TraverseRoot:
				idx = tt.Name
			case hcl.TraverseAttr:
				idx = tt.Name
			}

			n, diags = indexNode(n, cty.StringVal(idx), t.SourceRange().Ptr())
			if te.diags.Extend(diags) {
				return nil
			}
		case hcl.TraverseIndex:
			n, diags = indexNode(n, t.Key, &t.SrcRange)
			if te.diags.Extend(diags) {
				return nil
			}
		case hcl.TraverseSplat:
			te.diags.Append(&hcl.Diagnostic{
				Severity: hcl.DiagError,
				Summary:  "Unsupported operation",
				Detail:   "Splat operation is not supported",
				Subject:  &t.SrcRange,
			})
			return nil
		default:
			panic(fmt.Sprintf("unknown traverser type %T", t))
		}
	}
	return &hclsyntax.LiteralValueExpr{
		Val:      cty.CapsuleVal(tree.NodeType, n),
		SrcRange: trav.SourceRange(),
	}
}

func (te *traversalEvaluator) ResolveTraversals(n hclsyntax.Expression) hclsyntax.Expression {
	switch n := n.(type) {
	case *hclsyntax.ParenthesesExpr:
		expr := te.ResolveTraversals(n.Expression)
		if expr != n.Expression {
			newN := *n
			newN.Expression = expr
			return &newN
		}
		return n
	case *hclsyntax.ScopeTraversalExpr:
		res := te.EvalTraversal(n.Traversal, te.root)
		if res == nil {
			return n
		}
		return res

	case *hclsyntax.RelativeTraversalExpr:
		myNode := tryExprToNode(n.Source)
		if myNode == nil {
			return n
		}
		res := te.EvalTraversal(n.Traversal, myNode)
		if res == nil {
			return n
		}
		return res
	case *hclsyntax.IndexExpr:
		return te.EvalIndex(n)
	default:
		return n
	}
}

func ResolveTraversals(root tree.Node, expr hclsyntax.Expression, ctx *hcl.EvalContext) (tree.Node, diagnostics.Diag) {
	e := traversalEvaluator{
		ctx:  ctx,
		root: root,
	}
	newExpr := e.ResolveTraversals(expr)
	if e.diags.HasErrors() {
		return nil, e.diags
	}
	node := tryExprToNode(newExpr)
	if node != nil {
		return node, e.diags
	}
	val, diags := newExpr.Value(e.ctx)
	if e.diags.ExtendHcl(diags) {
		return nil, e.diags
	}
	node = tryValueToNode(val)
	if node != nil {
		return node, e.diags
	}
	e.diags.Append(&hcl.Diagnostic{
		Severity: hcl.DiagError,
		Summary:  "Failed to evaluate the reference",
		Detail:   fmt.Sprintf("Expected a node, got '%s'", val.Type().FriendlyName()),
	})
	return nil, e.diags
}
