package builtin

import (
	"context"
	"fmt"
	"log/slog"
	"strings"
	"testing"

	"github.com/hashicorp/hcl/v2"
	"github.com/stretchr/testify/assert"

	"github.com/blackstork-io/fabric/internal/testtools"
	"github.com/blackstork-io/fabric/pkg/diagnostics"
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

	tt := []struct {
		name          string
		path          string
		glob          string
		expectedData  plugin.Data
		expectedDiags [][]testtools.Assert
	}{
		{
			name: "invalid_json_file_with_glob",
			glob: "testdata/json/invalid.txt",
			expectedDiags: [][]testtools.Assert{{
				testtools.IsError,
				testtools.SummaryEquals("Failed to read the files"),
				testtools.DetailEquals("invalid character 'i' looking for beginning of object key string"),
			}},
		},
		{
			name: "invalid_json_file_with_path",
			path: "testdata/json/invalid.txt",
			expectedDiags: [][]testtools.Assert{{
				testtools.IsError,
				testtools.SummaryEquals("Failed to read a file"),
				testtools.DetailEquals("invalid character 'i' looking for beginning of object key string"),
			}},
		},
		{
			name: "no_params",
			expectedDiags: [][]testtools.Assert{{
				testtools.IsError,
				testtools.SummaryEquals("Failed to parse arguments"),
				testtools.DetailEquals("Either \"glob\" value or \"path\" value must be provided"),
			}},
		},
		{
			name:         "no_glob_matches",
			glob:         "testdata/json/unknown_dir/*.json",
			expectedData: plugin.ListData{},
		},
		{
			name: "no_path_match",
			path: "testdata/json/unknown_dir/does-not-exist.json",
			expectedDiags: [][]testtools.Assert{{
				testtools.IsError,
				testtools.SummaryEquals("Failed to read a file"),
				testtools.DetailEquals("open testdata/json/unknown_dir/does-not-exist.json: no such file or directory"),
			}},
		},
		{
			name: "load_one_file_with_path",
			path: "testdata/json/a.json",
			expectedData: plugin.MapData{
				"property_for": plugin.StringData("a.json"),
			},
		},
		{
			name: "glob_matches_one_file",
			glob: "testdata/json/a.json",
			expectedData: plugin.ListData{
				plugin.MapData{
					"file_name": plugin.StringData("a.json"),
					"file_path": plugin.StringData("testdata/json/a.json"),
					"content": plugin.MapData{
						"property_for": plugin.StringData("a.json"),
					},
				},
			},
		},
		{
			name: "glob_matches_multiple_files",
			glob: "testdata/json/dir/*.json",
			expectedData: plugin.ListData{
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
	}
	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			p := &plugin.Schema{
				DataSources: plugin.DataSources{
					"json": makeJSONDataSource(),
				},
			}

			args := make([]string, 0)
			if tc.path != "" {
				args = append(args, fmt.Sprintf("path = %q", tc.path))
			}
			if tc.glob != "" {
				args = append(args, fmt.Sprintf("glob = %q", tc.glob))
			}
			argsBody := strings.Join(args, ",")

			var diags diagnostics.Diag

			argVal, diag := testtools.Decode(t, p.DataSources["json"].Args, argsBody)
			diags.Extend(diag)

			var data plugin.Data
			if !diags.HasErrors() {
				ctx := context.Background()
				var dgs hcl.Diagnostics
				data, dgs = p.RetrieveData(ctx, "json", &plugin.RetrieveDataParams{Args: argVal})
				diags.ExtendHcl(dgs)
			}
			assert.Equal(t, tc.expectedData, data)
			testtools.CompareDiags(t, nil, diags, tc.expectedDiags)

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
