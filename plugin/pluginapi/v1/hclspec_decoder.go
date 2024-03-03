package pluginapiv1

import (
	"fmt"

	"github.com/hashicorp/hcl/v2/hcldec"
)

func decodeHclSpec(src *HclSpec) (hcldec.Spec, error) {
	switch {
	case src == nil || src.GetData() == nil:
		return nil, nil
	case src.GetBlock() != nil:
		return decodeHclBlockSpec(src.GetBlock())
	case src.GetBlockMap() != nil:
		return decodeHclBlockMapSpec(src.GetBlockMap())
	case src.GetDefault() != nil:
		return decodeHclDefaultSpec(src.GetDefault())
	case src.GetLiteral() != nil:
		return decodeHclLiteralSpec(src.GetLiteral())
	case src.GetBlockList() != nil:
		return decodeHclBlockListSpec(src.GetBlockList())
	case src.GetBlockSet() != nil:
		return decodeHclBlockSetSpec(src.GetBlockSet())
	case src.GetAttr() != nil:
		return decodeHclAttrSpec(src.GetAttr())
	case src.GetObject() != nil:
		return decodeHclObjectSpec(src.GetObject())
	case src.GetBlockAttrs() != nil:
		return decodeHclBlockAttrsSpec(src.GetBlockAttrs())
	default:
		return nil, fmt.Errorf("unsupported hcl spec: %T", src)
	}
}

func decodeHclBlockAttrsSpec(src *HclBlockAttrs) (*hcldec.BlockAttrsSpec, error) {
	t, err := decodeCtyType(src.GetType())
	if err != nil {
		return nil, err
	}
	return &hcldec.BlockAttrsSpec{
		TypeName:    src.GetName(),
		Required:    src.GetRequired(),
		ElementType: t,
	}, nil
}

func decodeHclObjectSpec(src *HclObject) (hcldec.Spec, error) {
	dst := make(hcldec.ObjectSpec)
	for k, a := range src.GetAttrs() {
		attr, err := decodeHclSpec(a)
		if err != nil {
			return nil, err
		}
		dst[k] = attr
	}
	return dst, nil
}

func decodeHclAttrSpec(src *HclAttr) (*hcldec.AttrSpec, error) {
	t, err := decodeCtyType(src.GetType())
	if err != nil {
		return nil, err
	}
	return &hcldec.AttrSpec{
		Name:     src.GetName(),
		Type:     t,
		Required: src.GetRequired(),
	}, nil
}

func decodeHclBlockSetSpec(src *HclBlockSet) (*hcldec.BlockSetSpec, error) {
	nested, err := decodeHclSpec(src.GetNested())
	if err != nil {
		return nil, err
	}
	return &hcldec.BlockSetSpec{
		TypeName: src.GetName(),
		Nested:   nested,
		MinItems: int(src.GetMinItems()),
		MaxItems: int(src.GetMaxItems()),
	}, nil
}

func decodeHclBlockListSpec(src *HclBlockList) (*hcldec.BlockListSpec, error) {
	nested, err := decodeHclSpec(src.GetNested())
	if err != nil {
		return nil, err
	}
	return &hcldec.BlockListSpec{
		TypeName: src.GetName(),
		Nested:   nested,
		MinItems: int(src.GetMaxItems()),
		MaxItems: int(src.GetMaxItems()),
	}, nil
}

func decodeHclLiteralSpec(src *HclLiteral) (*hcldec.LiteralSpec, error) {
	value, err := decodeCtyValue(src.GetValue())
	if err != nil {
		return nil, err
	}
	return &hcldec.LiteralSpec{
		Value: value,
	}, nil
}

func decodeHclDefaultSpec(src *HclDefault) (*hcldec.DefaultSpec, error) {
	def, err := decodeHclSpec(src.GetDefault())
	if err != nil {
		return nil, err
	}
	prm, err := decodeHclSpec(src.GetPrimary())
	if err != nil {
		return nil, err
	}
	return &hcldec.DefaultSpec{
		Primary: prm,
		Default: def,
	}, nil
}

func decodeHclBlockMapSpec(src *HclBlockMap) (*hcldec.BlockMapSpec, error) {
	nested, err := decodeHclSpec(src.GetNested())
	if err != nil {
		return nil, err
	}
	return &hcldec.BlockMapSpec{
		TypeName:   src.GetName(),
		Nested:     nested,
		LabelNames: src.GetLabels(),
	}, nil
}

func decodeHclBlockSpec(src *HclBlock) (*hcldec.BlockSpec, error) {
	nested, err := decodeHclSpec(src.GetNested())
	if err != nil {
		return nil, err
	}
	return &hcldec.BlockSpec{
		TypeName: src.GetName(),
		Nested:   nested,
		Required: src.GetRequired(),
	}, nil
}
