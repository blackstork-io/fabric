package gobfix

import (
	"encoding"
	"encoding/gob"
	"fmt"

	"github.com/hashicorp/hcl/v2/hcldec"
	"github.com/zclconf/go-cty/cty"

	"github.com/blackstork-io/fabric/pkg/utils"
)

type attrSpec struct {
	Name     string
	Type     typeWrapper
	Required bool
}

// toHcl implements specWrapper.
func (s *attrSpec) toHcl() hcldec.Spec {
	return &hcldec.AttrSpec{
		Name:     s.Name,
		Type:     s.Type.Type,
		Required: s.Required,
	}
}

type blockAttrsSpec struct {
	TypeName    string
	ElementType typeWrapper
	Required    bool
}

// toHcl implements specWrapper.
func (s *blockAttrsSpec) toHcl() hcldec.Spec {
	return &hcldec.BlockAttrsSpec{
		TypeName:    s.TypeName,
		ElementType: s.ElementType.Type,
		Required:    s.Required,
	}
}

type objectSpec map[string]Spec

// toHcl implements specWrapper.
func (s *objectSpec) toHcl() hcldec.Spec {
	res := make(hcldec.ObjectSpec, len(*s))
	for k, v := range *s {
		res[k] = v.toHcl()
	}
	return res
}

type Spec interface {
	toHcl() hcldec.Spec
}

type typeWrapper struct {
	cty.Type
}

// UnmarshalBinary implements encoding.BinaryUnmarshaler.
func (t *typeWrapper) UnmarshalBinary(data []byte) error {
	return t.Type.UnmarshalJSON(data)
}

// MarshalBinary implements encoding.BinaryMarshaler.
func (t *typeWrapper) MarshalBinary() (data []byte, err error) {
	return t.Type.MarshalJSON()
}

func objectSpecFromHcl(spec hcldec.ObjectSpec) objectSpec {
	if spec == nil {
		return nil
	}
	res := make(objectSpec, len(spec))
	for k, v := range spec {
		res[k] = FromHcl(v)
	}
	return res
}

func FromHcl(spec hcldec.Spec) Spec {
	if utils.IsNil(spec) {
		return nil
	}
	switch t := spec.(type) {
	case *hcldec.AttrSpec:
		return &attrSpec{
			Name:     t.Name,
			Type:     typeWrapper{t.Type},
			Required: t.Required,
		}
	case *hcldec.BlockAttrsSpec:
		return &blockAttrsSpec{
			TypeName:    t.TypeName,
			ElementType: typeWrapper{t.ElementType},
			Required:    t.Required,
		}
	case *hcldec.ObjectSpec:
		tmp := objectSpecFromHcl(*t)
		return &tmp
	default:
		panic(fmt.Sprintf("gobfix unimplemented for type %T", spec))
	}
}

func ToHcl(spec Spec) hcldec.Spec {
	if utils.IsNil(spec) {
		return nil
	}
	return spec.toHcl()
}

var (
	_ encoding.BinaryMarshaler   = &typeWrapper{}
	_ encoding.BinaryUnmarshaler = &typeWrapper{}

	_ Spec = &attrSpec{}
	_ Spec = &blockAttrsSpec{}
	_ Spec = &objectSpec{}
)

func init() {
	gob.Register(&attrSpec{})
	gob.Register(&blockAttrsSpec{})
	gob.Register(&objectSpec{})
}
