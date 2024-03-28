package evaluator

import (
	"context"
	"fmt"

	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/ext/customdecode"
	"github.com/hashicorp/hcl/v2/hcldec"
	"github.com/hashicorp/hcl/v2/hclsyntax"
	"github.com/zclconf/go-cty/cty"
	"github.com/zclconf/go-cty/cty/function"

	"github.com/blackstork-io/fabric/pkg/diagnostics"
)

type sharedEval struct {
	evalCtx *hcl.EvalContext
}

type eval struct {
	*sharedEval
	context.Context
	evalList *TracebackLinkedList
}

// Performs a shallow clone operation
func (e *eval) Clone() eval {
	return *e
}

// Shallow clones the eval and appends a new frame to it
func (e *eval) WithFrame(frame FrameIdentifier) eval {
	clone := *e
	clone.evalList = &TracebackLinkedList{
		parent: clone.evalList,
		frame:  frame,
	}
	return clone
}

type msgFrame string

func (m msgFrame) Eq(FrameIdentifier) bool {
	// message frames are informational-only, so they do not participate in the circular ref detection
	return false
}

type FriendlyNamer interface {
	FriendlyName() string
}

func Extract[V FriendlyNamer](expr hcl.Expression, evalCtx *hcl.EvalContext) (res V, diags diagnostics.Diag) {
	hcldec.ExprSpec

	val, diag := expr.Value(evalCtx)
	if diags.ExtendHcl(diag) {
		return
	}
	typ := val.Type()
	if typ.IsCapsuleType() {
		capVal := val.EncapsulatedValue()
		var ok bool
		res, ok = capVal.(V)
		if ok {
			return
		}
	}
	diags.Append(&hcl.Diagnostic{
		Severity:    hcl.DiagError,
		Summary:     "Incorrect type of expression",
		Detail:      fmt.Sprintf("Expected %s, got %s", res.FriendlyName(), typ.FriendlyName()),
		Subject:     expr.Range().Ptr(),
		Expression:  expr,
		EvalContext: evalCtx,
	})
	return
}

// top level command
func (e eval) Render(docPath string) (res string, diags diagnostics.Diag) {
	expr, diag := hclsyntax.ParseExpression([]byte(docPath), "<user-input>", hcl.InitialPos)
	if diags.ExtendHcl(diag) {
		return
	}
	evalCtx := e.sharedEval.evalCtx.NewChild()

	// variables like cur_doc, etc

	// wrong to set it here, just an example
	evalCtx.Functions = map[string]function.Function{
		"get": function.New(&function.Spec{
			Params: []function.Parameter{
				{
					Name:        "path",
					Description: "path to a block",
					Type:        customdecode.ExpressionClosureType,
				},
			},
			Impl: func(args []cty.Value, retType cty.Type) (cty.Value, error) {
				closure := customdecode.ExpressionClosureFromVal(args[0])
				// XXX: do the traversal on this thingy
				// closure.Expression
				// Must fully resolve (arbitrary exprs only in indexes)
				// while evaluating all unknowns with
				// closure.EvalContext
				_ = closure
				return cty.NilVal, nil
				// return evalWithLocals(vars, closure)
			},
		}),
	}

	// val, diag := Extract[*hclsyntax.Block](expr, e.evalCtx)
}

func (e eval) RenderDocument(doc *hclsyntax.Block) (res string, diags diagnostics.Diag) {
	// e.EnterFrame(&doc.Body.SrcRange)
	// par-exec-able
	for _, block := range doc.Body.Blocks {
		switch block.Type {
		case "content":
			res += e.CallContentPlugin(block)
		case "section":
			// res += e.RenderSection(block)
		default:
			// ignore all other blocks
			continue
		}
	}
	return
}

func (e eval) CallContentPlugin(plg *hclsyntax.Block) (res string) {
	// par-exec-able
	// e.EnterFrame(&plg.Body.SrcRange)
	val := e.renderedContent.Get(plg)
	val.Get(&e)

	// for _, block := range doc.Body.Blocks {
	// 	switch block.Type {
	// 	case "content":
	// 		res += e.CallContentPlugin(block)
	// 	case "section":
	// 		res += e.RenderSection(block)
	// 	default:
	// 		// ignore all other blocks
	// 		continue
	// 	}
	// }
	return
}

type Document struct {
	source    *hclsyntax.Block
	cv        *CachedVal
	pd        *ParsedDocument
	hadErrors bool
}

func (d *Document) parse(e eval) *ParsedDocument {
	d.cv.Enter(&e, false)
	defer d.cv.Exit()
	if d.pd != nil {
		return d.pd
	}
	if d.hadErrors {
		panic("hidden error not to duplicate the errors in the output")
	}

	// possibilities:
	// Wait
	// HasValue (read-return)
	// NoValue (uninit/failed)
	// CircularRef (returned by wait or calculate)

	// Actions:
	// Wait() -> Returns one of the other two, checks status
	// Calculate() -> Sets up for calculation, can return circ ref error
	//
}

type ParsedDocument struct{}
