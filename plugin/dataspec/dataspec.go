// Documentable wrapper types form
package dataspec

import (
	"bytes"
	"strings"

	"github.com/hashicorp/hcl/v2/hclwrite"

	"github.com/blackstork-io/fabric/pkg/diagnostics"
)

type (
	Blocks     []*Block
	Attributes map[string]*Attr
)

// Spec is the interface for all data specification types.
type Spec interface {
	WriteDoc(w *hclwrite.Body)
	ValidateSpec() diagnostics.Diag
}

var (
	_ Spec = (*AttrSpec)(nil)
	_ Spec = (*BlockSpec)(nil)
	_ Spec = (*RootSpec)(nil)
)

// RenderDoc renders the block documentation for spec.
func RenderDoc(spec BlockSpec, blockName string, labels ...string) string {
	// Special-casing the first line generation:
	// config "data" "csv" { -> config data csv {

	const placeholder = "P"
	var sb strings.Builder

	if strings.Contains(blockName, " ") {
		return "<error: block name contains spaces>"
	}

	f := hclwrite.NewEmptyFile()
	spec.WriteDoc(f.Body().AppendNewBlock(placeholder, nil).Body())
	doc := hclwrite.Format(f.Bytes())
	blockBodyStart := bytes.IndexByte(doc, '{')
	if blockBodyStart == -1 {
		return "<error: no block body generated>"
	}

	sb.WriteString(blockName)
	sb.WriteByte(' ')
	for _, label := range labels {
		hasSpaces := strings.Contains(blockName, " ")
		if hasSpaces {
			sb.WriteByte('"')
		}
		sb.WriteString(label)
		if hasSpaces {
			sb.WriteString(`" `)
		} else {
			sb.WriteByte(' ')
		}
	}
	sb.Write(doc[blockBodyStart:])

	return sb.String()
}
