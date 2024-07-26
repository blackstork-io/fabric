package builtin

import (
	"context"
	"fmt"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/blackstork-io/fabric/pkg/diagnostics"
	"github.com/blackstork-io/fabric/pkg/diagnostics/diagtest"
	"github.com/blackstork-io/fabric/plugin"
	"github.com/blackstork-io/fabric/plugin/plugindata"
	"github.com/blackstork-io/fabric/plugin/plugintest"
)

func Test_makeTXTDataSchema(t *testing.T) {
	schema := makeTXTDataSource()
	assert.Nil(t, schema.Config)
	assert.NotNil(t, schema.Args)
	assert.NotNil(t, schema.DataFunc)
}

func Test_fetchTXTData(t *testing.T) {
	tt := []struct {
		name          string
		path          string
		glob          string
		expectedData  plugindata.Data
		expectedDiags diagtest.Asserts
	}{
		{
			name:         "valid_path",
			path:         "testdata/txt/data.txt",
			expectedData: plugindata.String("data_content"),
		},
		{
			name: "with_glob_matches",
			glob: "testdata/txt/dat*.txt",
			expectedData: plugindata.List{
				plugindata.Map{
					"file_name": plugindata.String("data.txt"),
					"file_path": plugindata.String("testdata/txt/data.txt"),
					"content":   plugindata.String("data_content"),
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
			name:         "no_glob_matches",
			glob:         "testdata/txt/does-not-exist*.txt",
			expectedData: plugindata.List{},
		},
		{
			name: "invalid_path",
			path: "testdata/txt/does_not_exist.txt",
			expectedDiags: diagtest.Asserts{{
				diagtest.IsError,
				diagtest.SummaryEquals("Failed to read a file"),
				diagtest.DetailEquals("<nil>: Failed to open a file; open testdata/txt/does_not_exist.txt: no such file or directory"),
			}},
		},
		{
			name:         "empty_file_with_path",
			path:         "testdata/txt/empty.txt",
			expectedData: plugindata.String(""),
		},
		{
			name: "empty_file_with_glob",
			glob: "testdata/txt/empt*.txt",
			expectedData: plugindata.List{
				plugindata.Map{
					"file_name": plugindata.String("empty.txt"),
					"file_path": plugindata.String("testdata/txt/empty.txt"),
					"content":   plugindata.String(""),
				},
			},
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			p := &plugin.Schema{
				DataSources: plugin.DataSources{
					"txt": makeTXTDataSource(),
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

			argVal, diag := plugintest.Decode(t, p.DataSources["txt"].Args, argsBody)
			diags.Extend(diag)

			var data plugindata.Data
			if !diags.HasErrors() {
				ctx := context.Background()
				var dgs diagnostics.Diag
				data, dgs = p.RetrieveData(ctx, "txt", &plugin.RetrieveDataParams{Args: argVal})
				diags.Extend(dgs)
			}
			assert.Equal(t, tc.expectedData, data)
			tc.expectedDiags.AssertMatch(t, diags, nil)
		})
	}
}
