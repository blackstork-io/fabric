package pluginapiv1

import (
	"fmt"

	"github.com/zclconf/go-cty/cty"
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
	case t.IsListType() && src.GetList() != nil:
		return decodeCtyListValue(src.GetList())
	case t.IsMapType() && src.GetMap() != nil:
		return decodeCtyMapValue(src.GetMap())
	case t.IsSetType() && src.GetSet() != nil:
		return decodeCtySetValue(src.GetSet())
	case t.IsObjectType() && src.GetObject() != nil:
		return decodeCtyObjectValue(src.GetObject())
	case t.IsTupleType() && src.GetTuple() != nil:
		return decodeCtyTupleValue(src.GetTuple())
	default:
		return cty.NullVal(t), nil
	}
}

func decodeCtyTupleValue(src *CtyTupleValue) (cty.Value, error) {
	elements := make([]cty.Value, len(src.GetElements()))
	var err error
	for i, elem := range src.GetElements() {
		elements[i], err = decodeCtyValue(elem)
		if err != nil {
			return cty.NilVal, err
		}
	}
	return cty.TupleVal(elements), nil
}

func decodeCtyObjectValue(src *CtyObjectValue) (cty.Value, error) {
	attrs := make(map[string]cty.Value, len(src.GetAttrs()))
	var err error
	for k, v := range src.GetAttrs() {
		attrs[k], err = decodeCtyValue(v)
		if err != nil {
			return cty.NilVal, err
		}
	}
	return cty.ObjectVal(attrs), nil
}

func decodeCtySetValue(src *CtySetValue) (cty.Value, error) {
	elements := make([]cty.Value, len(src.GetElements()))
	var err error
	for i, elem := range src.GetElements() {
		elements[i], err = decodeCtyValue(elem)
		if err != nil {
			return cty.NilVal, err
		}
	}
	return cty.SetVal(elements), nil
}

func decodeCtyMapValue(src *CtyMapValue) (cty.Value, error) {
	elements := make(map[string]cty.Value, len(src.GetElements()))
	var err error
	for k, v := range src.GetElements() {
		elements[k], err = decodeCtyValue(v)
		if err != nil {
			return cty.NilVal, err
		}
	}
	return cty.MapVal(elements), nil
}

func decodeCtyListValue(src *CtyListValue) (cty.Value, error) {
	elements := make([]cty.Value, len(src.GetElements()))
	var err error
	for i, elem := range src.GetElements() {
		elements[i], err = decodeCtyValue(elem)
		if err != nil {
			return cty.NilVal, err
		}
	}
	return cty.ListVal(elements), nil
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
