package sqlite

import (
	"context"
	"path"
	"testing"

	"github.com/hashicorp/hcl/v2"
	"github.com/stretchr/testify/assert"
	"github.com/zclconf/go-cty/cty"

	"github.com/blackstork-io/fabric/pkg/diagnostics"
	"github.com/blackstork-io/fabric/plugin"
	"github.com/blackstork-io/fabric/plugin/dataspec"
)

func TestSqliteDataSchema(t *testing.T) {
	source := makeSqliteDataSource()
	assert.NotNil(t, source.Config)
	assert.NotNil(t, source.Args)
	assert.NotNil(t, source.DataFunc)
}

func TestSqliteDataCall(t *testing.T) {
	type result struct {
		data  plugin.Data
		diags diagnostics.Diag
	}
	tt := []struct {
		name     string
		cfg      map[string]cty.Value
		args     map[string]cty.Value
		before   func(tb testing.TB, fs testFS) string
		expected result
		canceled bool
	}{
		{
			name: "empty_database_uri",
			cfg: map[string]cty.Value{
				"database_uri": cty.StringVal(""),
			},
			expected: result{
				diags: diagnostics.Diag{
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
			cfg: map[string]cty.Value{
				"database_uri": cty.NullVal(cty.String),
			},
			expected: result{
				diags: diagnostics.Diag{
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
			cfg: map[string]cty.Value{
				"database_uri": cty.StringVal("file:./file.db"),
			},
			args: (map[string]cty.Value{
				"sql_query": cty.StringVal(""),
			}),
			expected: result{
				diags: diagnostics.Diag{
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
			cfg: (map[string]cty.Value{
				"database_uri": cty.StringVal("file:./file.db"),
			}),
			args: (map[string]cty.Value{
				"sql_query": cty.NullVal(cty.String),
			}),
			expected: result{
				diags: diagnostics.Diag{
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
			args: (map[string]cty.Value{
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
			expected: result{
				data: plugin.ListData{},
			},
		},
		{
			name: "non_empty_table",
			args: (map[string]cty.Value{
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
			expected: result{
				data: plugin.ListData{
					plugin.MapData{
						"id":       plugin.NumberData(1),
						"text_val": plugin.StringData("text_1"),
						"num_val":  plugin.NumberData(1),
						"bool_val": plugin.BoolData(true),
					},
					plugin.MapData{
						"id":       plugin.NumberData(2),
						"text_val": plugin.StringData("text_2"),
						"num_val":  plugin.NumberData(2),
						"bool_val": plugin.BoolData(false),
					},
				},
			},
		},
		{
			name: "with_sql_args",
			args: (map[string]cty.Value{
				"sql_query": cty.StringVal("SELECT * FROM testdata WHERE text_val = $1 AND num_val = $2 AND bool_val = $3 AND null_val IS $4;"),
				"sql_args": cty.TupleVal([]cty.Value{
					cty.StringVal("text_2"),
					cty.NumberIntVal(2),
					cty.BoolVal(false),
					cty.NullVal(cty.DynamicPseudoType),
				}),
			}),
			before: func(tb testing.TB, fs testFS) string {
				dsn := "file:" + path.Join(fs.path, "file.db")
				prepareTestDB(tb, testData{
					dsn:    dsn,
					schema: "CREATE TABLE testdata (id INTEGER PRIMARY KEY, text_val TEXT, num_val INTEGER, bool_val BOOLEAN, null_val TEXT DEFAULT NULL)",
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
			expected: result{
				data: plugin.ListData{
					plugin.MapData{
						"id":       plugin.NumberData(2),
						"text_val": plugin.StringData("text_2"),
						"num_val":  plugin.NumberData(2),
						"bool_val": plugin.BoolData(false),
						"null_val": nil,
					},
				},
			},
		},
		{
			name: "missing_sql_args",
			args: (map[string]cty.Value{
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
			expected: result{
				diags: diagnostics.Diag{
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
			args: (map[string]cty.Value{
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
			expected: result{
				diags: diagnostics.Diag{
					{
						Severity: hcl.DiagError,
						Summary:  "Failed to query database",
						Detail:   "no such table: testdata",
					},
				},
			},
		},
		{
			name: "canceled",
			cfg: (map[string]cty.Value{
				"database_uri": cty.StringVal("file:./file.db"),
			}),
			args: (map[string]cty.Value{
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
			expected: result{
				diags: diagnostics.Diag{{
					Severity: hcl.DiagError,
					Summary:  "Failed to query database",
					Detail:   "context canceled",
				}},
			},
			canceled: true,
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			p := Plugin("1.2.3")
			params := plugin.RetrieveDataParams{
				Config: dataspec.NewBlock([]string{"config"}, tc.cfg),
				Args:   dataspec.NewBlock([]string{"args"}, tc.args),
			}
			if tc.before != nil {
				fs := makeTestFS(t)
				dsn := tc.before(t, fs)
				params.Config = dataspec.NewBlock([]string{"config"}, map[string]cty.Value{
					"database_uri": cty.StringVal(dsn),
				})
			}
			ctx := context.Background()
			if tc.canceled {
				nextCtx, cancel := context.WithCancel(ctx)
				cancel()
				ctx = nextCtx
			}
			data, diags := p.RetrieveData(ctx, "sqlite", &params)
			assert.Equal(t, tc.expected, result{data, diags}, "unexpected result")
		})
	}
}
