package builtin

import (
	"context"
	"io"
	"log/slog"
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

			Either "glob" or "path" attribute must be set.

			When "path" is specified, only the content of the file is returned.
			When "glob" attribute is specified, the structure returned by the data source is a list of dicts that contain the content of the file and file metadata. For example:

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

func readTXTFile(path string) (plugin.Data, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, hcl.Diagnostics{{
			Severity: hcl.DiagError,
			Summary:  "Failed to open a file",
			Detail:   err.Error(),
		}}
	}
	defer f.Close()
	data, err := io.ReadAll(f)
	if err != nil {
		return nil, hcl.Diagnostics{{
			Severity: hcl.DiagError,
			Summary:  "Failed to read a file",
			Detail:   err.Error(),
		}}
	}
	return plugin.StringData(string(data)), nil
}

func readTXTFiles(ctx context.Context, pattern string) (plugin.Data, error) {
	paths, err := filepath.Glob(pattern)
	if err != nil {
		return nil, err
	}
	result := make(plugin.ListData, 0, len(paths))
	for _, path := range paths {
		fileData, err := readTXTFile(path)
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

func fetchTXTData(ctx context.Context, params *plugin.RetrieveDataParams) (plugin.Data, hcl.Diagnostics) {

	glob := params.Args.GetAttr("glob")
	path := params.Args.GetAttr("path")

	if !(path.IsNull() || path.AsString() == "") {
		slog.Debug("Reading a file from the path", "path", path.AsString())
		data, err := readTXTFile(path.AsString())
		if err != nil {
			slog.Error(
				"Error while reading a file",
				slog.String("path", path.AsString()),
				slog.Any("error", err),
			)
			return nil, hcl.Diagnostics{{
				Severity: hcl.DiagError,
				Summary:  "Failed to read a file",
				Detail:   err.Error(),
			}}
		}
		return data, nil
	} else if !glob.IsNull() && glob.AsString() != "" {
		slog.Debug("Reading the files that match the glob pattern", "glob", glob.AsString())
		data, err := readTXTFiles(ctx, glob.AsString())
		if err != nil {
			slog.Error(
				"Error while reading the files",
				slog.String("glob", glob.AsString()),
				slog.Any("error", err),
			)
			return nil, hcl.Diagnostics{{
				Severity: hcl.DiagError,
				Summary:  "Failed to read the files",
				Detail:   err.Error(),
			}}
		}
		return data, nil
	}
	slog.Error("Either \"glob\" value or \"path\" value must be provided")
	return nil, hcl.Diagnostics{{
		Severity: hcl.DiagError,
		Summary:  "Failed to parse provided arguments",
		Detail:   "Either \"glob\" value or \"path\" value must be provided",
	}}
}
