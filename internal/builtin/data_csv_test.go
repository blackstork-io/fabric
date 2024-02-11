package builtin

import (
	"context"
	"testing"

	"github.com/hashicorp/hcl/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/zclconf/go-cty/cty"

	"github.com/blackstork-io/fabric/plugin"
)

func Test_makeCSVDataSchema(t *testing.T) {
	schema := makeCSVDataSource()
	require.NotNil(t, schema, "expected data source csv")
	assert.NotNil(t, schema.DataFunc)
	assert.NotNil(t, schema.Args)
	assert.NotNil(t, schema.Config)
}

func Test_fetchCSVData(t *testing.T) {
	type results struct {
		Data  plugin.Data
		Diags hcl.Diagnostics
	}
	tt := []struct {
		name      string
		path      string
		delimiter string
		expected  results
	}{
		{
			name:      "comma_delim",
			path:      "testdata/csv/comma.csv",
			delimiter: ",",
			expected: results{
				Data: plugin.ListData{
					plugin.MapData{
						"id":     plugin.StringData("b8fa4bb0-6dd4-45ba-96e0-9a182b2b932e"),
						"active": plugin.BoolData(true),
						"name":   plugin.StringData("Stacey"),
						"age":    plugin.NumberData(26),
						"height": plugin.NumberData(1.98),
					},
					plugin.MapData{
						"id":     plugin.StringData("b0086c49-bcd8-4aae-9f88-4f46b128e709"),
						"active": plugin.BoolData(false),
						"name":   plugin.StringData("Myriam"),
						"age":    plugin.NumberData(33),
						"height": plugin.NumberData(1.81),
					},
					plugin.MapData{
						"id":     plugin.StringData("a12d2a8c-eebc-42b3-be52-1ab0a2969a81"),
						"active": plugin.BoolData(true),
						"name":   plugin.StringData("Oralee"),
						"age":    plugin.NumberData(31),
						"height": plugin.NumberData(2.23),
					},
				},
			},
		},
		{
			name:      "semicolon_delim",
			path:      "testdata/csv/semicolon.csv",
			delimiter: ";",
			expected: results{
				Data: plugin.ListData{
					plugin.MapData{
						"id":     plugin.StringData("b8fa4bb0-6dd4-45ba-96e0-9a182b2b932e"),
						"active": plugin.BoolData(true),
						"name":   plugin.StringData("Stacey"),
						"age":    plugin.NumberData(26),
						"height": plugin.NumberData(1.98),
					},
					plugin.MapData{
						"id":     plugin.StringData("b0086c49-bcd8-4aae-9f88-4f46b128e709"),
						"active": plugin.BoolData(false),
						"name":   plugin.StringData("Myriam"),
						"age":    plugin.NumberData(33),
						"height": plugin.NumberData(1.81),
					},
					plugin.MapData{
						"id":     plugin.StringData("a12d2a8c-eebc-42b3-be52-1ab0a2969a81"),
						"active": plugin.BoolData(true),
						"name":   plugin.StringData("Oralee"),
						"age":    plugin.NumberData(31),
						"height": plugin.NumberData(2.23),
					},
				},
			},
		},
		{
			name: "empty_path",
			expected: results{
				Diags: hcl.Diagnostics{
					{
						Severity: hcl.DiagError,
						Summary:  "path is required",
					},
				},
			},
		},
		{
			name:      "invalid_delimiter",
			path:      "testdata/csv/comma.csv",
			delimiter: "abc",
			expected: results{
				Diags: hcl.Diagnostics{
					{
						Severity: hcl.DiagError,
						Summary:  "delimiter must be a single character",
					},
				},
			},
		},
		{
			name: "default_delimiter",
			path: "testdata/csv/comma.csv",
			expected: results{
				Data: plugin.ListData{
					plugin.MapData{
						"id":     plugin.StringData("b8fa4bb0-6dd4-45ba-96e0-9a182b2b932e"),
						"active": plugin.BoolData(true),
						"name":   plugin.StringData("Stacey"),
						"age":    plugin.NumberData(26),
						"height": plugin.NumberData(1.98),
					},
					plugin.MapData{
						"id":     plugin.StringData("b0086c49-bcd8-4aae-9f88-4f46b128e709"),
						"active": plugin.BoolData(false),
						"name":   plugin.StringData("Myriam"),
						"age":    plugin.NumberData(33),
						"height": plugin.NumberData(1.81),
					},
					plugin.MapData{
						"id":     plugin.StringData("a12d2a8c-eebc-42b3-be52-1ab0a2969a81"),
						"active": plugin.BoolData(true),
						"name":   plugin.StringData("Oralee"),
						"age":    plugin.NumberData(31),
						"height": plugin.NumberData(2.23),
					},
				},
			},
		},
		{
			name: "invalid_path",
			path: "testdata/csv/does_not_exist.csv",
			expected: results{
				Diags: hcl.Diagnostics{
					{
						Severity: hcl.DiagError,
						Summary:  "Failed to read csv file",
						Detail:   "open testdata/csv/does_not_exist.csv: no such file or directory",
					},
				},
			},
		},

		{
			name: "invalid_csv",
			path: "testdata/csv/invalid.csv",
			expected: results{
				Diags: hcl.Diagnostics{
					{
						Severity: hcl.DiagError,
						Summary:  "Failed to read csv file",
						Detail:   "record on line 2: wrong number of fields",
					},
				},
			},
		},
		{
			name:      "empty_csv",
			path:      "testdata/csv/empty.csv",
			delimiter: ",",
			expected: results{
				Data: plugin.ListData{},
			},
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			p := &plugin.Schema{
				DataSources: plugin.DataSources{
					"csv": makeCSVDataSource(),
				},
			}
			delim := cty.StringVal(tc.delimiter)
			if tc.delimiter == "" {
				delim = cty.NullVal(cty.String)
			}
			args := cty.ObjectVal(map[string]cty.Value{
				"path": cty.StringVal(tc.path),
			})
			cfg := cty.ObjectVal(map[string]cty.Value{
				"delimiter": delim,
			})
			ctx := context.Background()
			data, diags := p.RetrieveData(ctx, "csv", &plugin.RetrieveDataParams{Config: cfg, Args: args})
			assert.Equal(t, tc.expected, results{data, diags})
		})
	}
}

func Test_readCSVFileCancellation(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	data, err := readCSVFile(ctx, "testdata/csv/comma.csv", defaultCSVDelimiter)
	assert.Nil(t, data)
	assert.Error(t, context.Canceled, err)
}
