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
