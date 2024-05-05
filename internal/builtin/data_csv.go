package builtin

import (
	"context"
	"encoding/csv"
	"encoding/json"
	"io"
	"os"
	"unicode/utf8"

	"github.com/hashicorp/hcl/v2"
	"github.com/zclconf/go-cty/cty"

	"github.com/blackstork-io/fabric/plugin"
	"github.com/blackstork-io/fabric/plugin/dataspec"
	"github.com/blackstork-io/fabric/plugin/dataspec/constraint"
)

func makeCSVDataSource() *plugin.DataSource {
	return &plugin.DataSource{
		DataFunc: fetchCSVData,
		Config: dataspec.ObjectSpec{
			&dataspec.AttrSpec{
				Name:         "delimiter",
				Type:         cty.String,
				DefaultVal:   cty.StringVal(","),
				MinInclusive: cty.NumberIntVal(1),
				MaxInclusive: cty.NumberIntVal(1),
				Doc:          `CSV field delimiter`,
			},
		},
		Args: dataspec.ObjectSpec{
			&dataspec.AttrSpec{
				Name:        "path",
				Type:        cty.String,
				Constraints: constraint.RequiredMeaningfull,
				ExampleVal:  cty.StringVal("path/to/file.csv"),
			},
		},
		Doc: `
		Imports and parses a csv file.

		We assume the table has a header and turn each line into a map based on the header titles.

		For example following table

		| column_A | column_B | column_C |
		| -------- | -------- | -------- |
		| Test     | true     | 42       |
		| Line 2   | false    | 4.2      |

		will be represented as the following structure:
		` + "```json" + `
		[
		  {"column_A": "Test", "column_B": true, "column_C": 42},
		  {"column_A": "Line 2", "column_B": false, "column_C": 4.2}
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
	path := params.Args.GetAttr("path").AsString()
	if path == "" {
		return nil, hcl.Diagnostics{{
			Severity: hcl.DiagError,
			Summary:  "path is required",
		}}
	}
	delim, diags := getDelim(params.Config)
	if diags != nil {
		return nil, diags
	}
	data, err := readCSVFile(ctx, path, delim)
	if err != nil {
		return nil, hcl.Diagnostics{{
			Severity: hcl.DiagError,
			Summary:  "Failed to read csv file",
			Detail:   err.Error(),
		}}
	}
	return data, nil
}

func readCSVFile(ctx context.Context, path string, sep rune) (plugin.ListData, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	r := csv.NewReader(f)
	r.Comma = sep
	result := make(plugin.ListData, 0)
	headers, err := r.Read()
	if err == io.EOF {
		return result, nil
	} else if err != nil {
		return nil, err
	}
	for {
		select {
		case <-ctx.Done(): // stop reading if the context is canceled
			return nil, ctx.Err()
		default:
			row, err := r.Read()
			if err == io.EOF {
				return result, nil
			} else if err != nil {
				return nil, err
			}

			m := make(plugin.MapData, len(headers))
			for j, header := range headers {
				if header == "" {
					continue
				}
				if j >= len(row) {
					m[header] = nil
					continue
				}
				if row[j] == "true" {
					m[header] = plugin.BoolData(true)
				} else if row[j] == "false" {
					m[header] = plugin.BoolData(false)
				} else {
					n := json.Number(row[j])
					if f, err := n.Float64(); err == nil {
						m[header] = plugin.NumberData(f)
					} else {
						m[header] = plugin.StringData(row[j])
					}
				}
			}
			result = append(result, m)
		}
	}
}
