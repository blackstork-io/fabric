package pluginapiv1

import (
	"fmt"

	"github.com/zclconf/go-cty/cty"
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
	default:
		return cty.NilType, fmt.Errorf("unsupported cty type: %T", src)
	}
}

func decodeCtyPrimitiveType(src *CtyPrimitiveType) (cty.Type, error) {
	switch src.GetKind() {
	case CtyPrimitiveKind_CTY_PRIMITIVE_KIND_BOOL:
		return cty.Bool, nil
	case CtyPrimitiveKind_CTY_PRIMITIVE_KIND_NUMBER:
		return cty.Number, nil
	case CtyPrimitiveKind_CTY_PRIMITIVE_KIND_STRING:
		return cty.String, nil
	default:
		return cty.NilType, fmt.Errorf("unsupported primitive cty type: %v", src.Kind)
	}
}

func decodeCtyListType(src *CtyListType) (cty.Type, error) {
	elemType, err := decodeCtyType(src.GetElement())
	if err != nil {
		return cty.NilType, err
	}
	return cty.List(elemType), nil
}

func decodeCtyMapType(src *CtyMapType) (cty.Type, error) {
	elemType, err := decodeCtyType(src.GetElement())
	if err != nil {
		return cty.NilType, err
	}
	return cty.Map(elemType), nil
}

func decodeCtySetType(src *CtySetType) (cty.Type, error) {
	elemType, err := decodeCtyType(src.GetElement())
	if err != nil {
		return cty.NilType, err
	}
	return cty.Set(elemType), nil
}

func decodeCtyObjectType(src *CtyObjectType) (cty.Type, error) {
	attrTypes := make(map[string]cty.Type)
	for name, attrType := range src.GetAttrs() {
		t, err := decodeCtyType(attrType)
		if err != nil {
			return cty.NilType, err
		}
		attrTypes[name] = t
	}
	return cty.Object(attrTypes), nil
}

func decodeCtyTupleType(src *CtyTupleType) (cty.Type, error) {
	elemTypes := make([]cty.Type, len(src.GetElements()))
	for i, elemType := range src.GetElements() {
		t, err := decodeCtyType(elemType)
		if err != nil {
			return cty.NilType, err
		}
		elemTypes[i] = t
	}
	return cty.Tuple(elemTypes), nil
}
