package ctyencoder

import (
	"context"
	"fmt"

	"github.com/hashicorp/hcl/v2"
	"github.com/zclconf/go-cty/cty"

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

var _ CollectionEncoder[plugin.Data] = nullEncoder{}

type mapEncoder plugin.MapData

func newMapEncoder(val cty.Value) CollectionEncoder[plugin.Data] {
	if val.IsNull() || !val.IsKnown() {
		return nullEncoder{}
	}
	return mapEncoder(make(plugin.MapData, val.LengthInt()))
}

var _ CollectionEncoder[plugin.Data] = mapEncoder{}

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

func newListEncoder(val cty.Value) CollectionEncoder[plugin.Data] {
	if val.IsNull() || !val.IsKnown() {
		return nullEncoder{}
	}
	l := listEncoder(make(plugin.ListData, 0, val.LengthInt()))
	return &l
}

var _ CollectionEncoder[plugin.Data] = mapEncoder{}

func (l *listEncoder) Encode() (plugin.Data, diagnostics.Diag) {
	return plugin.ListData(*l), nil
}

func (l *listEncoder) Add(k cty.Value, v plugin.Data) diagnostics.Diag {
	*l = append(*l, v)
	return nil
}

var pluginDataEncoder = Encoder[plugin.Data]{
	Encode: func(val cty.Value) (result plugin.Data, diags diagnostics.Diag) {
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
		diags.Append(&hcl.Diagnostic{
			Severity: hcl.DiagError,
			Summary:  "Incorrect data type",
			Detail:   fmt.Sprintf("%s is not supported here", ty.FriendlyName()),
		})
		return
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

// ToPluginData converts cty.Value to plugin.Data
func ToPluginData(ctx context.Context, dataCtx plugin.MapData, val cty.Value) (plugin.Data, diagnostics.Diag) {
	return (&EvalEncoder[plugin.Data]{
		Encoder: pluginDataEncoder,
		Ctx:     ctx,
		DataCtx: dataCtx,
	}).Encode(val)
}
