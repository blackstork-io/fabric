package dataspec

import (
	"fmt"

	"github.com/hashicorp/hcl/v2/hcldec"
	"github.com/hashicorp/hcl/v2/hclwrite"
)

// Wraps hcldec.ObjectSpec
type ObjectSpec []ObjectSpecChild

func (o ObjectSpec) IsEmpty() bool {
	return len(o) == 0
}

func (o ObjectSpec) getSpec() Spec {
	return o
}

func (o ObjectSpec) isRootSpec() rootSpecSigil {
	return rootSpecSigil{}
}

func (o ObjectSpec) HcldecSpec() hcldec.Spec {
	res := make(hcldec.ObjectSpec, len(o))
	for _, v := range o {
		res[v.KeyForObjectSpec()] = v.HcldecSpec()
	}
	return res
}

func (o ObjectSpec) WriteDoc(w *hclwrite.Body) {
	for i, spec := range o {
		switch st := spec.getSpec().(type) {
		case *AttrSpec:
			if i != 0 {
				w.AppendNewline()
			}
			st.WriteDoc(w)
		case *BlockSpec:
			if i != 0 {
				w.AppendNewline()
			}
			st.WriteDoc(w)
		case *OpaqueSpec:
			if i != 0 {
				w.AppendNewline()
			}
			st.WriteDoc(w)
		}
	}
}

func (o ObjectSpec) ValidateSpec() (errs []string) {
	names := make(map[string]struct{}, len(o))
	for _, spec := range o {
		if _, found := names[spec.KeyForObjectSpec()]; found {
			errs = append(errs, fmt.Sprintf("name %q is repeated within ObjectSpec", spec.KeyForObjectSpec()))
		} else {
			names[spec.KeyForObjectSpec()] = struct{}{}
		}
		sp := spec.getSpec()
		errs = append(errs, sp.ValidateSpec()...)
		switch st := sp.(type) {
		case *AttrSpec:
		case *BlockSpec:
		case *OpaqueSpec:
		default:
			errs = append(errs, fmt.Sprintf("invalid nesting: %T within ObjectSpec", st))
		}
	}
	return
}
