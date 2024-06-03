package utils

import (
	"sync"
	"sync/atomic"

	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/hclsyntax"

	"github.com/blackstork-io/fabric/pkg/diagnostics"
)

func ToHclsyntaxBody(body hcl.Body) *hclsyntax.Body {
	hclsyntaxBody, ok := body.(*hclsyntax.Body)
	if !ok {
		// Should never happen: hcl.Body for hcl documents is always *hclsyntax.Body
		panic("hcl.Body to *hclsyntax.Body failed")
	}
	return hclsyntaxBody
}

func EvalContextByVar(ctx *hcl.EvalContext, name string) *hcl.EvalContext {
	for ; ctx != nil; ctx = ctx.Parent() {
		if ctx.Variables == nil {
			continue
		}
		_, found := ctx.Variables[name]
		if found {
			return ctx
		}
	}
	return nil
}

func EvalContextByFunc(ctx *hcl.EvalContext, name string) *hcl.EvalContext {
	for ; ctx != nil; ctx = ctx.Parent() {
		if ctx.Functions == nil {
			continue
		}
		_, found := ctx.Functions[name]
		if found {
			return ctx
		}
	}
	return nil
}

type onceVal[V any] struct {
	state atomic.Int32
	fn    func() (V, diagnostics.Diag)
	res   V
	mu    sync.Mutex
}

func (o *onceVal[V]) do() (res V, diags diagnostics.Diag) {
	state := o.state.Load()
	switch {
	case state > 0:
		return o.res, nil
	case state < 0:
		diags = diagnostics.Diag{diagnostics.RepeatedError}
		return
	}
	o.mu.Lock()
	defer func() {
		o.fn = nil
		if state == 0 {
			// this is a panic
			o.state.Store(-1)
		}
		o.mu.Unlock()
	}()
	state = o.state.Load()
	switch {
	case state > 0:
		res = o.res
	case state < 0:
		diags = diagnostics.Diag{diagnostics.RepeatedError}
	default:
		res, diags = o.fn()
		if diags.HasErrors() {
			state = -1
		} else {
			o.res = res
			state = 1
		}
		o.state.Store(state)
	}
	return
}

// OnceVal returns a function that calls fn only once and caches the result.
// If fn returns diagnostics with errors, the function will return it only once,
// on subsequent calls it will return RepeatedError.
func OnceVal[V any](fn func() (V, diagnostics.Diag)) func() (V, diagnostics.Diag) {
	return (&onceVal[V]{fn: fn}).do
}
