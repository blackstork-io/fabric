package utils

import (
	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/hclsyntax"
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
