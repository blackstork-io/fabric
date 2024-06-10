package pluginapiv1

import (
	"fmt"

	"github.com/zclconf/go-cty/cty"

	"github.com/blackstork-io/fabric/eval/dataquery"
	"github.com/blackstork-io/fabric/pkg/utils"
	"github.com/blackstork-io/fabric/plugin"
)

func encodeCtyType(src cty.Type) (*CtyType, error) {
	switch {
	case src.IsPrimitiveType():
		return encodeCtyPrimitiveType(src)
	case src.IsListType() || src.IsMapType() || src.IsSetType():
		return encodeCtySingleType(src)
	case src.IsObjectType():
		return encodeCtyObjectType(src)
	case src.IsTupleType():
		return encodeCtyTupleType(src)
	case src.Equals(cty.DynamicPseudoType):
		return &CtyType{
			Data: &CtyType_DynamicPseudo{
				DynamicPseudo: &CtyDynamicPseudoType{},
			},
		}, nil
	case src.IsCapsuleType():
		return encodeCtyCapsuleType(src)
	default:
		return nil, fmt.Errorf("unsupported cty type: %s", src.FriendlyName())
	}
}

func encodeCtyCapsuleType(src cty.Type) (*CtyType, error) {
	switch {
	case plugin.EncapsulatedData.CtyTypeEqual(src):
		return &CtyType{
			Data: &CtyType_Encapsulated{
				Encapsulated: CtyCapsuleType_CAPSULE_PLUGIN_DATA,
			},
		}, nil
	case dataquery.DelayedEvalType.CtyTypeEqual(src):
		return &CtyType{
			Data: &CtyType_Encapsulated{
				Encapsulated: CtyCapsuleType_CAPSULE_DELAYED_EVAL,
			},
		}, nil
	default:
		return nil, fmt.Errorf("unsupported capsule cty type: %s", src.FriendlyName())
	}
}

func encodeCtyPrimitiveType(src cty.Type) (*CtyType, error) {
	var kind CtyPrimitiveType
	switch src {
	case cty.Bool:
		kind = CtyPrimitiveType_KIND_BOOL
	case cty.Number:
		kind = CtyPrimitiveType_KIND_NUMBER
	case cty.String:
		kind = CtyPrimitiveType_KIND_STRING
	default:
		return nil, fmt.Errorf("unsupported primitive cty type: %s", src.FriendlyName())
	}

	return &CtyType{
		Data: &CtyType_Primitive{
			Primitive: kind,
		},
	}, nil
}

func encodeCtySingleType(src cty.Type) (*CtyType, error) {
	elemType, err := encodeCtyType(src.ElementType())
	if err != nil {
		return nil, err
	}

	switch {
	case src.IsListType():
		return &CtyType{
			Data: &CtyType_List{
				List: elemType,
			},
		}, nil
	case src.IsMapType():
		return &CtyType{
			Data: &CtyType_Map{
				Map: elemType,
			},
		}, nil
	case src.IsSetType():
		return &CtyType{
			Data: &CtyType_Set{
				Set: elemType,
			},
		}, nil
	default:
		panic("unreachable")
	}
}

func encodeCtyObjectType(src cty.Type) (*CtyType, error) {
	attrs, err := utils.MapMapErr(src.AttributeTypes(), encodeCtyType)
	if err != nil {
		return nil, err
	}

	return &CtyType{
		Data: &CtyType_Object{
			Object: &CtyObjectType{
				Attrs: attrs,
			},
		},
	}, nil
}

func encodeCtyTupleType(src cty.Type) (*CtyType, error) {
	elements, err := utils.FnMapErr(src.TupleElementTypes(), encodeCtyType)
	if err != nil {
		return nil, err
	}
	return &CtyType{
		Data: &CtyType_Tuple{
			Tuple: &CtyTupleType{
				Elements: elements,
			},
		},
	}, nil
}
