package inline

import (
	"github.com/Masterminds/semver/v3"
	"github.com/blackstork-io/fabric/plugininterface/v1"
	"github.com/hashicorp/hcl/v2/hcldec"
	"github.com/zclconf/go-cty/cty"
)

var Version = semver.MustParse("0.1.0")

type Plugin struct{}

func (Plugin) GetPlugins() []plugininterface.Plugin {
	return []plugininterface.Plugin{
		{
			Namespace:  "blackstork",
			Kind:       "data",
			Name:       "inline",
			Version:    plugininterface.Version(*Version),
			ConfigSpec: nil,
			InvocationSpec: &hcldec.BlockAttrsSpec{
				TypeName:    "inline",
				Required:    true,
				ElementType: cty.DynamicPseudoType,
			},
		},
	}
}

func (p Plugin) Call(args plugininterface.Args) plugininterface.Result {
	return plugininterface.Result{
		Result: p.convert(args.Args),
	}
}

func (p Plugin) convert(v cty.Value) any {
	if v.IsNull() {
		return nil
	}
	t := v.Type()
	switch {
	case t == cty.String:
		return v.AsString()
	case t == cty.Number:
		if v.AsBigFloat().IsInt() {
			n, _ := v.AsBigFloat().Int64()
			return n
		} else {
			n, _ := v.AsBigFloat().Float64()
			return n
		}
	case t == cty.Bool:
		return v.True()
	case t.IsMapType() || t.IsObjectType():
		return p.convertMap(v)
	case t.IsListType():
		return p.convertList(v)
	default:
		return nil
	}
}

func (p Plugin) convertList(v cty.Value) []any {
	var result []any
	for _, v := range v.AsValueSlice() {
		result = append(result, p.convert(v))
	}
	return result
}

func (p Plugin) convertMap(v cty.Value) map[string]any {
	result := make(map[string]any)
	for k, v := range v.AsValueMap() {
		result[k] = p.convert(v)
	}
	return result
}
