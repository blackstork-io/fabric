package builtin

import (
	"context"
	"testing"

	"github.com/hashicorp/hcl/v2"
	"github.com/stretchr/testify/assert"
	"github.com/zclconf/go-cty/cty"

	"github.com/blackstork-io/fabric/plugin"
)

func Test_makeJSONDataSchema(t *testing.T) {
	schema := makeJSONDataSource()
	assert.Nil(t, schema.Config)
	assert.NotNil(t, schema.Args)
	assert.NotNil(t, schema.DataFunc)
}

func Test_fetchJSONData(t *testing.T) {
	type result struct {
		Data  plugin.Data
		Diags hcl.Diagnostics
	}
	tt := []struct {
		name     string
		glob     string
		expected result
	}{
		{
			name: "invalid_json_file",
			glob: "testdata/json/invalid.txt",
			expected: result{
				Diags: hcl.Diagnostics{{
					Severity: hcl.DiagError,
					Summary:  "Failed to read json files",
					Detail:   "invalid character 'i' looking for beginning of object key string",
				}},
			},
		},
		{
			name: "empty_glob",
			expected: result{
				Diags: hcl.Diagnostics{{
					Severity: hcl.DiagError,
					Summary:  "Failed to parse arguments",
					Detail:   "glob is required",
				}},
			},
		},
		{
			name: "empty_list",
			glob: "testdata/json/unknown_dir/*.json",
			expected: result{
				Data: plugin.ListData{},
			},
		},
		{
			name: "one_file",
			glob: "testdata/json/a.json",
			expected: result{
				Data: plugin.ListData{
					plugin.MapData{
						"filename": plugin.StringData("testdata/json/a.json"),
						"contents": plugin.MapData{
							"property_for": plugin.StringData("a.json"),
						},
					},
				},
			},
		},
		{
			name: "dir",
			glob: "testdata/json/dir/*.json",
			expected: result{
				Data: plugin.ListData{
					plugin.MapData{
						"filename": plugin.StringData("testdata/json/dir/b.json"),
						"contents": plugin.ListData{
							plugin.MapData{
								"id":           plugin.NumberData(1),
								"property_for": plugin.StringData("dir/b.json"),
							},
							plugin.MapData{
								"id":           plugin.NumberData(2),
								"property_for": plugin.StringData("dir/b.json"),
							},
						},
					},
					plugin.MapData{
						"filename": plugin.StringData("testdata/json/dir/c.json"),
						"contents": plugin.ListData{
							plugin.MapData{
								"id":           plugin.NumberData(3),
								"property_for": plugin.StringData("dir/c.json"),
							},
							plugin.MapData{
								"id":           plugin.NumberData(4),
								"property_for": plugin.StringData("dir/c.json"),
							},
						},
					},
				},
			},
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			p := &plugin.Schema{
				DataSources: plugin.DataSources{
					"json": makeJSONDataSource(),
				},
			}
			data, diags := p.RetrieveData(context.Background(), "json", &plugin.RetrieveDataParams{
				Args: cty.ObjectVal(map[string]cty.Value{
					"glob": cty.StringVal(tc.glob),
				}),
			})
			assert.Equal(t, tc.expected, result{data, diags})
		})
	}
}

func Test_readJSONFilesCancellation(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	data, err := readJSONFiles(ctx, "testdata/json/a.json")
	assert.Nil(t, data)
	assert.Error(t, context.Canceled, err)
}
