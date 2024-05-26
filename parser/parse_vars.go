package parser

import (
	"fmt"
	"slices"

	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/hclsyntax"
	"golang.org/x/exp/maps"

	"github.com/blackstork-io/fabric/parser/definitions"
	"github.com/blackstork-io/fabric/pkg/diagnostics"
)

func ParseVars(block *hclsyntax.Block) (vars []*hclsyntax.Attribute, diags diagnostics.Diag) {
	for _, subblock := range block.Body.Blocks {
		diags.Append(&hcl.Diagnostic{
			Severity: hcl.DiagWarning,
			Summary:  "Unsupported nesting",
			Detail:   fmt.Sprintf(`%s block does not support nested blocks, use nested maps instead`, definitions.BlockKindVars),
			Subject:  subblock.Range().Ptr(),
		})
	}
	// Attributes are in map form, we need to sort them in order of definition
	vars = maps.Values(block.Body.Attributes)
	slices.SortFunc(vars, func(a, b *hclsyntax.Attribute) int {
		return a.NameRange.Start.Byte - b.NameRange.Start.Byte
	})
	return
}
