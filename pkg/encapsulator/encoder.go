package encapsulator

import (
	"reflect"

	"github.com/zclconf/go-cty/cty"
)

// Non-generic Encoder interface
type EncoderI interface {
	CtyType() cty.Type
	CtyTypeEqual(ty cty.Type) bool
	ValCtyTypeEqual(v cty.Value) bool
	toCty(val any) cty.Value
}

type encoderCore struct {
	ctyType cty.Type
}

func (e *encoderCore) toCty(val any) cty.Value {
	return cty.CapsuleVal(e.ctyType, val)
}

// CtyType returns the cty type of the encoded values.
func (e *encoderCore) CtyType() cty.Type {
	return e.ctyType
}

// CtyTypeEqual checks if the given cty type can be decoded into T.
// This is the strictest check, requiring that the given cty type was produced by this Encoder.
func (e *encoderCore) CtyTypeEqual(ty cty.Type) bool {
	return e.ctyType.Equals(ty)
}

// ValCtyTypeEqual checks if the given cty value can be decoded into T.
// This is the strictest check, requiring that the given cty value was produced by this Encoder.
func (e *encoderCore) ValCtyTypeEqual(v cty.Value) bool {
	v, _ = v.Unmark()
	return isValid(v) && e.CtyTypeEqual(v.Type())
}

// NewEncoder creates an Encoder.
// Will use cty.Capsule if capsuleOps is nil, and cty.CapsuleWithOps otherwise
func newEncoderCore(goType reflect.Type, friendlyName string, capsuleOps *cty.CapsuleOps) encoderCore {
	var ctyType cty.Type

	if capsuleOps == nil {
		ctyType = cty.Capsule(friendlyName, goType)
	} else {
		ctyType = cty.CapsuleWithOps(friendlyName, goType, capsuleOps)
	}
	return encoderCore{
		ctyType: ctyType,
	}
}

// Encoder encapsulates *T into cty.Value.
type Encoder[T any] struct {
	encoderCore
}

// ToCty encapsulates *T into cty.Value.
func (e *Encoder[T]) ToCty(val *T) cty.Value {
	return e.toCty(val)
}

// ValToCty is a convenience function that encapsulates &val into cty.Value.
func (e *Encoder[T]) ValToCty(val T) cty.Value {
	return e.toCty(&val)
}

// NewEncoder creates an Encoder.
// Will use cty.Capsule if capsuleOps is nil, and cty.CapsuleWithOps otherwise
func NewEncoder[T any](friendlyName string, capsuleOps *CapsuleOps[T]) *Encoder[T] {
	return &Encoder[T]{
		encoderCore: newEncoderCore(
			reflect.TypeFor[T](),
			friendlyName,
			capsuleOps.asCtyCapsuleOps(),
		),
	}
}
