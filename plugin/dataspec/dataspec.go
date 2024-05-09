// Documentable wrapper tupes form
package dataspec

import (
	"bytes"
	"strings"

	"github.com/hashicorp/hcl/v2/hcldec"
	"github.com/hashicorp/hcl/v2/hclwrite"

	"github.com/blackstork-io/fabric/pkg/diagnostics"
)

type rootSpecSigil struct{}

func (rootSpecSigil) isRootSpec() rootSpecSigil {
	return rootSpecSigil{}
}

// RootSpec represents valid specs in root position (as .Args, .Config values).
type RootSpec interface {
	Spec
	// Must be safe to call on nil receiver
	IsEmpty() bool
	isRootSpec() rootSpecSigil
}

// RenderDoc renders the block documentation for spec.
func RenderDoc(spec RootSpec, blockName string, labels ...string) string {
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

// Spec is a documentable wrapper over `hcldec.Spec`s.
type Spec interface {
	HcldecSpec() hcldec.Spec
	WriteDoc(w *hclwrite.Body)
	getSpec() Spec
	ValidateSpec() diagnostics.Diag
}

// Reperesents types that could be included as children in Object.
type ObjectSpecChild interface {
	Spec
	KeyForObjectSpec() string
}

var (
	_ Spec            = (*OpaqueSpec)(nil)
	_ ObjectSpecChild = (*AttrSpec)(nil)
	_ ObjectSpecChild = (*BlockSpec)(nil)
	_ ObjectSpecChild = (*KeyForObjectSpec)(nil)
	_ RootSpec        = (*ObjDumpSpec)(nil)
	_ RootSpec        = ObjectSpec{}
)
