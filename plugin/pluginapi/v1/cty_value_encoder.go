package pluginapiv1

import (
	"fmt"

	"github.com/zclconf/go-cty/cty"
	"github.com/zclconf/go-cty/cty/convert"

	"github.com/blackstork-io/fabric/eval/dataquery"
	"github.com/blackstork-io/fabric/pkg/ctyencoder"
	"github.com/blackstork-io/fabric/pkg/diagnostics"
	"github.com/blackstork-io/fabric/plugin"
)

var grpcCtyValueDecoder = ctyencoder.Encoder[*CtyValue]{
	EncodeVal: func(val cty.Value) (res *CtyValue, diags diagnostics.Diag) {
		var err error
		if val == cty.NilVal {
			return nil, nil
		}
		ty := val.Type()
		if ty.IsPrimitiveType() {
			res, err = encodeCtyPrimitiveValue(val)
			diags.AppendErr(err, "Can't encode a value")
			return
		}
		if dataquery.DelayedEvalType.CtyTypeEqual(ty) {
			convVal, err := convert.Convert(val, plugin.EncapsulatedData.CtyType())
			if diags.AppendErr(err, "Can't encode a value") {
				return
			}
			res = &CtyValue{
				Type: &CtyType{
					Data: &CtyType_Encapsulated{Encapsulated: CtyCapsuleType_CAPSULE_DELAYED_EVAL},
				},
				Data: &CtyValue_PluginData{
					PluginData: encodeData(*plugin.EncapsulatedData.MustFromCty(convVal)),
				},
			}
			return
		}
		diags.Add("Unsupported value", "unsupported cty value: "+ty.FriendlyName())
		return
	},
	EncodePluginData: func(val plugin.Data) (*CtyValue, diagnostics.Diag) {
		return &CtyValue{
			Type: &CtyType{
				Data: &CtyType_Encapsulated{Encapsulated: CtyCapsuleType_CAPSULE_PLUGIN_DATA},
			},
			Data: &CtyValue_PluginData{
				PluginData: encodeData(val),
			},
		}, nil
	},
	MapEncoder:    newMapEncoder,
	ObjectEncoder: newMapEncoder,
	ListEncoder:   newListEncoder,
	TupleEncoder:  newListEncoder,
	SetEncoder:    newListEncoder,
}

func newMapEncoder(val cty.Value) ctyencoder.CollectionEncoder[*CtyValue] {
	ty, err := encodeCtyType(val.Type())
	me := &mapEncoder{
		ty:  ty,
		err: err,
	}
	if err == nil && !val.IsNull() {
		me.values = make(map[string]*CtyValue, val.LengthInt())
	}
	return me
}

type mapEncoder struct {
	ty     *CtyType
	err    error
	values map[string]*CtyValue
}

// Add implements ctyencoder.CollectionEncoder.
func (me *mapEncoder) Add(k cty.Value, v *CtyValue) diagnostics.Diag {
	if me.values == nil {
		return nil
	}
	me.values[k.AsString()] = v
	return nil
}

// Encode implements ctyencoder.CollectionEncoder.
func (me *mapEncoder) Encode() (*CtyValue, diagnostics.Diag) {
	if me.err != nil {
		return nil, diagnostics.FromErr(me.err, "Can't encode a value")
	}
	val := &CtyValue{
		Type: me.ty,
	}
	if me.values == nil {
		return val, nil
	}
	val.Data = &CtyValue_MapLike{
		MapLike: &CtyMapLike{
			Elements: me.values,
		},
	}
	return val, nil
}

func newListEncoder(val cty.Value) ctyencoder.CollectionEncoder[*CtyValue] {
	ty, err := encodeCtyType(val.Type())
	le := &listEncoder{
		ty:  ty,
		err: err,
	}
	if err == nil && !val.IsNull() {
		le.values = make([]*CtyValue, val.LengthInt())
	}
	return le
}

type listEncoder struct {
	ty     *CtyType
	err    error
	values []*CtyValue
}

func (le *listEncoder) Add(_ cty.Value, v *CtyValue) diagnostics.Diag {
	if le.values == nil {
		return nil
	}
	le.values = append(le.values, v)
	return nil
}

func (le *listEncoder) Encode() (*CtyValue, diagnostics.Diag) {
	if le.err != nil {
		return nil, diagnostics.FromErr(le.err, "Can't encode a value")
	}
	val := &CtyValue{
		Type: le.ty,
	}
	if le.values == nil {
		return val, nil
	}
	val.Data = &CtyValue_ListLike{
		ListLike: &CtyListLike{
			Elements: le.values,
		},
	}
	return val, nil
}

func encodeCtyValue(src cty.Value) (*CtyValue, error) {
	res, diags := grpcCtyValueDecoder.Encode(nil, src)
	if diags.HasErrors() {
		return nil, diags
	}
	return res, nil
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
