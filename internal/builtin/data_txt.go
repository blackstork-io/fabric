package builtin

import (
	"context"
	"io"
	"os"
	"path/filepath"

	"github.com/hashicorp/hcl/v2"
	"github.com/zclconf/go-cty/cty"

	"github.com/blackstork-io/fabric/plugin"
	"github.com/blackstork-io/fabric/plugin/dataspec"
)

func makeTXTDataSource() *plugin.DataSource {
	return &plugin.DataSource{
		DataFunc: fetchTXTData,
		Args: dataspec.ObjectSpec{
			&dataspec.AttrSpec{
				Name:       "glob",
				Type:       cty.String,
				Required:   false,
				ExampleVal: cty.StringVal("path/to/files*.txt"),
				Doc:        `A glob pattern to select TXT files for reading`,
			},
			&dataspec.AttrSpec{
				Name:       "path",
				Type:       cty.String,
				Required:   false,
				ExampleVal: cty.StringVal("data/disclaimer.txt"),
				Doc:        `A file path to a TXT file to read`,
			},
		},
		Doc: `
			Loads TXT files with the names that match a provided "glob" pattern or a single file from a provided path.

			Either "glob" value or "path" value must be provided.

			When "path" is specified, only the content of the file is returned.
			When "glob" is specified, the structure returned by the data source is a list of dicts with file data, for example:
			` + "```json" + `
			  [
			    {
			      "file_path": "path/file-a.txt",
			      "file_name": "file-a.txt",
			      "content": "foobar"
			    },
			    {
			      "file_path": "path/file-b.txt",
			      "file_name": "file-b.txt",
			      "content": "x\\ny\\nz"
			    }
			  ]
			` + "```",
	}
}

func fetchTXTData(ctx context.Context, params *plugin.RetrieveDataParams) (plugin.Data, hcl.Diagnostics) {
	path := params.Args.GetAttr("path")
	if path.IsNull() || path.AsString() == "" {
		return nil, hcl.Diagnostics{{
			Severity: hcl.DiagError,
			Summary:  "Failed to parse arguments",
			Detail:   "path is required",
		}}
	}
	f, err := os.Open(path.AsString())
	if err != nil {
		return nil, hcl.Diagnostics{{
			Severity: hcl.DiagError,
			Summary:  "Failed to open txt file",
			Detail:   err.Error(),
		}}
	}
	defer f.Close()
	data, err := io.ReadAll(f)
	if err != nil {
		return nil, hcl.Diagnostics{{
			Severity: hcl.DiagError,
			Summary:  "Failed to read txt file",
			Detail:   err.Error(),
		}}
	}
	return plugin.StringData(string(data)), nil
}


func readTXTFiles(ctx context.Context, pattern string, sep rune) (plugin.ListData, error) {
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
