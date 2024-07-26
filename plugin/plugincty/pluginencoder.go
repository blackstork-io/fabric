package plugincty

import (
	"log/slog"

	"github.com/hashicorp/hcl/v2"
	"github.com/zclconf/go-cty/cty"
	"github.com/zclconf/go-cty/cty/convert"

	"github.com/blackstork-io/fabric/pkg/ctyencoder"
	"github.com/blackstork-io/fabric/pkg/diagnostics"
	"github.com/blackstork-io/fabric/plugin/plugindata"
)

type nullEncoder struct{}

// Add implements CollectionEncoder.
func (n nullEncoder) Add(k cty.Value, v plugindata.Data) diagnostics.Diag {
	return nil
}

// Encode implements CollectionEncoder.
func (n nullEncoder) Encode() (plugindata.Data, diagnostics.Diag) {
	return nil, nil
}

var _ ctyencoder.CollectionEncoder[plugindata.Data] = nullEncoder{}

type mapEncoder plugindata.Map

func newMapEncoder(val cty.Value) (ctyencoder.CollectionEncoder[plugindata.Data], diagnostics.Diag) {
	if val.IsNull() || !val.IsKnown() {
		return nullEncoder{}, nil
	}
	return mapEncoder(make(plugindata.Map, val.LengthInt())), nil
}

var _ ctyencoder.CollectionEncoder[plugindata.Data] = mapEncoder{}

func (m mapEncoder) Encode() (plugindata.Data, diagnostics.Diag) {
	return plugindata.Map(m), nil
}

func (m mapEncoder) Add(k cty.Value, v plugindata.Data) diagnostics.Diag {
	if k.IsNull() || !k.Type().Equals(cty.String) {
		return diagnostics.Diag{&hcl.Diagnostic{
			Severity: hcl.DiagError,
			Summary:  "Incorrect key type",
			Detail:   k.Type().FriendlyName() + " is not supported here, only strings are allowed",
		}}
	}
	m[k.AsString()] = v
	return nil
}

type listEncoder plugindata.List

func newListEncoder(val cty.Value) (ctyencoder.CollectionEncoder[plugindata.Data], diagnostics.Diag) {
	if val.IsNull() || !val.IsKnown() {
		return nullEncoder{}, nil
	}
	l := listEncoder(make(plugindata.List, 0, val.LengthInt()))
	return &l, nil
}

var _ ctyencoder.CollectionEncoder[plugindata.Data] = mapEncoder{}

func (l *listEncoder) Encode() (plugindata.Data, diagnostics.Diag) {
	return plugindata.List(*l), nil
}

func (l *listEncoder) Add(k cty.Value, v plugindata.Data) diagnostics.Diag {
	*l = append(*l, v)
	return nil
}

var pluginDataEncoder = ctyencoder.Encoder[plugindata.Data]{
	EncodeVal: func(val cty.Value) (result plugindata.Data, diags diagnostics.Diag) {
		if val.IsNull() || !val.IsKnown() {
			return nil, nil
		}
		ty := val.Type()
		if ty.IsPrimitiveType() {
			switch ty {
			case cty.String:
				result = plugindata.String(val.AsString())
				return
			case cty.Number:
				n, _ := val.AsBigFloat().Float64()
				result = plugindata.Number(n)
				return
			case cty.Bool:
				result = plugindata.Bool(val.True())
				return
			}
		}
		if plugindata.EncapsulatedData.ValCtyTypeEqual(val) {
			result = *plugindata.EncapsulatedData.MustFromCty(val)
			return
		} else {
			res, err := convert.Convert(val, plugindata.EncapsulatedData.CtyType())
			if diags.AppendErr(err, "Failed to encode data") {
				slog.Error("Failed to encode", "in", val.GoString())
				return
			}
			result = *plugindata.EncapsulatedData.MustFromCty(res)
			return
		}
	},
	EncodePluginData: func(val plugindata.Data) (result plugindata.Data, diags diagnostics.Diag) {
		return val, nil
	},
	MapEncoder:    newMapEncoder,
	ObjectEncoder: newMapEncoder,
	ListEncoder:   newListEncoder,
	TupleEncoder:  newListEncoder,
	SetEncoder:    newListEncoder,
}

func Encode(val cty.Value) (plugindata.Data, diagnostics.Diag) {
	return pluginDataEncoder.Encode(nil, val)
}
