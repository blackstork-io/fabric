package builtin

import (
	"context"
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/hcl/v2"
	"github.com/stretchr/testify/assert"

	"github.com/blackstork-io/fabric/internal/testtools"
	"github.com/blackstork-io/fabric/pkg/diagnostics"
	"github.com/blackstork-io/fabric/plugin"
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
		expectedData  plugin.Data
		expectedDiags [][]testtools.Assert
	}{
		{
			name:         "valid_path",
			path:         "testdata/txt/data.txt",
			expectedData: plugin.StringData("data_content"),
		},
		{
			name: "with_glob_matches",
			glob: "testdata/txt/dat*.txt",
			expectedData: plugin.ListData{
				plugin.MapData{
					"file_name": plugin.StringData("data.txt"),
					"file_path": plugin.StringData("testdata/txt/data.txt"),
					"content":   plugin.StringData("data_content"),
				},
			},
		},
		{
			name: "no_path_no_glob",
			expectedDiags: [][]testtools.Assert{{
				testtools.IsError,
				testtools.SummaryEquals("Failed to parse provided arguments"),
				testtools.DetailEquals("Either \"glob\" value or \"path\" value must be provided"),
			}},
		},
		{
			name:         "no_glob_matches",
			glob:         "testdata/txt/does-not-exist*.txt",
			expectedData: plugin.ListData{},
		},
		{
			name: "invalid_path",
			path: "testdata/txt/does_not_exist.txt",
			expectedDiags: [][]testtools.Assert{{
				testtools.IsError,
				testtools.SummaryEquals("Failed to read a file"),
				testtools.DetailEquals("<nil>: Failed to open a file; open testdata/txt/does_not_exist.txt: no such file or directory"),
			}},
		},
		{
			name:         "empty_file_with_path",
			path:         "testdata/txt/empty.txt",
			expectedData: plugin.StringData(""),
		},
		{
			name: "empty_file_with_glob",
			glob: "testdata/txt/empt*.txt",
			expectedData: plugin.ListData{
				plugin.MapData{
					"file_name": plugin.StringData("empty.txt"),
					"file_path": plugin.StringData("testdata/txt/empty.txt"),
					"content":   plugin.StringData(""),
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

			argVal, diag := testtools.Decode(t, p.DataSources["txt"].Args, argsBody)
			diags.Extend(diag)

			var data plugin.Data
			if !diags.HasErrors() {
				ctx := context.Background()
				var dgs hcl.Diagnostics
				data, dgs = p.RetrieveData(ctx, "txt", &plugin.RetrieveDataParams{Args: argVal})
				diags.ExtendHcl(dgs)
			}
			assert.Equal(t, tc.expectedData, data)
			testtools.CompareDiags(t, nil, diags, tc.expectedDiags)

		})
	}
}
