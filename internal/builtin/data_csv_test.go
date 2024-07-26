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
	"github.com/blackstork-io/fabric/plugin/plugindata"
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
		expectedData  plugindata.Data
		expectedDiags diagtest.Asserts
	}{
		{
			name:      "comma_delim_path",
			path:      "testdata/csv/comma.csv",
			delimiter: ",",
			expectedData: plugindata.List{
				plugindata.Map{
					"id":     plugindata.String("b8fa4bb0-6dd4-45ba-96e0-9a182b2b932e"),
					"active": plugindata.Bool(true),
					"name":   plugindata.String("Stacey"),
					"age":    plugindata.Number(26),
					"height": plugindata.Number(1.98),
				},
				plugindata.Map{
					"id":     plugindata.String("b0086c49-bcd8-4aae-9f88-4f46b128e709"),
					"active": plugindata.Bool(false),
					"name":   plugindata.String("Myriam"),
					"age":    plugindata.Number(33),
					"height": plugindata.Number(1.81),
				},
				plugindata.Map{
					"id":     plugindata.String("a12d2a8c-eebc-42b3-be52-1ab0a2969a81"),
					"active": plugindata.Bool(true),
					"name":   plugindata.String("Oralee"),
					"age":    plugindata.Number(31),
					"height": plugindata.Number(2.23),
				},
			},
		},
		{
			name:      "semicolon_delim_path",
			path:      "testdata/csv/semicolon.csv",
			delimiter: ";",
			expectedData: plugindata.List{
				plugindata.Map{
					"id":     plugindata.String("b8fa4bb0-6dd4-45ba-96e0-9a182b2b932e"),
					"active": plugindata.Bool(true),
					"name":   plugindata.String("Stacey"),
					"age":    plugindata.Number(26),
					"height": plugindata.Number(1.98),
				},
				plugindata.Map{
					"id":     plugindata.String("b0086c49-bcd8-4aae-9f88-4f46b128e709"),
					"active": plugindata.Bool(false),
					"name":   plugindata.String("Myriam"),
					"age":    plugindata.Number(33),
					"height": plugindata.Number(1.81),
				},
				plugindata.Map{
					"id":     plugindata.String("a12d2a8c-eebc-42b3-be52-1ab0a2969a81"),
					"active": plugindata.Bool(true),
					"name":   plugindata.String("Oralee"),
					"age":    plugindata.Number(31),
					"height": plugindata.Number(2.23),
				},
			},
		},
		{
			name:      "comma_delim_glob",
			glob:      "testdata/csv/comm*.csv",
			delimiter: ",",
			expectedData: plugindata.List{
				plugindata.Map{
					"file_name": plugindata.String("comma.csv"),
					"file_path": plugindata.String("testdata/csv/comma.csv"),
					"content": plugindata.List{
						plugindata.Map{
							"id":     plugindata.String("b8fa4bb0-6dd4-45ba-96e0-9a182b2b932e"),
							"active": plugindata.Bool(true),
							"name":   plugindata.String("Stacey"),
							"age":    plugindata.Number(26),
							"height": plugindata.Number(1.98),
						},
						plugindata.Map{
							"id":     plugindata.String("b0086c49-bcd8-4aae-9f88-4f46b128e709"),
							"active": plugindata.Bool(false),
							"name":   plugindata.String("Myriam"),
							"age":    plugindata.Number(33),
							"height": plugindata.Number(1.81),
						},
						plugindata.Map{
							"id":     plugindata.String("a12d2a8c-eebc-42b3-be52-1ab0a2969a81"),
							"active": plugindata.Bool(true),
							"name":   plugindata.String("Oralee"),
							"age":    plugindata.Number(31),
							"height": plugindata.Number(2.23),
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
			expectedData: plugindata.List{
				plugindata.Map{
					"id":     plugindata.String("b8fa4bb0-6dd4-45ba-96e0-9a182b2b932e"),
					"active": plugindata.Bool(true),
					"name":   plugindata.String("Stacey"),
					"age":    plugindata.Number(26),
					"height": plugindata.Number(1.98),
				},
				plugindata.Map{
					"id":     plugindata.String("b0086c49-bcd8-4aae-9f88-4f46b128e709"),
					"active": plugindata.Bool(false),
					"name":   plugindata.String("Myriam"),
					"age":    plugindata.Number(33),
					"height": plugindata.Number(1.81),
				},
				plugindata.Map{
					"id":     plugindata.String("a12d2a8c-eebc-42b3-be52-1ab0a2969a81"),
					"active": plugindata.Bool(true),
					"name":   plugindata.String("Oralee"),
					"age":    plugindata.Number(31),
					"height": plugindata.Number(2.23),
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
			expectedData: plugindata.List{},
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
			var data plugindata.Data
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
