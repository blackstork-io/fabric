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
				ObjDump: &ObjDumpSpec{
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
				Opaque: &OpaqueSpec{
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
func encodeAttr(src *dataspec.AttrSpec) (*AttrSpec, error) {
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
	return &AttrSpec{
		Name:       src.Name,
		Type:       ty,
		DefaultVal: dv,
		ExampleVal: ev,
		Doc:        src.Doc,
	}, nil
}

func encodeBlock(src *dataspec.BlockSpec) (*BlockSpec, error) {
	nested, err := encodeSpec(src.Nested)
	if err != nil {
		return nil, err
	}
	return &BlockSpec{
		Name:     src.Name,
		Nested:   nested,
		Doc:      src.Doc,
		Required: src.Required,
	}, nil
}

func encodeObject(src dataspec.ObjectSpec) (*ObjectSpec, error) {
	res := make([]*ObjectSpec_ObjectSpecChild, 0, len(src))
	for _, s := range src {
		switch sT := s.(type) {
		case *dataspec.AttrSpec:
			v, err := encodeAttr(sT)
			if err != nil {
				return nil, err
			}
			res = append(res, &ObjectSpec_ObjectSpecChild{
				Data: &ObjectSpec_ObjectSpecChild_Attr{
					Attr: v,
				},
			})
		case *dataspec.BlockSpec:
			v, err := encodeBlock(sT)
			if err != nil {
				return nil, err
			}
			res = append(res, &ObjectSpec_ObjectSpecChild{
				Data: &ObjectSpec_ObjectSpecChild_Block{
					Block: v,
				},
			})
		case *dataspec.KeyForObjectSpec:
			v, err := encodeSpec(sT.Spec)
			if err != nil {
				return nil, err
			}
			res = append(res, &ObjectSpec_ObjectSpecChild{
				Data: &ObjectSpec_ObjectSpecChild_Named{
					Named: &KeyForObjectSpec{
						Spec: v,
						Key:  sT.KeyForObjectSpec(),
					},
				},
			})
		default:
			return nil, fmt.Errorf("unsupported named spec: %T", src)
		}
	}

	return &ObjectSpec{
		Specs: res,
	}, nil
}
