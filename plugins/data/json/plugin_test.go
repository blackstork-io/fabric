package json

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
	assert.Equal(t, "json", got.Name)
	assert.Equal(t, "data", got.Kind)
	assert.Equal(t, "blackstork", got.Namespace)
	assert.Equal(t, Version.String(), got.Version.Cast().String())
	assert.Nil(t, got.ConfigSpec)
	assert.NotNil(t, got.InvocationSpec)
}

func TestPlugin_Call(t *testing.T) {
	tt := []struct {
		name     string
		glob     string
		expected plugininterface.Result
	}{
		{
			name: "empty_list",
			glob: "unknown_dir/*.json",
			expected: plugininterface.Result{
				Result: []any{},
			},
		},
		{
			name: "one_file",
			glob: "testdata/a.json",
			expected: plugininterface.Result{
				Result: []any{
					map[string]any{
						"filename": "testdata/a.json",
						"contents": map[string]any{
							"property_for": "a.json",
						},
					},
				},
			},
		},
		{
			name: "dir",
			glob: "testdata/dir/*.json",
			expected: plugininterface.Result{
				Result: []any{
					map[string]any{
						"filename": "testdata/dir/b.json",
						"contents": []any{
							map[string]any{
								"id":           float64(1),
								"property_for": "dir/b.json",
							},
							map[string]any{
								"id":           float64(2),
								"property_for": "dir/b.json",
							},
						},
					},
					map[string]any{
						"filename": "testdata/dir/c.json",
						"contents": []any{
							map[string]any{
								"id":           float64(3),
								"property_for": "dir/c.json",
							},
							map[string]any{
								"id":           float64(4),
								"property_for": "dir/c.json",
							},
						},
					},
				},
			},
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			plugin := Plugin{}
			args := plugininterface.Args{
				Kind: "json",
				Name: "json",
				Args: cty.ObjectVal(map[string]cty.Value{
					"glob": cty.StringVal(tc.glob),
				}),
			}
			got := plugin.Call(args)
			assert.Equal(t, tc.expected, got)
		})
	}

}
