package builtin

import (
	"context"
	"fmt"
	"log/slog"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/blackstork-io/fabric/pkg/diagnostics"
	"github.com/blackstork-io/fabric/pkg/diagnostics/diagtest"
	"github.com/blackstork-io/fabric/plugin"
	"github.com/blackstork-io/fabric/plugin/plugintest"
)

func Test_makeCSVDataSchema(t *testing.T) {
	schema := makeCSVDataSource()
	require.NotNil(t, schema, "expected data source csv")
	assert.NotNil(t, schema.DataFunc)
	assert.NotNil(t, schema.Args)
	assert.NotNil(t, schema.Config)
}

func Test_fetchCSVData(t *testing.T) {
	slog.SetLogLoggerLevel(slog.LevelDebug)

	tt := []struct {
		name          string
		path          string
		glob          string
		delimiter     string
		expectedData  plugin.Data
		expectedDiags diagtest.Asserts
	}{
		{
			name:      "comma_delim_path",
			path:      "testdata/csv/comma.csv",
			delimiter: ",",
			expectedData: plugin.ListData{
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
		{
			name:      "semicolon_delim_path",
			path:      "testdata/csv/semicolon.csv",
			delimiter: ";",
			expectedData: plugin.ListData{
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
		{
			name:      "comma_delim_glob",
			glob:      "testdata/csv/comm*.csv",
			delimiter: ",",
			expectedData: plugin.ListData{
				plugin.MapData{
					"file_name": plugin.StringData("comma.csv"),
					"file_path": plugin.StringData("testdata/csv/comma.csv"),
					"content": plugin.ListData{
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
		},
		{
			name: "no_path_no_glob",
			expectedDiags: diagtest.Asserts{{
				diagtest.IsError,
				diagtest.SummaryEquals("Failed to parse provided arguments"),
				diagtest.DetailEquals("Either \"glob\" value or \"path\" value must be provided"),
			}},
		},
		{
			name:      "invalid_delimiter",
			path:      "testdata/csv/comma.csv",
			delimiter: "abc",
			expectedDiags: diagtest.Asserts{{
				diagtest.IsError,
				diagtest.DetailContains(
					`The length`, `"delimiter"`, `exactly 1`,
				),
			}},
		},
		{
			name: "default_delimiter",
			path: "testdata/csv/comma.csv",
			expectedData: plugin.ListData{
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
		{
			name: "invalid_path",
			path: "testdata/csv/does_not_exist.csv",
			expectedDiags: diagtest.Asserts{{
				diagtest.IsError,
				diagtest.SummaryContains("Failed to read a file"),
				diagtest.DetailContains("no such file or directory"),
			}},
		},

		{
			name: "invalid_csv",
			path: "testdata/csv/invalid.csv",
			expectedDiags: diagtest.Asserts{{
				diagtest.IsError,
				diagtest.SummaryContains("Failed to read a file"),
				diagtest.DetailEquals("record on line 2: wrong number of fields"),
			}},
		},
		{
			name:         "empty_csv",
			path:         "testdata/csv/empty.csv",
			delimiter:    ",",
			expectedData: plugin.ListData{},
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			p := &plugin.Schema{
				DataSources: plugin.DataSources{
					"csv": makeCSVDataSource(),
				},
			}
			config := ""
			if tc.delimiter != "" {
				config = fmt.Sprintf("delimiter = %q", tc.delimiter)
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
			argVal, diag := plugintest.Decode(t, p.DataSources["csv"].Args, argsBody)
			diags.Extend(diag)
			cfgVal, diag := plugintest.Decode(t, p.DataSources["csv"].Config, config)
			diags.Extend(diag)
			var data plugin.Data
			if !diags.HasErrors() {
				ctx := context.Background()
				var dgs diagnostics.Diag
				data, dgs = p.RetrieveData(ctx, "csv", &plugin.RetrieveDataParams{Config: cfgVal, Args: argVal})
				diags.Extend(dgs)
			}
			assert.Equal(t, tc.expectedData, data)
			tc.expectedDiags.AssertMatch(t, diags, nil)
		})
	}
}

func Test_readCSVFileCancellation(t *testing.T) {
	const defaultCSVDelimiter = ','
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	data, err := readAndDecodeCSVFile(ctx, "testdata/csv/comma.csv", defaultCSVDelimiter)
	assert.Nil(t, data)
	assert.Error(t, context.Canceled, err)
}
