package definitions

import (
	"github.com/hashicorp/hcl/v2/hclsyntax"
)

type Dynamic struct {
	ParseResult *ParsedDynamic
}

type ParsedDynamic struct {
	Block *hclsyntax.Block
	// Items is a list of items to be iterated over dynamically.
	// Always present.
	Items   *hclsyntax.Attribute
	Content []*ParsedContent
}
