package dataspec

import (
	"fmt"

	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/hcldec"
	"github.com/hashicorp/hcl/v2/hclsyntax"
	"github.com/zclconf/go-cty/cty"

	"github.com/blackstork-io/fabric/pkg/diagnostics"
)

// Wrapper over hcldec.Decode.
func Decode(body hcl.Body, spec RootSpec, ctx *hcl.EvalContext) (val cty.Value, diags diagnostics.Diag) {
	if _, ok := spec.(*ObjDumpSpec); ok {
		b, ok := body.(*hclsyntax.Body)
		if !ok {
			diags.Append(
				&hcl.Diagnostic{
					Severity: hcl.DiagError,
					Summary:  "This type of hcl.Body is not supported",
					Detail:   "Only works on native hcl format implementation",
					Subject:  body.MissingItemRange().Ptr(),
				},
			)
			return
		}
		return hclBodyToVal(b, ctx)
	}
	v, diag := hcldec.Decode(body, spec.HcldecSpec(), ctx)
	if diags.Extend(diag) {
		return
	}
	val = v
	return
}

type source struct {
	rng     *hcl.Range
	isBlock bool
}

func (s *source) kind() string {
	if s.isBlock {
		return "block"
	}
	return "attribute"
}

func hclBodyToVal(body *hclsyntax.Body, ctx *hcl.EvalContext) (val cty.Value, diags diagnostics.Diag) {
	obj := make(map[string]cty.Value, len(body.Attributes)+len(body.Blocks))
	sources := make(map[string]source, len(body.Attributes)+len(body.Blocks))

	for name, attr := range body.Attributes {
		attrVal, diag := attr.Expr.Value(ctx)
		if diags.Extend(diag) {
			continue
		}
		obj[name] = attrVal
		sources[name] = source{
			rng:     &attr.NameRange,
			isBlock: false,
		}
	}
	for _, block := range body.Blocks {
		if len(block.Labels) != 0 {
			diags.Append(&hcl.Diagnostic{
				Severity: hcl.DiagWarning,
				Summary:  "Block labels are ignored in this context",
				Detail:   "",
				Subject:  hcl.RangeBetween(block.LabelRanges[0], block.LabelRanges[len(block.LabelRanges)-1]).Ptr(),
			})
		}
		blockVal, diag := hclBodyToVal(block.Body, ctx)
		if src, found := sources[block.Type]; found {
			diag.Append(&hcl.Diagnostic{
				Severity: hcl.DiagError,
				Summary:  "Name conflict",
				Detail: fmt.Sprintf(
					"This block uses the same name as the %s at %s:%d:%d",
					src.kind(), src.rng.Filename, src.rng.Start.Line, src.rng.Start.Column,
				),
				Subject: &block.TypeRange,
			})
		} else {
			sources[block.Type] = source{
				rng:     &block.TypeRange,
				isBlock: true,
			}
		}
		if diags.Extend(diag) {
			continue
		}
		obj[block.Type] = blockVal
	}
	val = cty.ObjectVal(obj)
	return
}
