package pluginapiv1

import (
	"fmt"

	"github.com/zclconf/go-cty/cty"

	"github.com/blackstork-io/fabric/eval/dataquery"
	"github.com/blackstork-io/fabric/pkg/utils"
	"github.com/blackstork-io/fabric/plugin"
)

func decodeCtyType(src *CtyType) (cty.Type, error) {
	switch src := src.Data.(type) {
	case *CtyType_Primitive:
		return decodeCtyPrimitiveType(src.Primitive)
	case *CtyType_List:
		return decodeCtyListType(src.List)
	case *CtyType_Map:
		return decodeCtyMapType(src.Map)
	case *CtyType_Set:
		return decodeCtySetType(src.Set)
	case *CtyType_Object:
		return decodeCtyObjectType(src.Object)
	case *CtyType_Tuple:
		return decodeCtyTupleType(src.Tuple)
	case *CtyType_DynamicPseudo:
		return cty.DynamicPseudoType, nil
	case *CtyType_Encapsulated:
		switch src.Encapsulated {
		case CtyCapsuleType_CAPSULE_DELAYED_EVAL:
			return dataquery.DelayedEvalType.CtyType(), nil
		case CtyCapsuleType_CAPSULE_PLUGIN_DATA:
			return plugin.EncapsulatedData.CtyType(), nil
		default:
			return cty.NilType, fmt.Errorf("unsupported capsule cty type: %v", src.Encapsulated.String())
		}
	default:
		return cty.NilType, fmt.Errorf("unsupported cty type: %T", src)
	}
}

func decodeCtyPrimitiveType(src CtyPrimitiveType) (cty.Type, error) {
	switch src {
	case CtyPrimitiveType_KIND_BOOL:
		return cty.Bool, nil
	case CtyPrimitiveType_KIND_NUMBER:
		return cty.Number, nil
	case CtyPrimitiveType_KIND_STRING:
		return cty.String, nil
	default:
		return cty.NilType, fmt.Errorf("unsupported primitive cty type: %v", src)
	}
}

func decodeCtyListType(src *CtyType) (cty.Type, error) {
	elemType, err := decodeCtyType(src)
	if err != nil {
		return cty.NilType, err
	}
	return cty.List(elemType), nil
}

func decodeCtyMapType(src *CtyType) (cty.Type, error) {
	elemType, err := decodeCtyType(src)
	if err != nil {
		return cty.NilType, err
	}
	return cty.Map(elemType), nil
}

func decodeCtySetType(src *CtyType) (cty.Type, error) {
	elemType, err := decodeCtyType(src)
	if err != nil {
		return cty.NilType, err
	}
	return cty.Set(elemType), nil
}

func decodeCtyObjectType(src *CtyObjectType) (cty.Type, error) {
	attrTypes, err := utils.MapMapErr(src.GetAttrs(), decodeCtyType)
	if err != nil {
		return cty.NilType, err
	}
	return cty.Object(attrTypes), nil
}

func decodeCtyTupleType(src *CtyTupleType) (cty.Type, error) {
	elemTypes, err := utils.FnMapErr(src.GetElements(), decodeCtyType)
	if err != nil {
		return cty.NilType, err
	}
	return cty.Tuple(elemTypes), nil
}
