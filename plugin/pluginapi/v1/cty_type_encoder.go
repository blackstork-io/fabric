package pluginapiv1

import (
	"fmt"

	"github.com/zclconf/go-cty/cty"
)

func encodeCtyType(src cty.Type) (*CtyType, error) {
	switch {
	case src.IsPrimitiveType():
		return encodeCtyPrimitiveType(src)
	case src.IsListType():
		return encodeCtyListType(src)
	case src.IsMapType():
		return encodeCtyMapType(src)
	case src.IsSetType():
		return encodeCtySetType(src)
	case src.IsObjectType():
		return encodeCtyObjectType(src)
	case src.IsTupleType():
		return encodeCtyTupleType(src)
	case src.Equals(cty.DynamicPseudoType):
		return encodeCtyDynamicPseudoType(src)
	default:
		return nil, fmt.Errorf("unsupported cty type: %s", src.FriendlyName())
	}
}

func encodeCtyDynamicPseudoType(src cty.Type) (*CtyType, error) {
	return &CtyType{
		Data: &CtyType_DynamicPseudo{
			DynamicPseudo: &CtyDynamicPseudoType{},
		},
	}, nil
}

func encodeCtyPrimitiveType(src cty.Type) (*CtyType, error) {
	kind := CtyPrimitiveKind_CTY_PRIMITIVE_KIND_UNSPECIFIED
	switch src {
	case cty.Bool:
		kind = CtyPrimitiveKind_CTY_PRIMITIVE_KIND_BOOL
	case cty.Number:
		kind = CtyPrimitiveKind_CTY_PRIMITIVE_KIND_NUMBER
	case cty.String:
		kind = CtyPrimitiveKind_CTY_PRIMITIVE_KIND_STRING
	}
	if kind == CtyPrimitiveKind_CTY_PRIMITIVE_KIND_UNSPECIFIED {
		return nil, fmt.Errorf("unsupported primitive cty type: %s", src.FriendlyName())
	}
	return &CtyType{
		Data: &CtyType_Primitive{
			Primitive: &CtyPrimitiveType{
				Kind: kind,
			},
		},
	}, nil
}

func encodeCtyListType(src cty.Type) (*CtyType, error) {
	elemType, err := encodeCtyType(src.ElementType())
	if err != nil {
		return nil, err
	}
	return &CtyType{
		Data: &CtyType_List{
			List: &CtyListType{
				Element: elemType,
			},
		},
	}, nil
}

func encodeCtyMapType(src cty.Type) (*CtyType, error) {
	elemType, err := encodeCtyType(src.ElementType())
	if err != nil {
		return nil, err
	}
	return &CtyType{
		Data: &CtyType_Map{
			Map: &CtyMapType{
				Element: elemType,
			},
		},
	}, nil
}

func encodeCtySetType(src cty.Type) (*CtyType, error) {
	elemType, err := encodeCtyType(src.ElementType())
	if err != nil {
		return nil, err
	}
	return &CtyType{
		Data: &CtyType_Set{
			Set: &CtySetType{
				Element: elemType,
			},
		},
	}, nil
}

func encodeCtyObjectType(src cty.Type) (*CtyType, error) {
	srcAttrs := src.AttributeTypes()
	dstAttrs := make(map[string]*CtyType, len(srcAttrs))
	for k, v := range srcAttrs {
		ct, err := encodeCtyType(v)
		if err != nil {
			return nil, err
		}
		dstAttrs[k] = ct
	}
	return &CtyType{
		Data: &CtyType_Object{
			Object: &CtyObjectType{
				Attrs: dstAttrs,
			},
		},
	}, nil
}

func encodeCtyTupleType(src cty.Type) (*CtyType, error) {
	srcElems := src.TupleElementTypes()
	dstElems := make([]*CtyType, len(srcElems))
	for i, v := range srcElems {
		ct, err := encodeCtyType(v)
		if err != nil {
			return nil, err
		}
		dstElems[i] = ct
	}
	return &CtyType{
		Data: &CtyType_Tuple{
			Tuple: &CtyTupleType{
				Elements: dstElems,
			},
		},
	}, nil
}
