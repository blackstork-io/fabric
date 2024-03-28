package blocks

import (
	"reflect"
	"strings"

	"github.com/zclconf/go-cty/cty"
)

func capsuleTypeFor[V any]() cty.Type {
	ty := reflect.TypeOf((*V)(nil)).Elem()
	return cty.Capsule(
		strings.ToLower(ty.Name()),
		ty,
	)
}
