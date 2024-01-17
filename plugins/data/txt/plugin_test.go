package txt

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
	assert.Equal(t, "txt", got.Name)
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
			name: "valid_path",
			path: "testdata/data.txt",
			expected: plugininterface.Result{
				Result: "data_content",
			},
		},
		{
			name: "empty_path",
			expected: plugininterface.Result{
				Diags: hcl.Diagnostics{{
					Severity: hcl.DiagError,
					Summary:  "path is required",
				}},
			},
		},
		{
			name: "invalid_path",
			path: "testdata/does_not_exist.txt",
			expected: plugininterface.Result{
				Diags: hcl.Diagnostics{{
					Severity: hcl.DiagError,
					Summary:  "Failed to open txt file",
					Detail:   "open testdata/does_not_exist.txt: no such file or directory",
				}},
			},
		},
		{
			name: "empty_file",
			path: "testdata/empty.txt",
			expected: plugininterface.Result{
				Result: "",
			},
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			plugin := Plugin{}
			args := plugininterface.Args{
				Kind: "data",
				Name: "txt",
				Args: cty.ObjectVal(map[string]cty.Value{
					"path": cty.StringVal(tc.path),
				}),
			}
			got := plugin.Call(args)
			assert.Equal(t, tc.expected, got)
		})
	}
}
