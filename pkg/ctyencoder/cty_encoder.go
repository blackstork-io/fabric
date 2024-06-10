package ctyencoder

import (
	"github.com/zclconf/go-cty/cty"

	"github.com/blackstork-io/fabric/pkg/diagnostics"
	"github.com/blackstork-io/fabric/plugin"
)

// Encoder is a generic interface for encoding cty values into a specific type.
type Encoder[T any] struct {
	// Most general encoder, used when no other encoder is applicable.
	// This includes primitives, unknown types and nulls.
	EncodeVal        func(val cty.Value) (T, diagnostics.Diag)
	EncodePluginData func(val plugin.Data) (T, diagnostics.Diag)
	MapEncoder       func(val cty.Value) CollectionEncoder[T]
	ObjectEncoder    func(val cty.Value) CollectionEncoder[T]
	ListEncoder      func(val cty.Value) CollectionEncoder[T]
	TupleEncoder     func(val cty.Value) CollectionEncoder[T]
	SetEncoder       func(val cty.Value) CollectionEncoder[T]
}

func (e *Encoder[T]) Encode(path cty.Path, val cty.Value) (result T, diags diagnostics.Diag) {
	var diag diagnostics.Diag

	ty := val.Type()
	if ty == cty.NilType {
		result, diag = e.EncodeVal(val)
		addPathIfMissing(diag, path)
		diags.Extend(diag)
		return
	}

	var enc CollectionEncoder[T]
	isObj := false
	switch {
	case ty.IsObjectType():
		isObj = true
		enc = e.ObjectEncoder(val)
	case ty.IsMapType():
		enc = e.MapEncoder(val)
	case ty.IsListType():
		enc = e.ListEncoder(val)
	case ty.IsTupleType():
		enc = e.TupleEncoder(val)
	case ty.IsSetType():
		enc = e.SetEncoder(val)
	}
	if enc != nil {
		if !val.IsNull() && val.IsKnown() {
			path = append(path, nil)
			for it := val.ElementIterator(); it.Next(); {
				k, v := it.Element()
				path[len(path)-1] = valToStep(k, isObj)
				var res T
				res, diag = e.Encode(path, v)
				diags.Extend(diag)
				diags.Extend(enc.Add(k, res))
			}
			path = path[:len(path)-1]
		}
		result, diag = enc.Encode()
		addPathIfMissing(diag, path)
		diags.Extend(diag)
		return
	}

	if plugin.EncapsulatedData.ValDecodable(val) {
		data := *plugin.EncapsulatedData.MustFromCty(val)
		result, diag = e.EncodePluginData(data)
		addPathIfMissing(diag, path)
		diags.Extend(diag)
		return
	}

	result, diag = e.EncodeVal(val)
	addPathIfMissing(diag, path)
	diags.Extend(diag)
	return
}

// CollectionEncoder is an interface for encoding cty collections into a specific type.
type CollectionEncoder[T any] interface {
	// Will be called for each element in the collection.
	Add(k cty.Value, v T) diagnostics.Diag
	// Will be called after all elements are added.
	Encode() (T, diagnostics.Diag)
}

func addPathIfMissing(diags diagnostics.Diag, path cty.Path) {
	var extra diagnostics.PathExtra
	for _, diag := range diags {
		_, found := diagnostics.DiagnosticExtra[diagnostics.PathExtra](diag)
		if found {
			continue
		}
		if extra == nil {
			extra = diagnostics.NewPathExtra(path)
		}
		diagnostics.AddExtra(diag, extra)
	}
}

func valToStep(v cty.Value, isObj bool) cty.PathStep {
	if isObj {
		return cty.GetAttrStep{
			Name: v.AsString(),
		}
	}
	return cty.IndexStep{
		Key: v,
	}
}
