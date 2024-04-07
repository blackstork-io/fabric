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
)

const defaultCSVDelimiter = ','

func makeCSVDataSource() *plugin.DataSource {
	return &plugin.DataSource{
		DataFunc: fetchCSVData,
		Config: dataspec.ObjectSpec{
			&dataspec.AttrSpec{
				Name:     "delimiter",
				Type:     cty.String,
				Required: false,
			},
		},
		Args: dataspec.ObjectSpec{
			&dataspec.AttrSpec{
				Name:     "path",
				Type:     cty.String,
				Required: true,
			},
		},
	}
}

func getDelim(config cty.Value) (r rune, diags hcl.Diagnostics) {
	r = defaultCSVDelimiter
	if config.IsNull() {
		return
	}
	delim := config.GetAttr("delimiter")
	if delim.IsNull() {
		return
	}
	delimStr := delim.AsString()
	delimRune, runeLen := utf8.DecodeRuneInString(delimStr)
	if runeLen == 0 || len(delimStr) != runeLen {
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
	path := params.Args.GetAttr("path")
	if path.IsNull() || path.AsString() == "" {
		return nil, hcl.Diagnostics{{
			Severity: hcl.DiagError,
			Summary:  "path is required",
		}}
	}
	delim, diags := getDelim(params.Config)
	if diags != nil {
		return nil, diags
	}
	data, err := readCSVFile(ctx, path.AsString(), delim)
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
