package builtin

import (
	"context"
	"testing"
	"log/slog"

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
	slog.SetLogLoggerLevel(slog.LevelDebug)

	type result struct {
		Data  plugin.Data
		Diags hcl.Diagnostics
	}
	tt := []struct {
		name     string
		path     string
		glob     string
		expected result
	}{
		{
			name: "invalid_json_file_with_glob",
			glob: "testdata/json/invalid.txt",
			expected: result{
				Diags: hcl.Diagnostics{{
					Severity: hcl.DiagError,
					Summary:  "Failed to read JSON files",
					Detail:   "invalid character 'i' looking for beginning of object key string",
				}},
			},
		},
		{
			name: "invalid_json_file_with_path",
			path: "testdata/json/invalid.txt",
			expected: result{
				Diags: hcl.Diagnostics{{
					Severity: hcl.DiagError,
					Summary:  "Failed to read a JSON file",
					Detail:   "invalid character 'i' looking for beginning of object key string",
				}},
			},
		},
		{
			name: "no_params",
			expected: result{
				Diags: hcl.Diagnostics{{
					Severity: hcl.DiagError,
					Summary:  "Failed to parse arguments",
					Detail:   "Either \"glob\" value or \"path\" value must be provided",
				}},
			},
		},
		{
			name: "no_glob_matches",
			glob: "testdata/json/unknown_dir/*.json",
			expected: result{
				Data: plugin.ListData{},
			},
		},
		{
			name: "no_path_match",
			path: "testdata/json/unknown_dir/does-not-exist.json",
			expected: result{
				Diags: hcl.Diagnostics{{
					Severity: hcl.DiagError,
					Summary:  "Failed to read a JSON file",
					Detail:   "open testdata/json/unknown_dir/does-not-exist.json: no such file or directory",
				}},
			},
		},
		{
			name: "load_one_file_with_path",
			path: "testdata/json/a.json",
			expected: result{
				Data: plugin.MapData{
					"property_for": plugin.StringData("a.json"),
				},
			},
		},
		{
			name: "glob_matches_one_file",
			glob: "testdata/json/a.json",
			expected: result{
				Data: plugin.ListData{
					plugin.MapData{
						"file_name": plugin.StringData("a.json"),
						"file_path": plugin.StringData("testdata/json/a.json"),
						"content": plugin.MapData{
							"property_for": plugin.StringData("a.json"),
						},
					},
				},
			},
		},
		{
			name: "glob_matches_multiple_files",
			glob: "testdata/json/dir/*.json",
			expected: result{
				Data: plugin.ListData{
					plugin.MapData{
						"file_name": plugin.StringData("b.json"),
						"file_path": plugin.StringData("testdata/json/dir/b.json"),
						"content": plugin.ListData{
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
						"file_name": plugin.StringData("c.json"),
						"file_path": plugin.StringData("testdata/json/dir/c.json"),
						"content": plugin.ListData{
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
					"path": cty.StringVal(tc.path),
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
