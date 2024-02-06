package pluginapiv1

import (
	"fmt"

	"github.com/zclconf/go-cty/cty"
)

func encodeCtyValue(src cty.Value) (*CtyValue, error) {
	if src == cty.NilVal {
		return nil, nil
	}
	switch {
	case src.Type().IsPrimitiveType():
		return encodeCtyPrimitiveValue(src)
	case src.Type().IsListType():
		return encodeCtyListValue(src)
	case src.Type().IsMapType():
		return encodeCtyMapValue(src)
	case src.Type().IsSetType():
		return encodeCtySetValue(src)
	case src.Type().IsObjectType():
		return encodeCtyObjectValue(src)
	case src.Type().IsTupleType():
		return encodeCtyTupleValue(src)
	default:
		return nil, fmt.Errorf("unsupported cty value: %s", src.Type().FriendlyName())
	}
}

func encodeCtySetValue(src cty.Value) (*CtyValue, error) {
	t, err := encodeCtyType(src.Type())
	if err != nil {
		return nil, err
	}
	dst := CtyValue{
		Type: t,
	}
	if src.IsNull() {
		return &dst, nil
	}
	value := CtySetValue{}
	for it := src.ElementIterator(); it.Next(); {
		_, v := it.Element()
		elem, err := encodeCtyValue(v)
		if err != nil {
			return nil, err
		}
		value.Elements = append(value.Elements, elem)
	}
	dst.Data = &CtyValue_Set{
		Set: &value,
	}
	return &dst, nil
}

func encodeCtyTupleValue(src cty.Value) (*CtyValue, error) {
	t, err := encodeCtyType(src.Type())
	if err != nil {
		return nil, err
	}
	dst := CtyValue{
		Type: t,
	}
	if src.IsNull() {
		return &dst, nil
	}
	value := CtyTupleValue{}
	for it := src.ElementIterator(); it.Next(); {
		_, v := it.Element()
		elem, err := encodeCtyValue(v)
		if err != nil {
			return nil, err
		}
		value.Elements = append(value.Elements, elem)
	}
	dst.Data = &CtyValue_Tuple{
		Tuple: &value,
	}
	return &dst, nil
}

func encodeCtyMapValue(src cty.Value) (*CtyValue, error) {
	t, err := encodeCtyType(src.Type())
	if err != nil {
		return nil, err
	}
	dst := CtyValue{
		Type: t,
	}
	if src.IsNull() {
		return &dst, nil
	}
	value := CtyMapValue{
		Elements: make(map[string]*CtyValue),
	}
	for it := src.ElementIterator(); it.Next(); {
		k, v := it.Element()
		elem, err := encodeCtyValue(v)
		if err != nil {
			return nil, err
		}
		value.Elements[k.AsString()] = elem
	}
	dst.Data = &CtyValue_Map{
		Map: &value,
	}
	return &dst, nil
}

func encodeCtyObjectValue(src cty.Value) (*CtyValue, error) {
	t, err := encodeCtyType(src.Type())
	if err != nil {
		return nil, err
	}
	dst := CtyValue{
		Type: t,
	}
	if src.IsNull() {
		return &dst, nil
	}
	value := CtyObjectValue{
		Attrs: make(map[string]*CtyValue),
	}
	for it := src.ElementIterator(); it.Next(); {
		k, v := it.Element()
		elem, err := encodeCtyValue(v)
		if err != nil {
			return nil, err
		}
		value.Attrs[k.AsString()] = elem
	}
	dst.Data = &CtyValue_Object{
		Object: &value,
	}
	return &dst, nil
}

func encodeCtyListValue(src cty.Value) (*CtyValue, error) {
	t, err := encodeCtyType(src.Type())
	if err != nil {
		return nil, err
	}
	dst := CtyValue{
		Type: t,
	}
	if src.IsNull() {
		return &dst, nil
	}
	value := CtyListValue{}
	for it := src.ElementIterator(); it.Next(); {
		_, v := it.Element()
		elem, err := encodeCtyValue(v)
		if err != nil {
			return nil, err
		}
		value.Elements = append(value.Elements, elem)
	}
	dst.Data = &CtyValue_List{
		List: &value,
	}
	return &dst, nil
}

func encodeCtyPrimitiveValue(src cty.Value) (*CtyValue, error) {
	t, err := encodeCtyType(src.Type())
	if err != nil {
		return nil, err
	}
	dst := CtyValue{
		Type: t,
	}
	value := CtyPrimitiveValue{}
	switch {
	case src.IsNull():
		return &dst, nil
	case src.Type().Equals(cty.Bool):
		value.Data = &CtyPrimitiveValue_Bln{
			Bln: src.True(),
		}
	case src.Type().Equals(cty.String):
		value.Data = &CtyPrimitiveValue_Str{
			Str: src.AsString(),
		}
	case src.Type().Equals(cty.Number):
		if src.AsBigFloat().IsInt() {
			n, _ := src.AsBigFloat().Float64()
			value.Data = &CtyPrimitiveValue_Num{
				Num: n,
			}
		} else {
			return nil, fmt.Errorf("unsupported number cty value: %T", src)
		}
	default:
		return nil, fmt.Errorf("unsupported primitive cty value: %s", src.Type().FriendlyName())
	}
	dst.Data = &CtyValue_Primitive{
		Primitive: &value,
	}
	return &dst, nil
}
