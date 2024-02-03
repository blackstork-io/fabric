package pluginapiv1

import (
	"fmt"

	"github.com/hashicorp/hcl/v2/hcldec"
)

func encodeHclSpec(src hcldec.Spec) (*HclSpec, error) {
	if src == nil {
		return nil, nil
	}
	switch spec := src.(type) {
	case hcldec.ObjectSpec:
		return encodeHclObjectSpec(spec)
	case *hcldec.ObjectSpec:
		return encodeHclObjectSpec(*spec)
	case *hcldec.AttrSpec:
		return encodeHclAttrSpec(spec)
	case *hcldec.BlockAttrsSpec:
		return encodeHclBlockAttrsSpec(spec)
	case *hcldec.BlockSpec:
		return encodeHclBlockSpec(spec)
	case *hcldec.BlockMapSpec:
		return encodeHclBlockMapSpec(spec)
	case *hcldec.DefaultSpec:
		return encodeHclDefaultSpec(spec)
	case *hcldec.LiteralSpec:
		return encodeHclLiteralSpec(spec)
	case *hcldec.BlockSetSpec:
		return encodeHclBlockSetSpec(spec)
	case *hcldec.BlockListSpec:
		return encodeHclBlockListSpec(spec)
	default:
		return nil, fmt.Errorf("unsupported hcl spec: %T", src)
	}
}

func encodeHclBlockListSpec(src *hcldec.BlockListSpec) (*HclSpec, error) {
	nested, err := encodeHclSpec(src.Nested)
	if err != nil {
		return nil, err
	}
	return &HclSpec{
		Data: &HclSpec_BlockList{
			BlockList: &HclBlockList{
				Name:     src.TypeName,
				Nested:   nested,
				MinItems: int64(src.MinItems),
				MaxItems: int64(src.MaxItems),
			},
		},
	}, nil
}

func encodeHclBlockSetSpec(src *hcldec.BlockSetSpec) (*HclSpec, error) {
	nested, err := encodeHclSpec(src.Nested)
	if err != nil {
		return nil, err
	}
	return &HclSpec{
		Data: &HclSpec_BlockSet{
			BlockSet: &HclBlockSet{
				Name:     src.TypeName,
				Nested:   nested,
				MinItems: int64(src.MinItems),
				MaxItems: int64(src.MaxItems),
			},
		},
	}, nil
}

func encodeHclLiteralSpec(src *hcldec.LiteralSpec) (*HclSpec, error) {
	v, err := encodeCtyValue(src.Value)
	if err != nil {
		return nil, err
	}
	return &HclSpec{
		Data: &HclSpec_Literal{
			Literal: &HclLiteral{
				Value: v,
			},
		},
	}, nil
}

func encodeHclDefaultSpec(src *hcldec.DefaultSpec) (*HclSpec, error) {
	def, err := encodeHclSpec(src.Default)
	if err != nil {
		return nil, err
	}
	prm, err := encodeHclSpec(src.Primary)
	if err != nil {
		return nil, err
	}
	return &HclSpec{
		Data: &HclSpec_Default{
			Default: &HclDefault{
				Default: def,
				Primary: prm,
			},
		},
	}, nil
}

func encodeHclBlockMapSpec(src *hcldec.BlockMapSpec) (*HclSpec, error) {
	nested, err := encodeHclSpec(src.Nested)
	if err != nil {
		return nil, err
	}
	return &HclSpec{
		Data: &HclSpec_BlockMap{
			BlockMap: &HclBlockMap{
				Name:   src.TypeName,
				Nested: nested,
				Labels: src.LabelNames,
			},
		},
	}, nil
}

func encodeHclBlockSpec(src *hcldec.BlockSpec) (*HclSpec, error) {
	nested, err := encodeHclSpec(src.Nested)
	if err != nil {
		return nil, err
	}
	return &HclSpec{
		Data: &HclSpec_Block{
			Block: &HclBlock{
				Name:     src.TypeName,
				Required: src.Required,
				Nested:   nested,
			},
		},
	}, nil
}

func encodeHclBlockAttrsSpec(src *hcldec.BlockAttrsSpec) (*HclSpec, error) {
	t, err := encodeCtyType(src.ElementType)
	if err != nil {
		return nil, err
	}
	return &HclSpec{
		Data: &HclSpec_BlockAttrs{
			BlockAttrs: &HclBlockAttrs{
				Name:     src.TypeName,
				Type:     t,
				Required: src.Required,
			},
		},
	}, nil
}

func encodeHclObjectSpec(src hcldec.ObjectSpec) (*HclSpec, error) {
	dstAttrs := make(map[string]*HclSpec, len(src))
	var err error
	for k, v := range src {
		dstAttrs[k], err = encodeHclSpec(v)
		if err != nil {
			return nil, err
		}
	}
	return &HclSpec{
		Data: &HclSpec_Object{
			Object: &HclObject{
				Attrs: dstAttrs,
			},
		},
	}, nil
}

func encodeHclAttrSpec(src *hcldec.AttrSpec) (*HclSpec, error) {
	t, err := encodeCtyType(src.Type)
	if err != nil {
		return nil, err
	}
	return &HclSpec{
		Data: &HclSpec_Attr{
			Attr: &HclAttr{
				Name:     src.Name,
				Required: src.Required,
				Type:     t,
			},
		},
	}, nil
}
