package pluginapiv1

import (
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/zclconf/go-cty/cty"

	"github.com/blackstork-io/fabric/plugin"
)

func roundTrip(t *testing.T, val cty.Value) (decVal cty.Value, decTy cty.Type) {
	t.Helper()
	encVal, diag := encodeCtyValue(val)
	if diag.HasErrors() {
		t.Fatalf("Failed to encode value: %s", diag)
	}
	decVal, err := decodeCtyValue(encVal)
	if err != nil {
		t.Fatalf("Failed to decode value: %s", err)
	}
	if !decVal.RawEquals(val) {
		t.Fatalf("Roundtrip failed: got %s", decVal.GoString())
	}

	decTy = roundTripType(t, val.Type())
	return
}

func roundTripType(t *testing.T, ty cty.Type) (decTy cty.Type) {
	t.Helper()
	encTy, diag := encodeCtyType(ty)
	if diag.HasErrors() {
		t.Fatalf("Failed to encode type: %s", diag)
	}
	decTy, err := decodeCtyType(encTy)
	if err != nil {
		t.Fatalf("Failed to decode type: %s", err)
	}
	if !decTy.Equals(ty) {
		t.Fatalf("Roundtrip failed: got %s", decTy.GoString())
	}
	return
}

func TestCtyCodecRoundtrip(t *testing.T) {
	vals := []cty.Value{
		// all public cty.Value vars in cty
		cty.DynamicVal,
		cty.EmptyObjectVal,
		cty.EmptyTupleVal,
		cty.False,
		cty.NegativeInfinity,
		cty.PositiveInfinity,
		cty.True,
		cty.Zero,

		cty.StringVal("hello"),
		cty.StringVal(""),
		cty.NumberIntVal(42),
		cty.NumberIntVal(0),
		cty.NumberFloatVal(3.14),
		cty.NumberFloatVal(0),

		cty.ObjectVal(map[string]cty.Value{
			"foo": cty.StringVal("bar"),
		}),
		cty.ObjectVal(map[string]cty.Value{}),
		cty.ObjectVal(map[string]cty.Value{
			"foo": cty.StringVal("bar"),
			"baz": cty.NumberIntVal(42),
			"bar": cty.True,
		}),
		cty.ObjectVal(map[string]cty.Value{
			"foo": cty.NullVal(cty.String),
			"baz": cty.UnknownVal(cty.Number),
			"bar": cty.DynamicVal,
		}),
		cty.ObjectVal(map[string]cty.Value{
			"foo": cty.ObjectVal(map[string]cty.Value{
				"foo": cty.StringVal("bar"),
				"bar": cty.DynamicVal,
			}),
			"bar": cty.NullVal(cty.Object(map[string]cty.Type{
				"foo": cty.String,
				"bar": cty.DynamicPseudoType,
			})),
		}),

		cty.MapVal(map[string]cty.Value{
			"foo": cty.StringVal("bar"),
			"baz": cty.StringVal("foo"),
			"bar": cty.StringVal("baz"),
		}),
		cty.MapVal(map[string]cty.Value{
			"foo": cty.StringVal("bar"),
			"baz": cty.NullVal(cty.String),
			"bar": cty.UnknownVal(cty.String),
		}),
		cty.MapValEmpty(cty.Number),
		cty.MapValEmpty(cty.Map(cty.Tuple([]cty.Type{cty.Bool, cty.String}))),
		cty.MapValEmpty(cty.DynamicPseudoType),

		cty.ListVal([]cty.Value{
			cty.StringVal("foo"),
			cty.StringVal("foo"),
			cty.StringVal("bar"),
			cty.StringVal("bar"),
		}),
		cty.ListValEmpty(cty.String),
		cty.ListVal([]cty.Value{
			cty.NumberIntVal(42),
			cty.NumberIntVal(43),
		}),
		cty.ListValEmpty(cty.Number),
		cty.ListVal([]cty.Value{
			cty.NumberIntVal(42),
			cty.NumberIntVal(42),
			cty.NullVal(cty.Number),
			cty.UnknownVal(cty.Number),
		}),

		cty.SetVal([]cty.Value{
			cty.StringVal("foo"),
			cty.StringVal("foo"),
			cty.StringVal("bar"),
			cty.StringVal("bar"),
		}),
		cty.SetValEmpty(cty.String),
		cty.SetVal([]cty.Value{
			cty.NumberIntVal(42),
			cty.NumberIntVal(43),
		}),
		cty.SetValEmpty(cty.Number),
		cty.SetVal([]cty.Value{
			cty.NumberIntVal(42),
			cty.NumberIntVal(42),
			cty.NullVal(cty.Number),
			cty.UnknownVal(cty.Number),
		}),

		cty.TupleVal([]cty.Value{
			cty.StringVal("foo"),
			cty.NumberIntVal(42),
			cty.True,
			cty.DynamicVal,
		}),
		cty.TupleVal([]cty.Value{
			cty.True,
			cty.NullVal(cty.String),
			cty.UnknownVal(cty.Number),
		}),
		cty.TupleVal([]cty.Value{}),
	}
	for i, val := range vals {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			t.Logf("Testing value roundtrip: %s", val.GoString())
			roundTrip(t, val)
			t.Logf("Testing UnknownVal roundtrip")
			roundTrip(t, cty.UnknownVal(val.Type()))
			t.Logf("Testing NullVal roundtrip")
			roundTrip(t, cty.NullVal(val.Type()))
		})
	}
}

func TestCtyCodecNil(t *testing.T) {
	var err error
	t.Logf("Testing Nil roundtrip")
	encVal, diag := encodeCtyValue(cty.NilVal)
	if diag.HasErrors() {
		t.Fatalf("Failed to encode value: %s", diag)
	}
	decVal, err := decodeCtyValue(encVal)
	if err != nil {
		t.Fatalf("Failed to decode value: %s", err)
	}
	if decVal != cty.NilVal {
		t.Fatalf("Value roundtrip failed: got %s", decVal.GoString())
	}

	t.Logf("Testing NilType roundtrip")
	encTy, diag := encodeCtyType(cty.NilType)
	if diag.HasErrors() {
		t.Fatalf("Failed to encode type: %s", diag)
	}
	decTy, err := decodeCtyType(encTy)
	if err != nil {
		t.Fatalf("Failed to decode type: %s", err)
	}
	if decTy != cty.NilType {
		t.Fatalf("Type roundtrip failed: got %s", decTy.GoString())
	}
	return
}

func TestCtyCodecObjectOptional(t *testing.T) {
	ty := cty.ObjectWithOptionalAttrs(
		map[string]cty.Type{
			"required": cty.String,
			"optional": cty.String,
		},
		[]string{"optional"},
	)
	assert.True(t, ty.AttributeOptional("optional"))

	decTy := roundTripType(t, ty)
	assert.True(t, decTy.AttributeOptional("optional"))
}

func TestCtyCodecCapsule(t *testing.T) {
	var err error
	t.Logf("Testing Capsule roundtrip")

	data := plugin.Data(plugin.MapData{
		"foo": plugin.StringData("bar"),
		"bar": plugin.NumberData(42),
		"baz": plugin.BoolData(true),
		"qux": plugin.ListData{
			plugin.StringData("foo"),
			plugin.StringData("bar"),
		},
	})

	ctyVal := plugin.EncapsulatedData.ToCty(&data)
	encVal, diag := encodeCtyValue(ctyVal)
	if diag.HasErrors() {
		t.Fatalf("Failed to encode value: %s", diag)
	}
	decVal, err := decodeCtyValue(encVal)
	if err != nil {
		t.Fatalf("Failed to decode value: %s", err)
	}

	decData := plugin.EncapsulatedData.MustFromCty(decVal)
	assert.Equal(t, &data, decData)

	roundTripType(t, plugin.EncapsulatedData.CtyType())

	// t.Logf("Testing NilType roundtrip")
	// encTy, diag := encodeCtyType(cty.NilType)
	// if diag.HasErrors() {
	// 	t.Fatalf("Failed to encode type: %s", diag)
	// }
	// decTy, err := decodeCtyType(encTy)
	// if err != nil {
	// 	t.Fatalf("Failed to decode type: %s", err)
	// }
	// if decTy != cty.NilType {
	// 	t.Fatalf("Type roundtrip failed: got %s", decTy.GoString())
	// }
	// return
}
