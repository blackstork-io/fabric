package csv

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
	assert.Equal(t, "csv", got.Name)
	assert.Equal(t, "data", got.Kind)
	assert.Equal(t, "blackstork", got.Namespace)
	assert.Equal(t, Version.String(), got.Version.Cast().String())
	assert.Nil(t, got.ConfigSpec)
	assert.NotNil(t, got.InvocationSpec)
}

func TestPlugin_Call(t *testing.T) {
	tt := []struct {
		name      string
		path      string
		delimiter string
		expected  plugininterface.Result
	}{
		{
			name:      "comma_delim",
			path:      "testdata/comma.csv",
			delimiter: ",",
			expected: plugininterface.Result{
				Result: []map[string]any{
					{
						"id":     "b8fa4bb0-6dd4-45ba-96e0-9a182b2b932e",
						"active": true,
						"name":   "Stacey",
						"age":    int64(26),
						"height": float64(1.98),
					},
					{
						"id":     "b0086c49-bcd8-4aae-9f88-4f46b128e709",
						"active": false,
						"name":   "Myriam",
						"age":    int64(33),
						"height": float64(1.81),
					},
					{
						"id":     "a12d2a8c-eebc-42b3-be52-1ab0a2969a81",
						"active": true,
						"name":   "Oralee",
						"age":    int64(31),
						"height": float64(2.23),
					},
				},
			},
		},
		{
			name:      "semicolon_delim",
			path:      "testdata/semicolon.csv",
			delimiter: ";",
			expected: plugininterface.Result{
				Result: []map[string]any{
					{
						"id":     "b8fa4bb0-6dd4-45ba-96e0-9a182b2b932e",
						"active": true,
						"name":   "Stacey",
						"age":    int64(26),
						"height": float64(1.98),
					},
					{
						"id":     "b0086c49-bcd8-4aae-9f88-4f46b128e709",
						"active": false,
						"name":   "Myriam",
						"age":    int64(33),
						"height": float64(1.81),
					},
					{
						"id":     "a12d2a8c-eebc-42b3-be52-1ab0a2969a81",
						"active": true,
						"name":   "Oralee",
						"age":    int64(31),
						"height": float64(2.23),
					},
				},
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
			name:      "invalid_delimiter",
			path:      "testdata/comma.csv",
			delimiter: "abc",
			expected: plugininterface.Result{
				Diags: hcl.Diagnostics{{
					Severity: hcl.DiagError,
					Summary:  "delimiter must be a single character",
				}},
			},
		},
		{
			name: "default_delimiter",
			path: "testdata/comma.csv",
			expected: plugininterface.Result{
				Result: []map[string]any{
					{
						"id":     "b8fa4bb0-6dd4-45ba-96e0-9a182b2b932e",
						"active": true,
						"name":   "Stacey",
						"age":    int64(26),
						"height": float64(1.98),
					},
					{
						"id":     "b0086c49-bcd8-4aae-9f88-4f46b128e709",
						"active": false,
						"name":   "Myriam",
						"age":    int64(33),
						"height": float64(1.81),
					},
					{
						"id":     "a12d2a8c-eebc-42b3-be52-1ab0a2969a81",
						"active": true,
						"name":   "Oralee",
						"age":    int64(31),
						"height": float64(2.23),
					},
				},
			},
		},
		{
			name: "invalid_path",
			path: "testdata/does_not_exist.csv",
			expected: plugininterface.Result{
				Diags: hcl.Diagnostics{{
					Severity: hcl.DiagError,
					Summary:  "Failed to read csv file",
					Detail:   "open testdata/does_not_exist.csv: no such file or directory",
				}},
			},
		},

		{
			name: "invalid_csv",
			path: "testdata/invalid.csv",
			expected: plugininterface.Result{
				Diags: hcl.Diagnostics{{
					Severity: hcl.DiagError,
					Summary:  "Failed to read csv file",
					Detail:   "record on line 2: wrong number of fields",
				}},
			},
		},
		{
			name:      "empty_csv",
			path:      "testdata/empty.csv",
			delimiter: ",",
			expected: plugininterface.Result{
				Result: []map[string]any{},
			},
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			plugin := Plugin{}
			delim := cty.StringVal(tc.delimiter)
			if tc.delimiter == "" {
				delim = cty.NullVal(cty.String)
			}
			args := plugininterface.Args{
				Kind: "data",
				Name: "csv",
				Args: cty.ObjectVal(map[string]cty.Value{
					"path":      cty.StringVal(tc.path),
					"delimiter": delim,
				}),
			}
			got := plugin.Call(args)
			assert.Equal(t, tc.expected, got)
		})
	}
}
