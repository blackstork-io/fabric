package terraform

import (
	"testing"

	"github.com/blackstork-io/fabric/plugininterface/v1"
	"github.com/hashicorp/hcl/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/zclconf/go-cty/cty"
)

func TestPlugin_GetPlugins(t *testing.T) {
	plugin := Plugin{}
	plugins := plugin.GetPlugins()
	require.Len(t, plugins, 1, "expected 1 plugin")
	got := plugins[0]
	assert.Equal(t, "terraform_state_local", got.Name)
	assert.Equal(t, "data", got.Kind)
	assert.Equal(t, "blackstork", got.Namespace)
	assert.Equal(t, Version.String(), got.Version.Cast().String())
	assert.Nil(t, got.ConfigSpec)
	assert.NotNil(t, got.InvocationSpec)
}

func TestPlugin_Call(t *testing.T) {
	tt := []struct {
		name     string
		path     string
		expected plugininterface.Result
	}{
		{
			name: "notfound",
			path: "testdata/notfound.tfstate",
			expected: plugininterface.Result{
				Diags: hcl.Diagnostics{
					{
						Severity: hcl.DiagError,
						Summary:  "Failed to read terraform state",
						Detail:   "open testdata/notfound.tfstate: no such file or directory",
					},
				},
			},
		},
		{
			name: "empty_path",
			path: "",
			expected: plugininterface.Result{
				Diags: hcl.Diagnostics{
					{
						Severity: hcl.DiagError,
						Summary:  "Failed to parse arguments",
						Detail:   "path is required",
					},
				},
			},
		},
		{
			name: "valid",
			path: "testdata/terraform.tfstate",
			expected: plugininterface.Result{
				Result: map[string]any{
					"version": float64(1),
					"serial":  float64(0),
					"modules": []any{
						map[string]any{
							"path":      []any{"root"},
							"outputs":   map[string]any{},
							"resources": map[string]any{},
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
				Kind: "data",
				Name: "terraform_state_local",
				Args: cty.ObjectVal(map[string]cty.Value{
					"path": cty.StringVal(tc.path),
				}),
			}
			got := plugin.Call(args)
			assert.Equal(t, tc.expected, got)
		})
	}

}
