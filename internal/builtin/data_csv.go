package builtin

import (
	"context"
	"encoding/csv"
	"log/slog"
	"os"
	"path/filepath"
	"unicode/utf8"

	"github.com/hashicorp/hcl/v2"
	"github.com/zclconf/go-cty/cty"

	"github.com/blackstork-io/fabric/internal/builtin/utils"
	"github.com/blackstork-io/fabric/pkg/diagnostics"
	"github.com/blackstork-io/fabric/plugin"
	"github.com/blackstork-io/fabric/plugin/dataspec"
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
				Name:       "glob",
				Type:       cty.String,
				ExampleVal: cty.StringVal("path/to/file*.csv"),
				Doc:        `A glob pattern to select CSV files to read`,
			},
			&dataspec.AttrSpec{
				Name:       "path",
				Type:       cty.String,
				ExampleVal: cty.StringVal("path/to/file.csv"),
				Doc:        `A file path to a CSV file to read`,
			},
		},
		Doc: `
		Loads CSV files with the names that match a provided ` + "`glob`" + ` pattern or a single file from a provided path.

		Either ` + "`glob` or `path` attribute must be set." + `

		When ` + "`path`" + ` attribute is specified, the data source returns only the content of a file.
		When ` + "`glob`" + ` attribute is specified, the data source returns a list of dicts that contain the content of a file and file's metadata.

		**Note**: the data source assumes that CSV file has a header: the data source turns each line into a map with the column titles as keys.

		For example, CSV file with the following data:

		| column_A | column_B | column_C |
		| -------- | -------- | -------- |
		| Foo      | true     | 42       |
		| Bar      | false    | 4.2      |

		will be represented as the following data structure:
		` + "```json" + `
		[
		  {"column_A": "Foo", "column_B": true, "column_C": 42},
		  {"column_A": "Bar", "column_B": false, "column_C": 4.2}
		]
		` + "```" + `

		When ` + "`glob`" + ` is used and multiple files match the pattern, the data source will return a list of dicts, for example:

		` + "```json" + `
		[
		  {
		    "file_path": "path/file-a.csv",
		    "file_name": "file-a.csv",
		    "content": [
		      {"column_A": "Foo", "column_B": true, "column_C": 42},
		      {"column_A": "Bar", "column_B": false, "column_C": 4.2}
		    ]
		  },
		  {
		    "file_path": "path/file-b.csv",
		    "file_name": "file-b.csv",
		    "content": [
		      {"column_C": "Baz", "column_D": 1},
		      {"column_C": "Clu", "column_D": 2}
		    ]
		  },
		]
		` + "```",
	}
}

func getDelim(config cty.Value) (r rune, diags diagnostics.Diag) {
	delim := config.GetAttr("delimiter").AsString()
	delimRune, runeLen := utf8.DecodeRuneInString(delim)
	if runeLen == 0 || len(delim) != runeLen {
		diags = diagnostics.Diag{{
			Severity: hcl.DiagError,
			Summary:  "delimiter must be a single character",
		}}
		return
	}
	r = delimRune
	return
}

func fetchCSVData(ctx context.Context, params *plugin.RetrieveDataParams) (plugin.Data, diagnostics.Diag) {
	glob := params.Args.GetAttr("glob")
	path := params.Args.GetAttr("path")

	delim, err := getDelim(params.Config)
	if err != nil {
		slog.Error("Error while getting a delimiter value", slog.Any("error", err))
		return nil, err
	}

	if !(path.IsNull() || path.AsString() == "") {
		slog.Debug("Reading a file from a path", "path", path.AsString())
		data, err := readAndDecodeCSVFile(ctx, path.AsString(), delim)
		if err != nil {
			slog.Error(
				"Error while reading a CSV file",
				slog.String("path", path.AsString()),
				slog.Any("error", err),
			)
			return nil, diagnostics.Diag{{
				Severity: hcl.DiagError,
				Summary:  "Failed to read a file",
				Detail:   err.Error(),
			}}
		}
		return data, nil
	} else if !glob.IsNull() && glob.AsString() != "" {
		slog.Debug("Reading the files that match a glob pattern", "glob", glob.AsString())
		data, err := readCSVFiles(ctx, glob.AsString(), delim)
		if err != nil {
			slog.Error(
				"Error while reading the CSV files",
				slog.String("glob", glob.AsString()),
				slog.Any("error", err),
			)
			return nil, diagnostics.Diag{{
				Severity: hcl.DiagError,
				Summary:  "Failed to read the files",
				Detail:   err.Error(),
			}}
		}
		return data, nil
	}
	slog.Error("Either \"glob\" value or \"path\" value must be provided")
	return nil, diagnostics.Diag{{
		Severity: hcl.DiagError,
		Summary:  "Failed to parse provided arguments",
		Detail:   "Either \"glob\" value or \"path\" value must be provided",
	}}
}

func readCSVFiles(ctx context.Context, pattern string, delimiter rune) (plugin.ListData, error) {
	paths, err := filepath.Glob(pattern)
	if err != nil {
		return nil, err
	}
	result := make(plugin.ListData, 0, len(paths))
	for _, path := range paths {
		fileData, err := readAndDecodeCSVFile(ctx, path, delimiter)
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

func readAndDecodeCSVFile(ctx context.Context, path string, delimiter rune) (plugin.ListData, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	reader := csv.NewReader(f)
	reader.Comma = delimiter

	return utils.ParseCSVContent(ctx, reader)
}
