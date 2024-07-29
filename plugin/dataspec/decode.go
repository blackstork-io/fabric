package dataspec

import (
	"context"
	"fmt"
	"maps"

	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/ext/customdecode"
	"github.com/hashicorp/hcl/v2/hclsyntax"
	"github.com/zclconf/go-cty/cty"
	"github.com/zclconf/go-cty/cty/convert"

	"github.com/blackstork-io/fabric/cmd/fabctx"
	"github.com/blackstork-io/fabric/pkg/diagnostics"
	"github.com/blackstork-io/fabric/pkg/utils"
	"github.com/blackstork-io/fabric/plugin/dataspec/constraint"
	"github.com/blackstork-io/fabric/plugin/dataspec/deferred"
	"github.com/blackstork-io/fabric/plugin/plugindata"
)

func EvaluateDeferred(ctx context.Context, dataCtx plugindata.Map, val cty.Value) (cty.Value, diagnostics.Diag) {
	return (&transformer{
		ctx:     ctx,
		dataCtx: dataCtx,
	}).transform(val)
}

type transformer struct {
	ctx     context.Context
	dataCtx plugindata.Map
	path    cty.Path
}

func (t *transformer) transform(val cty.Value) (_ cty.Value, diags diagnostics.Diag) {
	ty := val.Type()
	var marks cty.ValueMarks
	val, marks = val.Unmark()
	switch {
	case deferred.Type.CtyTypeEqual(ty):
		// TODO: warn in these cases
		switch {
		case val.IsNull():
			val = cty.NullVal(cty.DynamicPseudoType)
		case !val.IsKnown():
			val = cty.UnknownVal(cty.DynamicPseudoType)
		default:
			eval := deferred.Type.MustFromCty(val)
			val, diags = eval.Eval(t.ctx, t.dataCtx)
			diags.Refine(diagnostics.AddPath(t.path))
		}
	case val.IsNull() || !val.IsKnown():
		// do nothing
	case ty.IsListType() || ty.IsSetType() || ty.IsTupleType():
		l := val.LengthInt()
		if l == 0 {
			break
		}
		vals := make([]cty.Value, 0, l)
		t.path = append(t.path, nil)
		for it := val.ElementIterator(); it.Next(); {
			key, subVal := it.Element()
			t.path[len(t.path)-1] = cty.IndexStep{
				Key: key,
			}
			subVal, diag := t.transform(subVal)
			diags.Extend(diag)
			vals = append(vals, subVal)
		}
		t.path = t.path[:len(t.path)-1]
		if diags.HasErrors() {
			break
		}
		switch {
		case ty.IsListType() && cty.CanListVal(vals):
			val = cty.ListVal(vals)
		case ty.IsSetType() && cty.CanSetVal(vals):
			val = cty.SetVal(vals)
		default:
			// Lists are replaced by tuples if no longer homogenious.
			// Sets are attempted to be converted to set type.
			val = cty.TupleVal(vals)
			if ty.IsSetType() {
				var err error
				val, err = convert.Convert(val, ty)
				if diags.AppendErr(err, "Failed to convert to set") {
					diags.Refine(diagnostics.AddPath(t.path))
				}
			}
		}
	case ty.IsMapType():
		l := val.LengthInt()
		if l == 0 {
			break
		}
		elems := make(map[string]cty.Value, l)
		t.path = append(t.path, nil)
		for it := val.ElementIterator(); it.Next(); {
			key, subVal := it.Element()
			t.path[len(t.path)-1] = cty.IndexStep{
				Key: key,
			}
			newVal, diag := t.transform(subVal)
			diags.Extend(diag)
			elems[key.AsString()] = newVal
		}
		t.path = t.path[:len(t.path)-1]
		if diags.HasErrors() {
			break
		}
		if cty.CanMapVal(elems) {
			val = cty.MapVal(elems)
		} else {
			val = cty.ObjectVal(elems)
		}
	case ty.IsObjectType():
		if ty.Equals(cty.EmptyObject) {
			break
		}
		atys := ty.AttributeTypes()
		elems := make(map[string]cty.Value, len(atys))
		t.path = append(t.path, nil)
		for name := range atys {
			subVal := val.GetAttr(name)
			t.path[len(t.path)-1] = cty.GetAttrStep{
				Name: name,
			}
			subVal, diag := t.transform(subVal)
			diags.Extend(diag)
			elems[name] = subVal
		}
		t.path = t.path[:len(t.path)-1]
		val = cty.ObjectVal(elems)
	}
	if diags.HasErrors() {
		val = cty.DynamicVal
	}
	return val.WithMarks(marks), diags
}

// Decodes hclsyntax.Block into a Block according to the given RootSpec.
// Basic validation is performed on the keys, values of attributes are not fully defined until deferred evaluation,
// so they are not type-checked in the resulting block.
func DecodeBlock(ctx context.Context, block *hclsyntax.Block, rootSpec *RootSpec) (res *Block, diags diagnostics.Diag) {
	return decodeBlock(block, rootSpec.BlockSpec(), fabctx.GetEvalContext(ctx))
}

// Decodes hclsyntax.Block into a Block according to the given RootSpec.
// Deferred evaluation is not allowed
func DecodeAndEvalBlock(ctx context.Context, block *hclsyntax.Block, rootSpec *RootSpec) (res *Block, diags diagnostics.Diag) {
	res, diags = decodeBlock(block, rootSpec.BlockSpec(), fabctx.GetEvalContext(ctx))
	if diags.HasErrors() {
		res = nil
		return
	}
	diags.Extend(EvalBlock(ctx, res, nil))
	if diags.HasErrors() {
		res = nil
		return
	}
	return
}

func EvalBlock(ctx context.Context, block *Block, dataCtx plugindata.Map) (diags diagnostics.Diag) {
	if block == nil {
		return
	}
	for _, block := range block.Blocks {
		diags.Extend(EvalBlock(ctx, block, dataCtx))
	}
	for _, attr := range block.Attrs {
		// deferred eval transform
		var diag diagnostics.Diag
		attr.Value, diag = EvaluateDeferred(ctx, dataCtx, attr.Value)
		if diags.Extend(diag.Refine(diagnostics.DefaultSubject(attr.ValueRange))) {
			continue
		}
		if attr.spec == nil {
			continue // can't convert
		}
		// convert

		var err error
		attr.Value, err = convert.Convert(attr.Value, attr.spec.Type)
		if err != nil {
			diags.Extend(diagnostics.FromErr(
				err,
				diagnostics.DefaultSummary("Incorrect attribute value type"),
				diagnostics.DefaultSubject(attr.ValueRange),
			))
			continue
		}
		if attr.spec.Constraints.Is(constraint.TrimSpace) {
			attr.Value = trimSpace(attr.Value)
		}
		diags.Extend(attr.spec.ValidateValue(attr.Value).Refine(diagnostics.DefaultSubject(attr.ValueRange)))
	}
	return
}

func decodeBlock(block *hclsyntax.Block, blockSpec *BlockSpec, ctx *hcl.EvalContext) (res *Block, diags diagnostics.Diag) {
	if block == nil {
		if blockSpec == nil {
			diags.Append(&hcl.Diagnostic{
				Severity: hcl.DiagError,
				Summary:  "Missing block and blockspec",
				Detail:   "Block and blockspec are both nil, can't decode a block",
			})
			return
		}
		if blockSpec.Required {
			diags.Append(&hcl.Diagnostic{
				Severity: hcl.DiagError,
				Summary:  "Missing block",
				Detail:   fmt.Sprintf("Block of type %s is required", formatHeader(blockSpec.Header.AsDocLabels())),
			})
			return
		}
		name, labels := blockSpec.Header.AsDocLabels()
		labels = append(labels, "<default values>")
		block = &hclsyntax.Block{
			Type:   name,
			Labels: labels,
			Body:   &hclsyntax.Body{},
		}
	}
	res = &Block{
		Header:        make([]string, 1, len(block.Labels)+1),
		HeaderRanges:  make([]hcl.Range, 1, len(block.Labels)+1),
		Attrs:         make(Attributes, len(block.Body.Attributes)),
		Blocks:        make(Blocks, 0, len(block.Body.Blocks)),
		ContentsRange: block.Body.SrcRange,
	}
	res.Header[0] = block.Type
	res.HeaderRanges[0] = block.TypeRange
	res.Header = append(res.Header, block.Labels...)
	res.HeaderRanges = append(res.HeaderRanges, block.LabelRanges...)

	if blockSpec == nil {
		for _, blk := range block.Body.Blocks {
			b, diag := decodeBlock(blk, nil, ctx)
			res.Blocks = append(res.Blocks, b)
			diags.Extend(diag)
		}
		for k, attr := range block.Body.Attributes {
			val, diag := attr.Expr.Value(ctx)
			diags.Extend(diag)
			res.Attrs[k] = &Attr{
				Name:       attr.Name,
				NameRange:  attr.NameRange,
				Value:      val,
				ValueRange: attr.Expr.Range(),
			}
		}
		return
	}

	wasUsed := make([]bool, len(blockSpec.Blocks))
nextBlock:
	for _, subB := range block.Body.Blocks {
		for i, bSpec := range blockSpec.Blocks {
			if !bSpec.Repeatable && wasUsed[i] {
				continue
			}
			if !bSpec.Header.Match(subB.Type, subB.Labels) {
				continue
			}
			wasUsed[i] = true
			b, diag := decodeBlock(subB, bSpec, ctx)
			diags.Extend(diag)
			res.Blocks = append(res.Blocks, b)
			continue nextBlock
		}
		if !blockSpec.AllowUnspecifiedBlocks {
			diags.Append(&hcl.Diagnostic{
				Severity: hcl.DiagError,
				Summary:  "Unexpected block",
				Detail:   fmt.Sprintf("%s can not contain this block", formatHeader(block.Type, block.Labels)),
				Subject:  &subB.TypeRange,
				Context:  &block.Body.SrcRange,
			})
			continue nextBlock
		}
		b, diag := decodeBlock(subB, nil, ctx)
		res.Blocks = append(res.Blocks, b)
		diags.Extend(diag)
	}
	for i, bSpec := range blockSpec.Blocks {
		if bSpec.Required && !wasUsed[i] {
			diags.Append(&hcl.Diagnostic{
				Severity: hcl.DiagError,
				Summary:  "Missing block",
				Detail:   fmt.Sprintf("%s requires a block of type %s", formatHeader(block.Type, block.Labels), formatHeader(bSpec.Header.AsDocLabels())),
				Subject:  block.Body.MissingItemRange().Ptr(),
			})
		}
	}

	attrs := maps.Clone(block.Body.Attributes)
	for _, spec := range blockSpec.Attrs {
		attr, found := utils.Pop(attrs, spec.Name)
		if !found {
			if spec.Constraints.Is(constraint.Required) {
				diags.Append(&hcl.Diagnostic{
					Severity: hcl.DiagError,
					Summary:  "Missing required attribute",
					Detail:   fmt.Sprintf("The attribute %q is required, but no definition was found.", spec.Name),
					Subject:  block.Body.MissingItemRange().Ptr(),
				})
			} else if spec.DefaultVal != cty.NilVal {
				rng := block.DefRange()

				res.Attrs[spec.Name] = &Attr{
					Name:       spec.Name,
					NameRange:  rng,
					Value:      spec.DefaultVal,
					ValueRange: rng,
					Secret:     spec.Secret,
					spec:       spec,
				}
			}
			continue
		}
		if spec.Deprecated != "" {
			diags = append(diags, &hcl.Diagnostic{
				Severity: hcl.DiagWarning,
				Summary:  "Deprecated attribute",
				Detail:   fmt.Sprintf("The attribute %q is deprecated: %s", spec.Name, spec.Deprecated),
				Subject:  &attr.NameRange,
				Context:  &attr.SrcRange,
			})
		}
		val, diag := DecodeAttr(spec, attr, ctx)

		if diags.Extend(diag.Refine(diagnostics.DefaultSubject(attr.Expr.Range()))) {
			continue
		}

		res.Attrs[spec.Name] = &Attr{
			Name:       spec.Name,
			NameRange:  attr.NameRange,
			Value:      val,
			ValueRange: attr.Expr.Range(),
			Secret:     spec.Secret,
			spec:       spec,
		}
	}
	for name, attr := range attrs {
		if !blockSpec.AllowUnspecifiedAttributes {
			diags.Append(&hcl.Diagnostic{
				Severity: hcl.DiagError,
				Summary:  "Unsupported attribute",
				Detail: fmt.Sprintf(
					"Unsupported attribute %q",
					name,
				),
				Subject: &attr.NameRange,
				Context: hcl.RangeBetween(attr.NameRange, attr.Expr.Range()).Ptr(),
			})
			continue
		}
		val, diag := attr.Expr.Value(ctx)
		if !diags.Extend(diag) {
			res.Attrs[name] = &Attr{
				Name:       name,
				NameRange:  attr.NameRange,
				Value:      val,
				ValueRange: attr.Expr.Range(),
			}
		}
	}
	return
}

func DecodeAttr(spec *AttrSpec, attr *hclsyntax.Attribute, ctx *hcl.EvalContext) (val cty.Value, diags diagnostics.Diag) {
	var decodeFn customdecode.CustomExpressionDecoderFunc
	if spec != nil {
		decodeFn = customdecode.CustomExpressionDecoderForType(spec.Type)
	}
	var diag hcl.Diagnostics
	if decodeFn != nil {
		val, diag = decodeFn(attr.Expr, ctx)
	} else {
		val, diag = attr.Expr.Value(ctx)
	}

	return val, diagnostics.Diag(diag).Refine(diagnostics.DefaultSubject(attr.Expr.Range()))
}
