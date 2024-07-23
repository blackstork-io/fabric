package sqlite

import (
	"context"
	"database/sql"
	"encoding/base64"
	"fmt"
	"time"

	"github.com/hashicorp/hcl/v2"
	_ "github.com/mattn/go-sqlite3"
	"github.com/zclconf/go-cty/cty"

	"github.com/blackstork-io/fabric/pkg/diagnostics"
	"github.com/blackstork-io/fabric/plugin"
	"github.com/blackstork-io/fabric/plugin/dataspec"
	"github.com/blackstork-io/fabric/plugin/dataspec/constraint"
)

func makeSqliteDataSource() *plugin.DataSource {
	return &plugin.DataSource{
		Config: &dataspec.RootSpec{
			Attrs: []*dataspec.AttrSpec{
				{
					Name:        "database_uri",
					Type:        cty.String,
					Constraints: constraint.RequiredMeaningful,
					Secret:      true,
				},
			},
		},
		Args: &dataspec.RootSpec{
			Attrs: []*dataspec.AttrSpec{
				{
					Name:        "sql_query",
					Type:        cty.String,
					Constraints: constraint.RequiredNonNull,
				},
				{
					Name: "sql_args",
					Type: cty.List(cty.DynamicPseudoType),
				},
			},
		},
		DataFunc: fetchSqliteData,
	}
}

func parseSqliteArgs(args *dataspec.Block) (string, []any, error) {
	sqlQuery := args.GetAttrVal("sql_query")
	if sqlQuery.IsNull() || sqlQuery.AsString() == "" {
		return "", nil, fmt.Errorf("sql_query is required")
	}
	sqlArgs := args.GetAttrVal("sql_args")
	if sqlArgs.IsNull() || sqlArgs.LengthInt() == 0 {
		return sqlQuery.AsString(), nil, nil
	}
	argsList := sqlArgs.AsValueSlice()
	argsResult := make([]any, len(argsList))
	for i, arg := range argsList {
		switch {
		case arg.IsNull():
			argsResult[i] = nil
		case arg.Type() == cty.Number:
			n, _ := arg.AsBigFloat().Float64()
			argsResult[i] = n
		case arg.Type() == cty.String:
			argsResult[i] = arg.AsString()
		case arg.Type() == cty.Bool:
			argsResult[i] = arg.True()
		default:
			return "", nil, fmt.Errorf("sql_args must be a list of strings, numbers, or booleans")
		}
	}
	return sqlQuery.AsString(), argsResult, nil
}

func fetchSqliteData(ctx context.Context, params *plugin.RetrieveDataParams) (plugin.Data, diagnostics.Diag) {
	dbURI := params.Config.GetAttrVal("database_uri").AsString()
	sqlQuery, sqlArgs, err := parseSqliteArgs(params.Args)
	if err != nil {
		return nil, diagnostics.Diag{{
			Severity: hcl.DiagError,
			Summary:  "Invalid arguments",
			Detail:   err.Error(),
		}}
	}

	db, err := sql.Open("sqlite3", dbURI)
	if err != nil {
		return nil, diagnostics.Diag{{
			Severity: hcl.DiagError,
			Summary:  "Failed to open database",
			Detail:   err.Error(),
		}}
	}
	defer db.Close()
	rows, err := db.QueryContext(ctx, sqlQuery, sqlArgs...)
	if err != nil {
		return nil, diagnostics.Diag{{
			Severity: hcl.DiagError,
			Summary:  "Failed to query database",
			Detail:   err.Error(),
		}}
	}
	// read columns
	columns, err := rows.Columns()
	if err != nil {
		return nil, diagnostics.Diag{{
			Severity: hcl.DiagError,
			Summary:  "Failed to get column names",
			Detail:   err.Error(),
		}}
	}
	result := make(plugin.ListData, 0)

	// read rows
	for rows.Next() {
		// create a map of column name to column value
		columnValArr := make([]nullData, len(columns))
		columnPtrArr := make([]any, len(columns))
		for i := range columns {
			columnPtrArr[i] = &columnValArr[i]
		}
		err = rows.Scan(columnPtrArr...)
		if err != nil {
			return nil, diagnostics.Diag{{
				Severity: hcl.DiagError,
				Summary:  "Failed to scan row",
				Detail:   err.Error(),
			}}
		}
		row := make(plugin.MapData)
		for i, column := range columns {
			if columnValArr[i].valid {
				row[column] = columnValArr[i].data
			} else {
				row[column] = nil
			}
		}
		result = append(result, row)
	}
	if err := rows.Err(); err != nil {
		return nil, diagnostics.Diag{{
			Severity: hcl.DiagError,
			Summary:  "Failed to read rows",
			Detail:   err.Error(),
		}}
	}
	return result, nil
}

type nullData struct {
	data  plugin.Data
	valid bool
}

func (n *nullData) Scan(value any) error {
	if value == nil {
		n.valid = false
		return nil
	}
	switch v := value.(type) {
	case []byte:
		n.data = plugin.StringData(base64.StdEncoding.EncodeToString(v))
	case string:
		n.data = plugin.StringData(v)
	case int64:
		n.data = plugin.NumberData(v)
	case float64:
		n.data = plugin.NumberData(v)
	case bool:
		n.data = plugin.BoolData(v)
	case time.Time:
		n.data = plugin.StringData(v.Format(time.RFC3339))
	default:
		return fmt.Errorf("unsupported type: %T", value)
	}
	n.valid = true
	return nil
}
