package pluginapiv1

import (
	"fmt"

	"github.com/hashicorp/hcl/v2"

	"github.com/blackstork-io/fabric/pkg/diagnostics"
	"github.com/blackstork-io/fabric/plugin/dataspec"
)

func encodeSpec(src dataspec.Spec) (*Spec, diagnostics.Diag) {
	switch sp := src.(type) {
	case nil:
		return nil, nil
	case *dataspec.AttrSpec:
		attr, diags := encodeAttr(sp)
		return &Spec{
			Data: &Spec_Attr{
				Attr: attr,
			},
		}, diags
	case *dataspec.BlockSpec:
		block, diags := encodeBlock(sp)
		return &Spec{
			Data: &Spec_Block{
				Block: block,
			},
		}, diags
	case dataspec.ObjectSpec:
		obj, diags := encodeObject(sp)
		return &Spec{
			Data: &Spec_ObjSpec{
				ObjSpec: obj,
			},
		}, diags
	case *dataspec.ObjDumpSpec:
		return &Spec{
			Data: &Spec_ObjDump{
				ObjDump: &ObjDumpSpec{
					Doc: sp.Doc,
				},
			},
		}, nil
	case *dataspec.OpaqueSpec:
		v, err := encodeHclSpec(sp.Spec)
		return &Spec{
			Data: &Spec_Opaque{
				Opaque: &OpaqueSpec{
					Doc:  sp.Doc,
					Spec: v,
				},
			},
		}, diagnostics.FromErr(err, "Failed to encode hcl spec")
	default:
		return nil, diagnostics.Diag{{
			Severity: hcl.DiagError,
			Summary:  "Unsupported spec",
			Detail:   fmt.Sprintf("unsupported spec: %T", src),
		}}
	}
}

// case *spec.KeyForObjectSpecr:
func encodeAttr(src *dataspec.AttrSpec) (*AttrSpec, diagnostics.Diag) {
	ty, diags := encodeCtyType(src.Type)
	dv, diag := encodeCtyValue(src.DefaultVal)
	diags.Extend(diag)
	ev, diag := encodeCtyValue(src.ExampleVal)
	diags.Extend(diag)
	oneof := make([]*CtyValue, len(src.OneOf))
	for i, v := range src.OneOf {
		oneof[i], diag = encodeCtyValue(v)
		diags.Extend(diag)
	}
	diags.Extend(diag)
	min, diag := encodeCtyValue(src.MinInclusive)
	diags.Extend(diag)
	max, diag := encodeCtyValue(src.MaxInclusive)
	diags.Extend(diag)
	if diags.HasErrors() {
		return nil, diags
	}
	return &AttrSpec{
		Name:         src.Name,
		Type:         ty,
		DefaultVal:   dv,
		ExampleVal:   ev,
		Doc:          src.Doc,
		Constraints:  uint32(src.Constraints),
		OneOf:        oneof,
		MinInclusive: min,
		MaxInclusive: max,
		Deprecated:   src.Deprecated,
		Secret:       src.Secret,
	}, diags
}

func encodeBlock(src *dataspec.BlockSpec) (*BlockSpec, diagnostics.Diag) {
	nested, diags := encodeSpec(src.Nested)
	if diags.HasErrors() {
		return nil, diags
	}
	return &BlockSpec{
		Name:     src.Name,
		Nested:   nested,
		Doc:      src.Doc,
		Required: src.Required,
	}, diags
}

func encodeObject(src dataspec.ObjectSpec) (_ *ObjectSpec, diags diagnostics.Diag) {
	res := make([]*ObjectSpec_ObjectSpecChild, 0, len(src))
	for _, s := range src {
		switch sT := s.(type) {
		case *dataspec.AttrSpec:
			v, diag := encodeAttr(sT)
			diags.Extend(diag)
			res = append(res, &ObjectSpec_ObjectSpecChild{
				Data: &ObjectSpec_ObjectSpecChild_Attr{
					Attr: v,
				},
			})
		case *dataspec.BlockSpec:
			v, diag := encodeBlock(sT)
			diags.Extend(diag)
			res = append(res, &ObjectSpec_ObjectSpecChild{
				Data: &ObjectSpec_ObjectSpecChild_Block{
					Block: v,
				},
			})
		case *dataspec.KeyForObjectSpec:
			v, diag := encodeSpec(sT.Spec)
			diags.Extend(diag)
			res = append(res, &ObjectSpec_ObjectSpecChild{
				Data: &ObjectSpec_ObjectSpecChild_Named{
					Named: &KeyForObjectSpec{
						Spec: v,
						Key:  sT.KeyForObjectSpec(),
					},
				},
			})
		default:
			diags.Add("Unsupported spec", fmt.Sprintf("unsupported named spec: %T", src))
		}
	}
	if diags.HasErrors() {
		return nil, diags
	}
	return &ObjectSpec{
		Specs: res,
	}, diags
}
