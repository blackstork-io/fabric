package dataspec

import (
	"fmt"
	"maps"

	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/ext/customdecode"
	"github.com/hashicorp/hcl/v2/hclsyntax"
	"github.com/zclconf/go-cty/cty"
	"github.com/zclconf/go-cty/cty/convert"

	"github.com/blackstork-io/fabric/pkg/diagnostics"
	"github.com/blackstork-io/fabric/pkg/utils"
	"github.com/blackstork-io/fabric/plugin/dataspec/constraint"
)

func Decode(block *hclsyntax.Block, rootSpec *RootSpec, ctx *hcl.EvalContext) (res *Block, diags diagnostics.Diag) {
	return DecodeBlock(block, rootSpec.BlockSpec(), ctx)
}

func DecodeBlock(block *hclsyntax.Block, blockSpec *BlockSpec, ctx *hcl.EvalContext) (res *Block, diags diagnostics.Diag) {
	if block == nil {
		if blockSpec != nil && blockSpec.Required {
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
			b, diag := DecodeBlock(blk, nil, ctx)
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

	// hcldec.AttrSpec

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
			b, diag := DecodeBlock(subB, bSpec, ctx)
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
		b, diag := DecodeBlock(subB, nil, ctx)
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
		val, diag := decodeAttr(spec, attr, ctx)
		diag.DefaultSubject(attr.Expr.Range().Ptr())
		if diags.Extend(diag) {
			continue
		}

		res.Attrs[spec.Name] = &Attr{
			Name:       spec.Name,
			NameRange:  attr.NameRange,
			Value:      val,
			ValueRange: attr.Expr.Range(),
			Secret:     spec.Secret,
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

func decodeAttr(spec *AttrSpec, attr *hclsyntax.Attribute, ctx *hcl.EvalContext) (val cty.Value, diags diagnostics.Diag) {
	var diag hcl.Diagnostics
	if decodeFn := customdecode.CustomExpressionDecoderForType(spec.Type); decodeFn != nil {
		val, diag = decodeFn(attr.Expr, ctx)
		diags = diagnostics.Diag(diag)
		return
	}

	val, diag = attr.Expr.Value(ctx)
	if diags.Extend(diag) {
		return
	}

	var err error
	val, err = convert.Convert(val, spec.Type)
	if err != nil {
		diags.Append(&hcl.Diagnostic{
			Severity: hcl.DiagError,
			Summary:  "Incorrect attribute value type",
			Detail: fmt.Sprintf(
				"Inappropriate value for attribute %q: %s.",
				spec.Name, err.Error(),
			),
			Subject:     attr.Expr.Range().Ptr(),
			Context:     hcl.RangeBetween(attr.NameRange, attr.Expr.Range()).Ptr(),
			Expression:  attr.Expr,
			EvalContext: ctx,
		})
		return
	}
	if spec.Constraints.Is(constraint.TrimSpace) {
		val = trimSpace(val)
	}
	diags.Extend(spec.ValidateValue(val))
	return
}
