package eval

import (
	"context"
	"fmt"
	"maps"

	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/hclsyntax"
	"github.com/zclconf/go-cty/cty"

	"github.com/blackstork-io/fabric/cmd/fabctx"
	"github.com/blackstork-io/fabric/parser"
	"github.com/blackstork-io/fabric/parser/definitions"
	"github.com/blackstork-io/fabric/pkg/diagnostics"
	"github.com/blackstork-io/fabric/pkg/utils"
	"github.com/blackstork-io/fabric/plugin/dataspec"
	"github.com/blackstork-io/fabric/plugin/dataspec/deferred"
	"github.com/blackstork-io/fabric/plugin/plugindata"
)

type Dynamic struct {
	block     *hclsyntax.Block
	condition *dataspec.Attr
	items     *dataspec.Attr
	children  []*Content
}

var dynamicBlockCond = &dataspec.AttrSpec{
	Name: "condition",
	Type: cty.Bool,
	Doc:  "Condition indicating whether dynamic block should be rendered",
}

var dynamicBlockItems = &dataspec.AttrSpec{
	Name: "items",
	Type: plugindata.Encapsulated.CtyType(),
	Doc:  "Items to be iterated over (list or map)",
}

func LoadDynamic(ctx context.Context, providers ContentProviders, node *definitions.ParsedDynamic) (_ *Dynamic, diags diagnostics.Diag) {
	var diag diagnostics.Diag
	block := &Dynamic{
		block:    node.Block,
		children: make([]*Content, 0, len(node.Content)),
	}
	evalCtx := fabctx.GetEvalContext(deferred.WithQueryFuncs(ctx))
	block.condition, diag = dataspec.DecodeAttr(evalCtx, node.Condition, dynamicBlockCond)
	diags.Extend(diag)
	if node.Items != nil {
		block.items, diag = dataspec.DecodeAttr(evalCtx, node.Items, dynamicBlockItems)
		diags.Extend(diag)
	}

	for _, child := range node.Content {
		decoded, diag := LoadContent(ctx, providers, child)
		diags.Extend(diag)
		block.children = append(block.children, decoded)
	}
	return block, diags
}

func addDynamicVars(child *Content, dynVarVals *definitions.ParsedVars) *Content {
	if child.Plugin != nil {
		chAction := *child.Plugin
		chAction.Vars = chAction.Vars.MergeWithBaseVars(dynVarVals)

		return &Content{
			Plugin: &chAction,
		}
	}
	if child.Section != nil {
		chSection := *child.Section
		chSection.vars = chSection.vars.MergeWithBaseVars(dynVarVals)
		chSection.children = utils.FnMap(chSection.children, func(child *Content) *Content {
			return addDynamicVars(child, dynVarVals)
		})
		return &Content{
			Section: &chSection,
		}
	}
	if child.Dynamic != nil {
		if child.Dynamic.items == nil {
			chDynamic := *child.Dynamic
			chDynamic.children = utils.FnMap(chDynamic.children, func(child *Content) *Content {
				return addDynamicVars(child, dynVarVals)
			})
			return &Content{
				Dynamic: &chDynamic,
			}
		}
	}
	return child
}

func unwrapDynamicContent(ctx context.Context, children []*Content, dataCtx plugindata.Map, dynVarVals *definitions.ParsedVars) (res []*Content, diags diagnostics.Diag) {
	ctx = deferred.WithQueryFuncs(ctx)
	res = make([]*Content, 0, len(children))
	// unwrap dynamic content
	for _, child := range children {
		if !dynVarVals.Empty() {
			child = addDynamicVars(child, dynVarVals)
		}
		if child.Dynamic == nil {
			res = append(res, child)
			continue
		}
		dynDataCtx := dataCtx
		// found dynamic content
		if !dynVarVals.Empty() {
			// apply dynamic vars to child dynamic var data context
			// this is the case of nested dynamic blocks, child dynamic block
			// can access parent dynamic block vars
			dynDataCtx = maps.Clone(dynDataCtx)
			if diags.Extend(ApplyVars(ctx, dynVarVals, dynDataCtx)) {
				continue
			}
		}

		val, diag := dataspec.EvalAttr(ctx, child.Dynamic.condition, dynDataCtx)
		if diags.Extend(diag) {
			continue
		}
		if val.IsNull() || val.False() {
			continue
		}
		if child.Dynamic.items == nil {
			// no dynamic vars defined, pass parent's
			content, diag := unwrapDynamicContent(ctx, child.Dynamic.children, dataCtx, dynVarVals)
			diags.Extend(diag)
			res = append(res, content...)
			continue
		}
		// iterate over items
		val, diag = dataspec.EvalAttr(ctx, child.Dynamic.items, dynDataCtx)
		if diags.Extend(diag) || val.IsNull() {
			continue
		}
		data := plugindata.Encapsulated.MustFromCty(val)
		if data == nil {
			diags.Append(&hcl.Diagnostic{
				Severity: hcl.DiagError,
				Summary:  "Invalid dynamic block items",
				Detail:   "Dynamic block items must be a list or a map, got nil",
				Subject:  child.Dynamic.items.ValueRange.Ptr(),
			})
			continue
		}
		switch dt := (*data).(type) {
		case nil:
			continue
		case plugindata.List:
			for idx, item := range dt {
				dynVars, diag := parseDynVars(ctx, plugindata.Number(idx), item, child.Dynamic.items.ValueRange)
				if diags.Extend(diag) {
					// infallible
					return
				}
				ndc, diag := unwrapDynamicContent(
					ctx, child.Dynamic.children, dataCtx,
					dynVars,
				)
				if diags.Extend(diag) {
					// stop dynamic block processing on error: it's likely that
					// the error will be repeated for each item and only add noise
					break
				}
				res = append(res, ndc...)
			}
		case plugindata.Map:
			for key, item := range dt {
				dynVars, diag := parseDynVars(ctx, plugindata.String(key), item, child.Dynamic.items.ValueRange)
				if diags.Extend(diag) {
					// infallible
					return
				}
				ndc, diag := unwrapDynamicContent(
					ctx, child.Dynamic.children, dataCtx,
					dynVars,
				)
				if diags.Extend(diag) {
					// stop dynamic block processing on error: it's likely that
					// the error will be repeated for each item and only add noise
					break
				}
				res = append(res, ndc...)
			}
		default:
			diags.Append(&hcl.Diagnostic{
				Severity: hcl.DiagError,
				Summary:  "Invalid dynamic block items",
				Detail:   fmt.Sprintf("Dynamic block items must be a list or a map, got %T", dt),
				Subject:  child.Dynamic.items.ValueRange.Ptr(),
			})
			continue
		}
	}
	return
}

func parseDynVars(ctx context.Context, idx, val plugindata.Data, rng hcl.Range) (parsed *definitions.ParsedVars, diags diagnostics.Diag) {
	// use existing vars parser by creating a synthetic (dynamic_)vars block
	const itemIndex = "dynamic_item_index"
	const item = "dynamic_item"
	return parser.ParseVars(ctx, &hclsyntax.Block{
		Type:            "dynamic_vars",
		TypeRange:       rng,
		OpenBraceRange:  rng,
		CloseBraceRange: rng,
		Body: &hclsyntax.Body{
			SrcRange: rng,
			EndRange: rng,
			Attributes: map[string]*hclsyntax.Attribute{
				itemIndex: {
					Name: itemIndex,
					Expr: &hclsyntax.LiteralValueExpr{
						Val: plugindata.Encapsulated.ValToCty(idx),
					},
					SrcRange:    rng,
					NameRange:   rng,
					EqualsRange: rng,
				},
				item: {
					Name: item,
					Expr: &hclsyntax.LiteralValueExpr{
						Val: plugindata.Encapsulated.ValToCty(val),
					},
					SrcRange:    rng,
					NameRange:   rng,
					EqualsRange: rng,
				},
			},
		},
	}, nil)
}
