package builtin

import (
	"context"
	"log/slog"

	"github.com/hashicorp/hcl/v2"
	"github.com/zclconf/go-cty/cty"

	"github.com/blackstork-io/fabric/plugin"
)

func makeInlineDataSource() *plugin.DataSource {
	return &plugin.DataSource{
		DataFunc: fetchInlineData,
		Args:     nil,
	}
}

func fetchInlineData(ctx context.Context, params *plugin.RetrieveDataParams) (plugin.Data, hcl.Diagnostics) {
	if params.Args.IsNull() {
		return nil, hcl.Diagnostics{{
			Severity: hcl.DiagError,
			Summary:  "Failed to parse arguments",
			Detail:   "inline data is required",
		}}
	}
	if !params.Args.Type().IsMapType() && !params.Args.Type().IsObjectType() {
		return nil, hcl.Diagnostics{{
			Severity: hcl.DiagError,
			Summary:  "Failed to parse arguments",
			Detail:   "inline data must be a map",
		}}
	}
	return convertInline(params.Args), nil
}

func convertInline(v cty.Value) plugin.Data {
	if v.IsNull() {
		return nil
	}
	t := v.Type()
	switch {
	case t == cty.String:
		return plugin.StringData(v.AsString())
	case t == cty.Number:
		n, _ := v.AsBigFloat().Float64()
		return plugin.NumberData(n)
	case t == cty.Bool:
		return plugin.BoolData(v.True())
	case t.IsMapType() || t.IsObjectType():
		return convertInlineMap(v)
	case t.IsListType() || t.IsTupleType():
		return convertInlineList(v)
	default:
		slog.Warn("convertInline: failed to match on type", "type", t.GoString(), "value", v.GoString())
		return nil
	}
}

func convertInlineList(v cty.Value) plugin.ListData {
	var result plugin.ListData
	for _, v := range v.AsValueSlice() {
		result = append(result, convertInline(v))
	}
	return result
}

func convertInlineMap(v cty.Value) plugin.MapData {
	result := make(plugin.MapData)
	for k, v := range v.AsValueMap() {
		result[k] = convertInline(v)
	}
	return result
}
