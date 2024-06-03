package pluginapiv1

import (
	"fmt"

	"github.com/zclconf/go-cty/cty"

	"github.com/blackstork-io/fabric/pkg/utils"
	"github.com/blackstork-io/fabric/plugin"
)

func decodeCtyValue(src *CtyValue) (cty.Value, error) {
	if src == nil {
		return cty.NilVal, nil
	}
	t, err := decodeCtyType(src.GetType())
	if err != nil {
		return cty.NilVal, err
	}
	switch {
	case t.IsPrimitiveType() && src.GetPrimitive() != nil:
		return decodeCtyPrimitiveValue(src.GetPrimitive())
	case t.IsObjectType() || t.IsMapType():
		return decodeCtyMapLike(t, src.GetMapLike().GetElements())
	case t.IsListType() || t.IsSetType() || t.IsTupleType():
		return decodeCtyListLike(t, src.GetListLike().GetElements())
	case plugin.EncapsulatedData.CtyTypeEqual(t):
		if data := src.GetEncapsulated().GetValue(); data != nil {
			return plugin.EncapsulatedData.ValToCty(decodeData(data)), nil
		}
	}
	return cty.NullVal(t), nil
}

func decodeCtyMapLike(t cty.Type, src map[string]*CtyValue) (cty.Value, error) {
	if src == nil {
		return cty.NullVal(t), nil
	}
	elements, err := utils.MapMapErr(src, decodeCtyValue)
	if err != nil {
		return cty.NilVal, err
	}
	if t.IsObjectType() {
		return cty.ObjectVal(elements), nil
	} else if t.IsMapType() {
		return cty.MapVal(elements), nil
	}
	return cty.NilVal, fmt.Errorf("Unsupported cty map-like type: %s", t.FriendlyName())
}

func decodeCtyListLike(t cty.Type, src []*CtyValue) (cty.Value, error) {
	if src == nil {
		return cty.NullVal(t), nil
	}
	elements, err := utils.FnMapErr(src, decodeCtyValue)
	if err != nil {
		return cty.NilVal, err
	}
	if t.IsListType() {
		return cty.ListVal(elements), nil
	} else if t.IsSetType() {
		return cty.SetVal(elements), nil
	} else if t.IsTupleType() {
		return cty.TupleVal(elements), nil
	}
	return cty.NilVal, fmt.Errorf("Unsupported cty list-like type: %s", t.FriendlyName())
}

func decodeCtyPrimitiveValue(src *CtyPrimitiveValue) (cty.Value, error) {
	switch data := src.GetData().(type) {
	case *CtyPrimitiveValue_Bln:
		return cty.BoolVal(data.Bln), nil
	case *CtyPrimitiveValue_Num:
		return cty.NumberFloatVal(data.Num), nil
	case *CtyPrimitiveValue_Str:
		return cty.StringVal(data.Str), nil
	default:
		return cty.NilVal, fmt.Errorf("unsupported primitive cty type: %T", src)
	}
}
