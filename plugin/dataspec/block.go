package dataspec

import (
	"fmt"

	"github.com/hashicorp/hcl/v2/hcldec"
	"github.com/hashicorp/hcl/v2/hclwrite"
)

// BlockSpec represents a nested block (hcldec.BlockSpec).
type BlockSpec struct {
	Name     string
	Nested   Spec
	Doc      string
	Required bool
	// not supported yet
	labels []string
}

func (b *BlockSpec) getSpec() Spec {
	return b
}

func (b *BlockSpec) KeyForObjectSpec() string {
	return b.Name
}

func (b *BlockSpec) WriteDoc(w *hclwrite.Body) {
	tokens := comment(nil, b.Doc)
	if len(tokens) != 0 {
		tokens = appendCommentNewLine(tokens)
	}
	if b.Required {
		tokens = comment(tokens, "Required")
	} else {
		tokens = comment(tokens, "Optional")
	}
	w.AppendUnstructuredTokens(tokens)
	block := w.AppendNewBlock(b.Name, b.labels)
	b.Nested.WriteDoc(block.Body())
}

func (b *BlockSpec) HcldecSpec() hcldec.Spec {
	return &hcldec.BlockSpec{
		TypeName: b.Name,
		Nested:   b.Nested.HcldecSpec(),
		Required: b.Required,
	}
}

func (b *BlockSpec) ValidateSpec() (errs []string) {
	switch st := b.Nested.(type) {
	case ObjectSpec:
	case *OpaqueSpec:
	default:
		errs = append(errs, fmt.Sprintf("invalid nesting: %T within Block spec", st))
	}
	return
}
