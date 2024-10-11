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
	// Nil if not present.
	Items *hclsyntax.Attribute
	// Condition is a condition that must be met for the block to be included.
	// Always present, defaults to true.
	Condition *hclsyntax.Attribute
	Content   []*ParsedContent
}
