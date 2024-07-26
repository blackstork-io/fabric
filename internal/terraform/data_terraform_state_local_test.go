package terraform

import (
	"context"
	"testing"

	"github.com/hashicorp/hcl/v2"
	"github.com/stretchr/testify/assert"
	"github.com/zclconf/go-cty/cty"

	"github.com/blackstork-io/fabric/pkg/diagnostics"
	"github.com/blackstork-io/fabric/plugin"
	"github.com/blackstork-io/fabric/plugin/dataspec"
	"github.com/blackstork-io/fabric/plugin/plugindata"
)

func Test_terraformStateLocalDataSchema(t *testing.T) {
	source := makeTerraformStateLocalDataSource()
	assert.Nil(t, source.Config)
	assert.NotNil(t, source.Args)
	assert.NotNil(t, source.DataFunc)
}

func Test_fetchTerraformStateLocalData_Call(t *testing.T) {
	type result struct {
		data  plugindata.Data
		diags diagnostics.Diag
	}
	tt := []struct {
		name     string
		path     string
		expected result
	}{
		{
			name: "notfound",
			path: "testdata/notfound.tfstate",
			expected: result{
				diags: diagnostics.Diag{
					{
						Severity: hcl.DiagError,
						Summary:  "Failed to read terraform state",
						Detail:   "open testdata/notfound.tfstate: no such file or directory",
					},
				},
			},
		},
		{
			name: "empty_path",
			path: "",
			expected: result{
				diags: diagnostics.Diag{
					{
						Severity: hcl.DiagError,
						Summary:  "Failed to parse arguments",
						Detail:   "path is required",
					},
				},
			},
		},
		{
			name: "valid",
			path: "testdata/terraform.tfstate",
			expected: result{
				data: plugindata.Map{
					"version": plugindata.Number(1),
					"serial":  plugindata.Number(0),
					"modules": plugindata.List{
						plugindata.Map{
							"path": plugindata.List{
								plugindata.String("root"),
							},
							"outputs":   plugindata.Map{},
							"resources": plugindata.Map{},
						},
					},
				},
			},
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			p := Plugin("1.2.3")
			var got result
			got.data, got.diags = p.RetrieveData(context.Background(), "terraform_state_local", &plugin.RetrieveDataParams{
				Args: dataspec.NewBlock([]string{"args"}, map[string]cty.Value{
					"path": cty.StringVal(tc.path),
				}),
			})
			assert.Equal(t, tc.expected, got)
		})
	}
}
