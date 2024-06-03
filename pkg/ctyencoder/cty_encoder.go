package ctyencoder

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/hashicorp/hcl/v2"
	"github.com/zclconf/go-cty/cty"

	"github.com/blackstork-io/fabric/parser/definitions"
	"github.com/blackstork-io/fabric/pkg/diagnostics"
	"github.com/blackstork-io/fabric/plugin"
)

// Encoder is a generic interface for encoding cty values into a specific type.
type Encoder[T any] struct {
	// Most general encoder, used when no other encoder is applicable.
	// This includes primitives, unknown types and nulls.
	Encode           func(val cty.Value) (T, diagnostics.Diag)
	EncodePluginData func(val plugin.Data) (T, diagnostics.Diag)
	MapEncoder       func(val cty.Value) CollectionEncoder[T]
	ObjectEncoder    func(val cty.Value) CollectionEncoder[T]
	ListEncoder      func(val cty.Value) CollectionEncoder[T]
	TupleEncoder     func(val cty.Value) CollectionEncoder[T]
	SetEncoder       func(val cty.Value) CollectionEncoder[T]
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

// EvalEncoder is a generic interface for encoding cty values into a specific type.
// Can evaluate expressions that require data context (if provided).
type EvalEncoder[T any] struct {
	Encoder Encoder[T]
	Ctx     context.Context
	DataCtx plugin.MapData
}

func (e *EvalEncoder[T]) Encode(val cty.Value) (result T, diags diagnostics.Diag) {
	return e.encode(nil, val)
}

func (e *EvalEncoder[T]) encode(path cty.Path, val cty.Value) (result T, diags diagnostics.Diag) {
	const maxNestedEvals = 100 // arbitrary limit to prevent infinite recursion
	var diag diagnostics.Diag

	// var curEnc StructureEncoder[T]
	for nestedEvals := 0; ; nestedEvals++ {
		ty := val.Type()
		if ty != cty.NilType {
			var enc CollectionEncoder[T]
			isObj := false
			switch {
			case ty.IsObjectType():
				isObj = true
				enc = e.Encoder.ObjectEncoder(val)
			case ty.IsMapType():
				enc = e.Encoder.MapEncoder(val)
			case ty.IsListType():
				enc = e.Encoder.ListEncoder(val)
			case ty.IsTupleType():
				enc = e.Encoder.TupleEncoder(val)
			case ty.IsSetType():
				enc = e.Encoder.SetEncoder(val)
			}
			if enc != nil {
				if !val.IsNull() && val.IsKnown() {
					path = append(path, nil)
					for it := val.ElementIterator(); it.Next(); {
						k, v := it.Element()
						path[len(path)-1] = valToStep(k, isObj)
						var res T
						res, diag = e.encode(path, v)
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
				t := plugin.EncapsulatedData.MustFromCty(val)
				slog.Error("plugin.EncapsulatedData", "v", (*t), "path", path)
				data := *plugin.EncapsulatedData.MustFromCty(val)
				result, diag = e.Encoder.EncodePluginData(data)
				addPathIfMissing(diag, path)
				diags.Extend(diag)
				return
			}
			if definitions.DataCtxEvalNeededType.ValDecodable(val) {
				toEval := definitions.DataCtxEvalNeededType.MustFromCty(val)
				if e.DataCtx == nil {
					diags.Append(&hcl.Diagnostic{
						Severity: hcl.DiagError,
						Summary:  "Incorrect data type",
						Detail:   "Attempted to evaluate expression that requires data context without it being present. You can't use this expression in this place.",
						Subject:  toEval.Range(),
						Extra:    diagnostics.NewPathExtra(path),
					})
					return
				}
				if nestedEvals >= maxNestedEvals {
					diags.Append(&hcl.Diagnostic{
						Severity: hcl.DiagError,
						Summary:  "Nesting limit reached",
						Detail: fmt.Sprintf(
							"%s required unwrapping more than %d times; probably an infinite recursion",
							ty.FriendlyName(),
							maxNestedEvals,
						),
						Subject: toEval.Range(),
						Extra:   diagnostics.NewPathExtra(path),
					})
					return
				}
				val, diag = toEval.Eval(e.Ctx, e.DataCtx)
				addPathIfMissing(diag, path)
				diags.Extend(diag)
				continue
			}
		}
		result, diag = e.Encoder.Encode(val)
		addPathIfMissing(diag, path)
		diags.Extend(diag)
		return
	}
}
