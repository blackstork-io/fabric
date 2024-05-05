package builtin

import (
	"context"
	"encoding/csv"
	"encoding/json"
	"io"
	"os"
	"path/filepath"
	"unicode/utf8"

	"github.com/hashicorp/hcl/v2"
	"github.com/zclconf/go-cty/cty"

	"github.com/blackstork-io/fabric/plugin"
	"github.com/blackstork-io/fabric/plugin/dataspec"
)

func makeCSVDataSource() *plugin.DataSource {
	return &plugin.DataSource{
		DataFunc: fetchCSVData,
		Config: dataspec.ObjectSpec{
			&dataspec.AttrSpec{
				Name:       "delimiter",
				Type:       cty.String,
				Required:   false,
				DefaultVal: cty.StringVal(","),
				Doc:        `Must be a one-character string`,
			},
		},
		Args: dataspec.ObjectSpec{
			&dataspec.AttrSpec{
				Name:       "glob",
				Type:       cty.String,
				Required:   true,
				ExampleVal: cty.StringVal("path/to/files*.csv"),
				Doc:        `A glob pattern to select CSV files for reading`,
			},
		},
		Doc: `
			Loads CSV files with the names that match a provided "glob" pattern.

			We assume that every CSV file has a header. Each line of the CSV file is converted into a map, with keys that match the header titles.

			For example, the following CSV data

			| column_A | column_B | column_C |
			| -------- | -------- | -------- |
			| Test     | true     | 42       |
			| Line 2   | false    | 4.2      |

			will be represented as this data structure:
			` + "```json" + `
			  [
			    {
			      "file_path": "path/file-a.csv",
			      "file_name": "file-a.csv",
			      "content": [
			        {"column_A": "Test", "column_B": true, "column_C": 42},
			        {"column_A": "Line 2", "column_B": false, "column_C": 4.2}
			      ]
			    }
			  ]
			` + "```",
	}
}

func getDelim(config cty.Value) (r rune, diags hcl.Diagnostics) {
	delim := config.GetAttr("delimiter").AsString()
	delimRune, runeLen := utf8.DecodeRuneInString(delim)
	if runeLen == 0 || len(delim) != runeLen {
		diags = hcl.Diagnostics{{
			Severity: hcl.DiagError,
			Summary:  "delimiter must be a single character",
		}}
		return
	}
	r = delimRune
	return
}

func fetchCSVData(ctx context.Context, params *plugin.RetrieveDataParams) (plugin.Data, hcl.Diagnostics) {
	glob := params.Args.GetAttr("glob")
	if glob.IsNull() || glob.AsString() == "" {
		return nil, hcl.Diagnostics{{
			Severity: hcl.DiagError,
			Summary:  "Failed to parse arguments",
			Detail:   "glob value is required",
		}}
	}
	delim, diags := getDelim(params.Config)
	if diags != nil {
		return nil, diags
	}
	data, err := readCSVFiles(ctx, glob.AsString(), delim)
	if err != nil {
		return nil, hcl.Diagnostics{{
			Severity: hcl.DiagError,
			Summary:  "Failed to read CSV files",
			Detail:   err.Error(),
		}}
	}
	return data, nil
}

func readCSVFiles(ctx context.Context, pattern string, sep rune) (plugin.ListData, error) {
	paths, err := filepath.Glob(pattern)
	if err != nil {
		return nil, err
	}
	result := make(plugin.ListData, 0, len(paths))
	for _, path := range paths {
		fileData, err := readCSVFile(ctx, path, sep)
		if err != nil {
			return result, err
		}
		result = append(result, plugin.MapData{
			"file_path": plugin.StringData(path),
			"file_name": plugin.StringData(filepath.Base(path)),
			"content":   fileData,
		})
	}
	return result, nil
}

func readCSVFile(ctx context.Context, path string, sep rune) (plugin.ListData, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	rowMaps := make(plugin.ListData, 0)

	reader := csv.NewReader(f)
	reader.Comma = sep

	headers, err := reader.Read()
	if err == io.EOF {
		return rowMaps, nil
	} else if err != nil {
		return nil, err
	}

	for {
		select {
		case <-ctx.Done(): // stop reading if the context is canceled
			return nil, ctx.Err()
		default:
			row, err := reader.Read()
			if err == io.EOF {
				return rowMaps, nil
			} else if err != nil {
				return nil, err
			}
			rowMap := make(plugin.MapData, len(headers))
			for j, header := range headers {
				if header == "" {
					continue
				}
				if j >= len(row) {
					rowMap[header] = nil
					continue
				}
				if row[j] == "true" {
					rowMap[header] = plugin.BoolData(true)
				} else if row[j] == "false" {
					rowMap[header] = plugin.BoolData(false)
				} else {
					n := json.Number(row[j])
					if f, err := n.Float64(); err == nil {
						rowMap[header] = plugin.NumberData(f)
					} else {
						rowMap[header] = plugin.StringData(row[j])
					}
				}
			}
			rowMaps = append(rowMaps, rowMap)
		}
	}
}
