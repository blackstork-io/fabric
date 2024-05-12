package builtin

import (
	"context"
	"testing"

	"github.com/hashicorp/hcl/v2"
	"github.com/stretchr/testify/assert"
	"github.com/zclconf/go-cty/cty"

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
	type result struct {
		Data  plugin.Data
		Diags diagnostics.Diag
	}
	tt := []struct {
		name     string
		path     string
		expected result
	}{
		{
			name: "valid_path",
			path: "testdata/txt/data.txt",
			expected: result{
				Data: plugin.StringData("data_content"),
			},
		},
		{
			name: "empty_path",
			expected: result{
				Diags: diagnostics.Diag{{
					Severity: hcl.DiagError,
					Summary:  "Failed to parse arguments",
					Detail:   "path is required",
				}},
			},
		},
		{
			name: "invalid_path",
			path: "testdata/txt/does_not_exist.txt",
			expected: result{
				Diags: diagnostics.Diag{{
					Severity: hcl.DiagError,
					Summary:  "Failed to open txt file",
					Detail:   "open testdata/txt/does_not_exist.txt: no such file or directory",
				}},
			},
		},
		{
			name: "empty_file",
			path: "testdata/txt/empty.txt",
			expected: result{
				Data: plugin.StringData(""),
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
			data, diags := p.RetrieveData(context.Background(), "txt", &plugin.RetrieveDataParams{
				Args: cty.ObjectVal(map[string]cty.Value{
					"path": cty.StringVal(tc.path),
				}),
			})
			assert.Equal(t, tc.expected, result{data, diags})
		})
	}
}
