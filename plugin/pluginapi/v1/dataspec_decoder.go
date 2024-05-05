package pluginapiv1

import (
	"fmt"

	"github.com/blackstork-io/fabric/pkg/utils"
	"github.com/blackstork-io/fabric/plugin/dataspec"
	"github.com/blackstork-io/fabric/plugin/dataspec/constraint"
)

func decodeRootSpec(src *Spec) (dataspec.RootSpec, error) {
	sp, err := decodeSpec(src)
	if err != nil {
		return nil, err
	}
	rs, ok := sp.(dataspec.RootSpec)
	if !ok {
		return nil, fmt.Errorf("attempted to encode non-root spec %T", sp)
	}
	return rs, nil
}

func decodeSpec(src *Spec) (dataspec.Spec, error) {
	switch data := src.GetData().(type) {
	case nil:
		return nil, nil
	case *Spec_Attr:
		return decodeAttrSpec(data.Attr)
	case *Spec_Block:
		return decodeBlockSpec(data.Block)
	case *Spec_ObjSpec:
		return decodeObjSpec(data.ObjSpec)
	case *Spec_ObjDump:
		return decodeObjDumpSpec(data.ObjDump)
	case *Spec_Opaque:
		return decodeOpaqueSpec(data.Opaque)
	default:
		return nil, fmt.Errorf("unsupported spec: %T", src)
	}
}

func decodeAttrSpec(src *AttrSpec) (*dataspec.AttrSpec, error) {
	t, err := decodeCtyType(src.GetType())
	if err != nil {
		return nil, err
	}
	def, err := decodeCtyValue(src.GetDefaultVal())
	if err != nil {
		return nil, err
	}
	ex, err := decodeCtyValue(src.GetExampleVal())
	if err != nil {
		return nil, err
	}
	oneof, err := utils.FnMapErr(src.GetOneOf(), decodeCtyValue)
	if err != nil {
		return nil, err
	}
	min, err := decodeCtyValue(src.GetMinInclusive())
	if err != nil {
		return nil, err
	}
	max, err := decodeCtyValue(src.GetMaxInclusive())
	if err != nil {
		return nil, err
	}
	return &dataspec.AttrSpec{
		Name:         src.GetName(),
		Type:         t,
		DefaultVal:   def,
		ExampleVal:   ex,
		Doc:          src.GetDoc(),
		Constraints:  constraint.Constraints(src.GetConstraints()),
		OneOf:        oneof,
		MinInclusive: min,
		MaxInclusive: max,
		Deprecated:   src.GetDeprecated(),
	}, nil
}

func decodeBlockSpec(src *BlockSpec) (*dataspec.BlockSpec, error) {
	nested, err := decodeSpec(src.GetNested())
	if err != nil {
		return nil, err
	}
	return &dataspec.BlockSpec{
		Name:     src.GetName(),
		Nested:   nested,
		Required: src.GetRequired(),
		Doc:      src.GetDoc(),
	}, nil
}

func decodeObjSpec(src *ObjectSpec) (dataspec.ObjectSpec, error) {
	encodedSpecs := src.GetSpecs()
	specs := make(dataspec.ObjectSpec, 0, len(encodedSpecs))
	for _, s := range encodedSpecs {
		switch sT := s.GetData().(type) {
		case nil:
			continue
		case *ObjectSpec_ObjectSpecChild_Named:
			parsedSpec, err := decodeSpec(sT.Named.GetSpec())
			if err != nil {
				return nil, err
			}
			specs = append(specs, dataspec.UnderKey(sT.Named.GetKey(), parsedSpec))
		case *ObjectSpec_ObjectSpecChild_Attr:
			parsedSpec, err := decodeAttrSpec(sT.Attr)
			if err != nil {
				return nil, err
			}
			specs = append(specs, parsedSpec)
		case *ObjectSpec_ObjectSpecChild_Block:
			parsedSpec, err := decodeBlockSpec(sT.Block)
			if err != nil {
				return nil, err
			}
			specs = append(specs, parsedSpec)
		default:
			return nil, fmt.Errorf("unsupported named spec: %T", src)
		}
	}

	return specs, nil
}

func decodeObjDumpSpec(objDump *ObjDumpSpec) (*dataspec.ObjDumpSpec, error) {
	return &dataspec.ObjDumpSpec{
		Doc: objDump.GetDoc(),
	}, nil
}

func decodeOpaqueSpec(opaque *OpaqueSpec) (*dataspec.OpaqueSpec, error) {
	res := &dataspec.OpaqueSpec{
		Doc: opaque.GetDoc(),
	}
	var err error
	res.Spec, err = decodeHclSpec(opaque.GetSpec())
	if err != nil {
		return nil, err
	}
	return res, nil
}
