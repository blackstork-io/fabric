package encapsulator

import (
	"fmt"
	"reflect"

	"github.com/zclconf/go-cty/cty"
)

type Encapsulator[T any] struct {
	typ cty.Type
}

// Checks if the value is a valid, non-null cty capsule value.
func (e *Encapsulator[T]) ValIs(v cty.Value) bool {
	if v.IsNull() || !v.IsKnown() {
		return false
	}
	return e.Is(v.Type())
}

// Checks if the given ty is the same as encapsulated type.
func (e *Encapsulator[T]) Is(ty cty.Type) bool {
	return e.typ.Equals(ty)
}

// Converts val into corresponding cty capsule value.
func (e *Encapsulator[T]) Type() cty.Type {
	return e.typ
}

// Converts val into corresponding cty capsule value.
func (e *Encapsulator[T]) ToCty(val *T) cty.Value {
	return cty.CapsuleVal(e.typ, val)
}

// Extracts the encapsulated T from cty.Value. Panics if the value incorrect.
func (e *Encapsulator[T]) MustFromCty(v cty.Value) *T {
	return v.EncapsulatedValue().(*T)
}

// Extracts the encapsulated T from cty.Value. Returns an error if the value is incorrect.
func (e *Encapsulator[T]) FromCty(v cty.Value) (*T, error) {
	if v.IsNull() {
		return nil, fmt.Errorf("null value")
	}
	typ := v.Type()
	if !typ.Equals(e.typ) {
		return nil, fmt.Errorf("wrong type: expected %s, got %s", e.typ.FriendlyName(), typ.FriendlyName())
	}
	return v.EncapsulatedValue().(*T), nil
}

func New[T any](friendlyName string) *Encapsulator[T] {
	return &Encapsulator[T]{
		typ: cty.Capsule(friendlyName, reflect.TypeFor[T]()),
	}
}

func NewWithOps[T any](friendlyName string, ops *cty.CapsuleOps) *Encapsulator[T] {
	return &Encapsulator[T]{
		typ: cty.CapsuleWithOps(friendlyName, reflect.TypeFor[T](), ops),
	}
}
