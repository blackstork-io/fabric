package pluginapiv1

import (
	"fmt"

	"github.com/blackstork-io/fabric/plugin/dataspec"
)

func encodeSpec(src dataspec.Spec) (*Spec, error) {
	switch sp := src.(type) {
	case nil:
		return nil, nil
	case *dataspec.AttrSpec:
		attr, err := encodeAttr(sp)
		if err != nil {
			return nil, err
		}
		return &Spec{
			Data: &Spec_Attr{
				Attr: attr,
			},
		}, nil
	case *dataspec.BlockSpec:
		block, err := encodeBlock(sp)
		if err != nil {
			return nil, err
		}
		return &Spec{
			Data: &Spec_Block{
				Block: block,
			},
		}, nil
	case dataspec.ObjectSpec:
		obj, err := encodeObject(sp)
		if err != nil {
			return nil, err
		}
		return &Spec{
			Data: &Spec_ObjSpec{
				ObjSpec: obj,
			},
		}, nil
	case *dataspec.ObjDumpSpec:
		return &Spec{
			Data: &Spec_ObjDump{
				ObjDump: &ObjDump{
					Doc: sp.Doc,
				},
			},
		}, nil
	case *dataspec.OpaqueSpec:
		v, err := encodeHclSpec(sp.Spec)
		if err != nil {
			return nil, err
		}
		return &Spec{
			Data: &Spec_Opaque{
				Opaque: &Opaque{
					Doc:  sp.Doc,
					Spec: v,
				},
			},
		}, nil
	default:
		return nil, fmt.Errorf("unsupported spec: %T", src)
	}
}

// case *spec.KeyForObjectSpecr:
func encodeAttr(src *dataspec.AttrSpec) (*Attr, error) {
	ty, err := encodeCtyType(src.Type)
	if err != nil {
		return nil, err
	}
	dv, err := encodeCtyValue(src.DefaultVal)
	if err != nil {
		return nil, err
	}
	ev, err := encodeCtyValue(src.ExampleVal)
	if err != nil {
		return nil, err
	}
	return &Attr{
		Name:       src.Name,
		Type:       ty,
		DefaultVal: dv,
		ExampleVal: ev,
		Doc:        src.Doc,
		Required:   src.Required,
	}, nil
}

func encodeBlock(src *dataspec.BlockSpec) (*Block, error) {
	nested, err := encodeSpec(src.Nested)
	if err != nil {
		return nil, err
	}
	return &Block{
		Name:     src.Name,
		Nested:   nested,
		Doc:      src.Doc,
		Required: src.Required,
	}, nil
}

func encodeObject(src dataspec.ObjectSpec) (*Object, error) {
	res := make([]*NamedSpec, 0, len(src))
	for _, s := range src {
		switch sT := s.(type) {
		case *dataspec.AttrSpec:
			v, err := encodeAttr(sT)
			if err != nil {
				return nil, err
			}
			res = append(res, &NamedSpec{
				Data: &NamedSpec_Attr{
					Attr: v,
				},
			})
		case *dataspec.BlockSpec:
			v, err := encodeBlock(sT)
			if err != nil {
				return nil, err
			}
			res = append(res, &NamedSpec{
				Data: &NamedSpec_Block{
					Block: v,
				},
			})
		case *dataspec.KeyForObjectSpec:
			v, err := encodeSpec(sT.Spec)
			if err != nil {
				return nil, err
			}
			res = append(res, &NamedSpec{
				Data: &NamedSpec_Named{
					Named: &Namer{
						Spec: v,
						Name: sT.KeyForObjectSpec(),
					},
				},
			})
		default:
			return nil, fmt.Errorf("unsupported named spec: %T", src)
		}
	}

	return &Object{
		Specs: res,
	}, nil
}
