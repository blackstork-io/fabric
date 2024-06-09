package plugincty

import (
	"fmt"
	"log/slog"

	"github.com/hashicorp/hcl/v2"
	"github.com/zclconf/go-cty/cty"
	"github.com/zclconf/go-cty/cty/convert"

	"github.com/blackstork-io/fabric/pkg/ctyencoder"
	"github.com/blackstork-io/fabric/pkg/diagnostics"
	"github.com/blackstork-io/fabric/plugin"
)

type nullEncoder struct{}

// Add implements CollectionEncoder.
func (n nullEncoder) Add(k cty.Value, v plugin.Data) diagnostics.Diag {
	return nil
}

// Encode implements CollectionEncoder.
func (n nullEncoder) Encode() (plugin.Data, diagnostics.Diag) {
	return nil, nil
}

var _ ctyencoder.CollectionEncoder[plugin.Data] = nullEncoder{}

type mapEncoder plugin.MapData

func newMapEncoder(val cty.Value) ctyencoder.CollectionEncoder[plugin.Data] {
	if val.IsNull() || !val.IsKnown() {
		return nullEncoder{}
	}
	return mapEncoder(make(plugin.MapData, val.LengthInt()))
}

var _ ctyencoder.CollectionEncoder[plugin.Data] = mapEncoder{}

func (m mapEncoder) Encode() (plugin.Data, diagnostics.Diag) {
	return plugin.MapData(m), nil
}

func (m mapEncoder) Add(k cty.Value, v plugin.Data) diagnostics.Diag {
	if k.IsNull() || !k.Type().Equals(cty.String) {
		return diagnostics.Diag{&hcl.Diagnostic{
			Severity: hcl.DiagError,
			Summary:  "Incorrect key type",
			Detail:   fmt.Sprintf("%s is not supported here, only strings are allowed", k.Type().FriendlyName()),
		}}
	}
	m[k.AsString()] = v
	return nil
}

type listEncoder plugin.ListData

func newListEncoder(val cty.Value) ctyencoder.CollectionEncoder[plugin.Data] {
	if val.IsNull() || !val.IsKnown() {
		return nullEncoder{}
	}
	l := listEncoder(make(plugin.ListData, 0, val.LengthInt()))
	return &l
}

var _ ctyencoder.CollectionEncoder[plugin.Data] = mapEncoder{}

func (l *listEncoder) Encode() (plugin.Data, diagnostics.Diag) {
	return plugin.ListData(*l), nil
}

func (l *listEncoder) Add(k cty.Value, v plugin.Data) diagnostics.Diag {
	*l = append(*l, v)
	return nil
}

var pluginDataEncoder = ctyencoder.Encoder[plugin.Data]{
	EncodeVal: func(val cty.Value) (result plugin.Data, diags diagnostics.Diag) {
		if val.IsNull() || !val.IsKnown() {
			return nil, nil
		}
		ty := val.Type()
		if ty.IsPrimitiveType() {
			switch ty {
			case cty.String:
				result = plugin.StringData(val.AsString())
				return
			case cty.Number:
				n, _ := val.AsBigFloat().Float64()
				result = plugin.NumberData(n)
				return
			case cty.Bool:
				result = plugin.BoolData(val.True())
				return
			}
		}
		if plugin.EncapsulatedData.ValCtyTypeEqual(val) {
			result = *plugin.EncapsulatedData.MustFromCty(val)
			return
		} else {
			slog.Error("convert", "in", val.Type().FriendlyName(), "out", plugin.EncapsulatedData.CtyType().FriendlyName())
			res, err := convert.Convert(val, plugin.EncapsulatedData.CtyType())
			if diags.AppendErr(err, "Failed to encode data") {
				// panic("Failed to encode")
				slog.Error("Failed to encode", "in", val.GoString())

				return
			}
			result = *plugin.EncapsulatedData.MustFromCty(res)
			return
		}
	},
	EncodePluginData: func(val plugin.Data) (result plugin.Data, diags diagnostics.Diag) {
		return val, nil
	},
	MapEncoder:    newMapEncoder,
	ObjectEncoder: newMapEncoder,
	ListEncoder:   newListEncoder,
	TupleEncoder:  newListEncoder,
	SetEncoder:    newListEncoder,
}

func Encode(val cty.Value) (plugin.Data, diagnostics.Diag) {
	return pluginDataEncoder.Encode(nil, val)
}
