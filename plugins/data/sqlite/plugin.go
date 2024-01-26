package sqlite

import (
	"database/sql"
	"errors"

	"github.com/Masterminds/semver/v3"
	"github.com/blackstork-io/fabric/plugininterface/v1"
	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/hcldec"
	_ "github.com/mattn/go-sqlite3"
	"github.com/zclconf/go-cty/cty"
)

var Version = semver.MustParse("0.1.0")

type Plugin struct{}

func (Plugin) GetPlugins() []plugininterface.Plugin {
	return []plugininterface.Plugin{
		{
			Namespace: "blackstork",
			Kind:      "data",
			Name:      "sqlite",
			Version:   plugininterface.Version(*Version),
			ConfigSpec: &hcldec.ObjectSpec{
				"database_uri": &hcldec.AttrSpec{
					Name:     "database_uri",
					Type:     cty.String,
					Required: true,
				},
			},
			InvocationSpec: &hcldec.ObjectSpec{
				"sql_query": &hcldec.AttrSpec{
					Name:     "sql_query",
					Type:     cty.String,
					Required: true,
				},
				"sql_args": &hcldec.AttrSpec{
					Name:     "sql_args",
					Type:     cty.List(cty.DynamicPseudoType),
					Required: false,
				},
			},
		},
	}
}

func (Plugin) parseConfig(cfg cty.Value) (string, error) {
	dbURI := cfg.GetAttr("database_uri")
	if dbURI.IsNull() || dbURI.AsString() == "" {
		return "", errors.New("database_uri is required")
	}
	return dbURI.AsString(), nil
}

func (Plugin) parseArgs(args cty.Value) (string, []any, error) {
	sqlQuery := args.GetAttr("sql_query")
	if sqlQuery.IsNull() || sqlQuery.AsString() == "" {
		return "", nil, errors.New("sql_query is required")
	}
	sqlArgs := args.GetAttr("sql_args")
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
			return "", nil, errors.New("sql_args must be a list of strings, numbers, or booleans")
		}
	}
	return sqlQuery.AsString(), argsResult, nil

}
func (p Plugin) Call(args plugininterface.Args) plugininterface.Result {
	dbURI, err := p.parseConfig(args.Config)
	if err != nil {
		return plugininterface.Result{
			Diags: hcl.Diagnostics{
				{
					Severity: hcl.DiagError,
					Summary:  "Invalid configuration",
					Detail:   err.Error(),
				},
			},
		}
	}
	sqlQuery, sqlArgs, err := p.parseArgs(args.Args)
	if err != nil {
		return plugininterface.Result{
			Diags: hcl.Diagnostics{
				{
					Severity: hcl.DiagError,
					Summary:  "Invalid arguments",
					Detail:   err.Error(),
				},
			},
		}
	}

	db, err := sql.Open("sqlite3", dbURI)
	if err != nil {
		return plugininterface.Result{
			Diags: hcl.Diagnostics{
				{
					Severity: hcl.DiagError,
					Summary:  "Failed to open database",
					Detail:   err.Error(),
				},
			},
		}
	}
	defer db.Close()
	rows, err := db.Query(sqlQuery, sqlArgs...)
	if err != nil {
		return plugininterface.Result{
			Diags: hcl.Diagnostics{
				{
					Severity: hcl.DiagError,
					Summary:  "Failed to query database",
					Detail:   err.Error(),
				},
			},
		}
	}
	// read columns
	columns, err := rows.Columns()
	if err != nil {
		return plugininterface.Result{
			Diags: hcl.Diagnostics{
				{
					Severity: hcl.DiagError,
					Summary:  "Failed to get column names",
					Detail:   err.Error(),
				},
			},
		}
	}
	result := make([]map[string]any, 0)
	// read rows
	for rows.Next() {
		// create a map of column name to column value
		columnValArr := make([]any, len(columns))
		columnPtrArr := make([]any, len(columns))
		for i := range columns {
			columnPtrArr[i] = &columnValArr[i]
		}
		err = rows.Scan(columnPtrArr...)
		if err != nil {
			return plugininterface.Result{
				Diags: hcl.Diagnostics{
					{
						Severity: hcl.DiagError,
						Summary:  "Failed to scan row",
						Detail:   err.Error(),
					},
				},
			}
		}
		row := make(map[string]any)
		for i, column := range columns {
			row[column] = columnValArr[i]
		}
		result = append(result, row)
	}
	if err := rows.Err(); err != nil {
		return plugininterface.Result{
			Diags: hcl.Diagnostics{
				{
					Severity: hcl.DiagError,
					Summary:  "Failed to read rows",
					Detail:   err.Error(),
				},
			},
		}
	}
	return plugininterface.Result{
		Result: result,
	}
}
