package inline

import (
	"testing"

	"github.com/blackstork-io/fabric/plugininterface/v1"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/zclconf/go-cty/cty"
)

func TestPlugin_GetPlugins(t *testing.T) {
	plugin := Plugin{}
	plugins := plugin.GetPlugins()
	require.Len(t, plugins, 1, "expected 1 plugin")
	got := plugins[0]
	assert.Equal(t, "inline", got.Name)
	assert.Equal(t, "data", got.Kind)
	assert.Equal(t, "blackstork", got.Namespace)
	assert.Equal(t, Version.String(), got.Version.Cast().String())
	assert.Nil(t, got.ConfigSpec)
	assert.NotNil(t, got.InvocationSpec)
}

func TestPlugin_Call(t *testing.T) {
	args := plugininterface.Args{
		Kind: "data",
		Name: "inline",
		Args: cty.ObjectVal(map[string]cty.Value{
			"foo": cty.StringVal("bar"),
			"baz": cty.NumberIntVal(1),
			"qux": cty.BoolVal(true),
			"quux": cty.ListVal([]cty.Value{
				cty.StringVal("corge"),
				cty.StringVal("grault"),
				cty.StringVal("garply"),
			}),
			"quuz": cty.ObjectVal(map[string]cty.Value{
				"garply": cty.StringVal("waldo"),
				"fred":   cty.NumberFloatVal(3.123),
				"plugh":  cty.BoolVal(false),
			}),
			"xyzzy": cty.NullVal(cty.String),
		}),
	}
	plugin := Plugin{}
	got := plugin.Call(args)
	assert.Equal(t, plugininterface.Result{
		Result: map[string]any{
			"foo": "bar",
			"baz": int64(1),
			"qux": true,
			"quux": []any{
				"corge",
				"grault",
				"garply",
			},
			"quuz": map[string]any{
				"garply": "waldo",
				"fred":   float64(3.123),
				"plugh":  false,
			},
			"xyzzy": nil,
		},
	}, got)

}
