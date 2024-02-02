package sqlite

import (
	"path"
	"testing"

	"github.com/blackstork-io/fabric/plugininterface/v1"
	"github.com/hashicorp/hcl/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/zclconf/go-cty/cty"
)

func TestPlugin_GetPlugins(t *testing.T) {
	plugin := Plugin{}
	plugins := plugin.GetPlugins()
	require.Len(t, plugins, 1, "expected 1 plugin")
	got := plugins[0]
	assert.Equal(t, "sqlite", got.Name)
	assert.Equal(t, "data", got.Kind)
	assert.Equal(t, "blackstork", got.Namespace)
	assert.Equal(t, Version.String(), got.Version.Cast().String())
	assert.NotNil(t, got.ConfigSpec)
	assert.NotNil(t, got.InvocationSpec)
}

func TestPlugin_Call(t *testing.T) {
	tt := []struct {
		name     string
		cfg      cty.Value
		args     cty.Value
		before   func(tb testing.TB, fs testFS) string
		expected plugininterface.Result
	}{
		{
			name: "empty_database_uri",
			cfg: cty.ObjectVal(map[string]cty.Value{
				"database_uri": cty.StringVal(""),
			}),
			expected: plugininterface.Result{
				Diags: hcl.Diagnostics{
					{
						Severity: hcl.DiagError,
						Summary:  "Invalid configuration",
						Detail:   "database_uri is required",
					},
				},
			},
		},
		{
			name: "nil_database_uri",
			cfg: cty.ObjectVal(map[string]cty.Value{
				"database_uri": cty.NullVal(cty.String),
			}),
			expected: plugininterface.Result{
				Diags: hcl.Diagnostics{
					{
						Severity: hcl.DiagError,
						Summary:  "Invalid configuration",
						Detail:   "database_uri is required",
					},
				},
			},
		},
		{
			name: "empty_sql_query",
			cfg: cty.ObjectVal(map[string]cty.Value{
				"database_uri": cty.StringVal("file:./file.db"),
			}),
			args: cty.ObjectVal(map[string]cty.Value{
				"sql_query": cty.StringVal(""),
			}),
			expected: plugininterface.Result{
				Diags: hcl.Diagnostics{
					{
						Severity: hcl.DiagError,
						Summary:  "Invalid arguments",
						Detail:   "sql_query is required",
					},
				},
			},
		},
		{
			name: "nil_sql_query",
			cfg: cty.ObjectVal(map[string]cty.Value{
				"database_uri": cty.StringVal("file:./file.db"),
			}),
			args: cty.ObjectVal(map[string]cty.Value{
				"sql_query": cty.NullVal(cty.String),
			}),
			expected: plugininterface.Result{
				Diags: hcl.Diagnostics{
					{
						Severity: hcl.DiagError,
						Summary:  "Invalid arguments",
						Detail:   "sql_query is required",
					},
				},
			},
		},
		{
			name: "empty_table",
			args: cty.ObjectVal(map[string]cty.Value{
				"sql_query": cty.StringVal("SELECT * FROM testdata"),
				"sql_args":  cty.ListValEmpty(cty.DynamicPseudoType),
			}),
			before: func(tb testing.TB, fs testFS) string {
				dsn := "file:" + path.Join(fs.path, "file.db")
				prepareTestDB(tb, testData{
					dsn:    dsn,
					schema: "CREATE TABLE testdata (id INTEGER PRIMARY KEY, text_val TEXT)",
					data:   []map[string]any{},
				})
				return dsn
			},
			expected: plugininterface.Result{
				Result: []map[string]any{},
			},
		},
		{
			name: "non_empty_table",
			args: cty.ObjectVal(map[string]cty.Value{
				"sql_query": cty.StringVal("SELECT * FROM testdata"),
				"sql_args":  cty.ListValEmpty(cty.DynamicPseudoType),
			}),
			before: func(tb testing.TB, fs testFS) string {
				dsn := "file:" + path.Join(fs.path, "file.db")
				prepareTestDB(tb, testData{
					dsn:    dsn,
					schema: "CREATE TABLE testdata (id INTEGER PRIMARY KEY, text_val TEXT, num_val INTEGER, bool_val BOOLEAN)",
					data: []map[string]any{
						{
							"id":       int64(1),
							"text_val": "text_1",
							"num_val":  int64(1),
							"bool_val": true,
						},
						{
							"id":       int64(2),
							"text_val": "text_2",
							"num_val":  int64(2),
							"bool_val": false,
						},
					},
				})
				return dsn
			},
			expected: plugininterface.Result{
				Result: []map[string]any{
					{
						"id":       int64(1),
						"text_val": "text_1",
						"num_val":  int64(1),
						"bool_val": true,
					},
					{
						"id":       int64(2),
						"text_val": "text_2",
						"num_val":  int64(2),
						"bool_val": false,
					},
				},
			},
		},
		{
			name: "with_sql_args",
			args: cty.ObjectVal(map[string]cty.Value{
				"sql_query": cty.StringVal("SELECT * FROM testdata WHERE bool_val = $1;"),
				"sql_args":  cty.ListVal([]cty.Value{cty.BoolVal(false)}),
			}),
			before: func(tb testing.TB, fs testFS) string {
				dsn := "file:" + path.Join(fs.path, "file.db")
				prepareTestDB(tb, testData{
					dsn:    dsn,
					schema: "CREATE TABLE testdata (id INTEGER PRIMARY KEY, text_val TEXT, num_val INTEGER, bool_val BOOLEAN)",
					data: []map[string]any{
						{
							"id":       int64(1),
							"text_val": "text_1",
							"num_val":  int64(1),
							"bool_val": true,
						},
						{
							"id":       int64(2),
							"text_val": "text_2",
							"num_val":  int64(2),
							"bool_val": false,
						},
					},
				})
				return dsn
			},
			expected: plugininterface.Result{
				Result: []map[string]any{
					{
						"id":       int64(2),
						"text_val": "text_2",
						"num_val":  int64(2),
						"bool_val": false,
					},
				},
			},
		},
		{
			name: "missing_sql_args",
			args: cty.ObjectVal(map[string]cty.Value{
				"sql_query": cty.StringVal("SELECT * FROM testdata WHERE bool_val = $1;"),
				"sql_args":  cty.ListValEmpty(cty.DynamicPseudoType),
			}),
			before: func(tb testing.TB, fs testFS) string {
				dsn := "file:" + path.Join(fs.path, "file.db")
				prepareTestDB(tb, testData{
					dsn:    dsn,
					schema: "CREATE TABLE testdata (id INTEGER PRIMARY KEY, text_val TEXT, num_val INTEGER, bool_val BOOLEAN)",
					data: []map[string]any{
						{
							"id":       int64(1),
							"text_val": "text_1",
							"num_val":  int64(1),
							"bool_val": true,
						},
						{
							"id":       int64(2),
							"text_val": "text_2",
							"num_val":  int64(2),
							"bool_val": false,
						},
					},
				})
				return dsn
			},
			expected: plugininterface.Result{
				Diags: hcl.Diagnostics{
					{
						Severity: hcl.DiagError,
						Summary:  "Failed to query database",
						Detail:   "not enough args to execute query: want 1 got 0",
					},
				},
			},
		},
		{
			name: "table_not_found",
			args: cty.ObjectVal(map[string]cty.Value{
				"sql_query": cty.StringVal("SELECT * FROM testdata"),
				"sql_args":  cty.ListValEmpty(cty.DynamicPseudoType),
			}),
			before: func(tb testing.TB, fs testFS) string {
				dsn := "file:" + path.Join(fs.path, "file.db")
				prepareTestDB(tb, testData{
					dsn:    dsn,
					schema: "CREATE TABLE testdata_other (id INTEGER PRIMARY KEY)",
					data:   []map[string]any{},
				})
				return dsn
			},
			expected: plugininterface.Result{
				Diags: hcl.Diagnostics{
					{
						Severity: hcl.DiagError,
						Summary:  "Failed to query database",
						Detail:   "no such table: testdata",
					},
				},
			},
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			plugin := Plugin{}
			args := plugininterface.Args{
				Kind:   "data",
				Name:   "sqlite",
				Config: tc.cfg,
				Args:   tc.args,
			}
			if tc.before != nil {
				fs := makeTestFS(t)
				dsn := tc.before(t, fs)
				args.Config = cty.ObjectVal(map[string]cty.Value{
					"database_uri": cty.StringVal(dsn),
				})
			}
			got := plugin.Call(args)
			assert.Equal(t, tc.expected, got)
		})
	}

}
