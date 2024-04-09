// Documentable wrapper tupes form
package dataspec

import (
	"github.com/hashicorp/hcl/v2/hcldec"
	"github.com/hashicorp/hcl/v2/hclwrite"
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
	f := hclwrite.NewEmptyFile()
	spec.WriteDoc(f.Body().AppendNewBlock(blockName, labels).Body())
	return string(hclwrite.Format(f.Bytes()))
}

// Spec is a documentable wrapper over `hcldec.Spec`s.
type Spec interface {
	HcldecSpec() hcldec.Spec
	WriteDoc(w *hclwrite.Body)
	getSpec() Spec
	Validate() []string
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
