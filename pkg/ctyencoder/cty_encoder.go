package ctyencoder

import (
	"github.com/zclconf/go-cty/cty"

	"github.com/blackstork-io/fabric/pkg/diagnostics"
	"github.com/blackstork-io/fabric/plugin/plugindata"
)

// Encoder is a generic interface for encoding cty values into a specific type.
type Encoder[T any] struct {
	// Most general encoder, used when no other encoder is applicable.
	// This includes primitives, unknown types and nulls.
	EncodeVal        func(val cty.Value) (T, diagnostics.Diag)
	EncodePluginData func(val plugindata.Data) (T, diagnostics.Diag)
	MapEncoder       func(val cty.Value) (CollectionEncoder[T], diagnostics.Diag)
	ObjectEncoder    func(val cty.Value) (CollectionEncoder[T], diagnostics.Diag)
	ListEncoder      func(val cty.Value) (CollectionEncoder[T], diagnostics.Diag)
	TupleEncoder     func(val cty.Value) (CollectionEncoder[T], diagnostics.Diag)
	SetEncoder       func(val cty.Value) (CollectionEncoder[T], diagnostics.Diag)
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
		enc, diag = e.ObjectEncoder(val)
	case ty.IsMapType():
		enc, diag = e.MapEncoder(val)
	case ty.IsListType():
		enc, diag = e.ListEncoder(val)
	case ty.IsTupleType():
		enc, diag = e.TupleEncoder(val)
	case ty.IsSetType():
		enc, diag = e.SetEncoder(val)
	}
	addPathIfMissing(diag, path)
	if diags.Extend(diag) {
		return
	}
	if enc != nil {
		if !val.IsNull() && val.IsKnown() {
			path = append(path, nil)
			for it := val.ElementIterator(); it.Next(); {
				k, v := it.Element()
				path[len(path)-1] = valToStep(k, isObj)
				var res T
				res, diag = e.Encode(path, v)
				if !diag.HasErrors() {
					diag.Extend(enc.Add(k, res))
				}
				addPathIfMissing(diag, path)
				diags.Extend(diag)
			}
			path = path[:len(path)-1]
		}
		result, diag = enc.Encode()
	} else if plugindata.EncapsulatedData.ValDecodable(val) {
		data := *plugindata.EncapsulatedData.MustFromCty(val)
		result, diag = e.EncodePluginData(data)
	} else {
		result, diag = e.EncodeVal(val)
	}
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
