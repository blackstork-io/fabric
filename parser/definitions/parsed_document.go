package definitions

import (
	"context"
	"slices"

	"github.com/hashicorp/hcl/v2"

	"github.com/blackstork-io/fabric/parser/evaluation"
	"github.com/blackstork-io/fabric/pkg/diagnostics"
	"github.com/blackstork-io/fabric/plugin"
)

type ParsedDocument struct {
	Meta    *MetaBlock
	Content []Renderable
	Data    []*ParsedData
}

// result has a shape map[plugin_name]map[block_name]plugin_result.
func (d *ParsedDocument) evalData(ctx context.Context, caller evaluation.DataCaller) (result plugin.MapData, diags diagnostics.Diag) {
	// TODO: can be parallel:

	result = plugin.MapData{}
	for _, node := range d.Data {
		res, diag := caller.CallData(
			ctx,
			node.PluginName,
			node.Config,
			node.Invocation,
		)
		if diags.Extend(diag) {
			continue
		}

		var pluginNameRes plugin.MapData
		if m, found := result[node.PluginName]; found {
			pluginNameRes = m.(plugin.MapData)
		} else {
			pluginNameRes = plugin.MapData{}
			result[node.PluginName] = pluginNameRes
		}

		if _, found := pluginNameRes[node.BlockName]; found {
			diags.Append(&hcl.Diagnostic{
				Severity: hcl.DiagWarning,
				Summary:  "Data conflict",
				Detail: ("Result of this block overwrites results from the previous invocation. " +
					"Creating multiple anonymous 'data ref' with the same 'base' is ill-advised, " +
					"we recommend naming all 'data ref' blocks uniquely"),
				Subject: node.Invocation.DefRange().Ptr(),
			})
		}
		pluginNameRes[node.BlockName] = res
	}
	return
}

func (d *ParsedDocument) Render(ctx context.Context, caller evaluation.PluginCaller) (result []string, diags diagnostics.Diag) {
	dataResult, diags := d.evalData(ctx, caller)
	if diags.HasErrors() {
		return
	}
	var res evaluation.Result
	document := plugin.ConvMapData{
		BlockKindContent: &res,
	}
	if d.Meta != nil {
		document[BlockKindMeta] = d.Meta.AsJQData()
	}

	dataCtx := evaluation.NewDataContext(plugin.ConvMapData{
		BlockKindData:     dataResult,
		BlockKindDocument: document,
	})
	slices.SortStableFunc(d.Content, func(a, b Renderable) int {
		ao, _ := caller.ContentInvocationOrder(ctx, a.Name())
		bo, _ := caller.ContentInvocationOrder(ctx, b.Name())
		return ao.Weight() - bo.Weight()
	})
	for _, content := range d.Content {
		diags.Extend(
			content.Render(ctx, caller, dataCtx.Share(), &res),
		)
	}
	result = res.AsGoType()
	return
}
