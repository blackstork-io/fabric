package evaluator

import "github.com/hashicorp/hcl/v2/hclsyntax"

// TODO: our

type DefinedDocument struct {
	block *hclsyntax.Block
}

func (d *DefinedDocument) FriendlyName() string {
	return "document block"
}
