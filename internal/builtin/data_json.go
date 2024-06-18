package builtin

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"os"
	"path/filepath"

	"github.com/hashicorp/hcl/v2"
	"github.com/zclconf/go-cty/cty"

	"github.com/blackstork-io/fabric/pkg/diagnostics"
	"github.com/blackstork-io/fabric/plugin"
	"github.com/blackstork-io/fabric/plugin/dataspec"
)

func makeJSONDataSource() *plugin.DataSource {
	return &plugin.DataSource{
		DataFunc: fetchJSONData,
		Args: dataspec.ObjectSpec{
			&dataspec.AttrSpec{
				Name:       "glob",
				Type:       cty.String,
				ExampleVal: cty.StringVal("path/to/file*.json"),
				Doc:        `A glob pattern to select JSON files to read`,
			},
			&dataspec.AttrSpec{
				Name:       "path",
				Type:       cty.String,
				ExampleVal: cty.StringVal("path/to/file.json"),
				Doc:        `A file path to a JSON file to read`,
			},
		},
		Doc: `
		Loads JSON files with the names that match a provided ` + "`glob`" + ` pattern or a single file from a provided ` + "`path`" + `value.

		Either ` + "`glob`" + ` or ` + "`path`" + ` argument must be set.

		When ` + "`path`" + ` argument is specified, the data source returns only the content of a file.
		When ` + "`glob`" + ` argument is specified, the data source returns a list of dicts that contain the content of a file and file's metadata. For example:

		` + "```json" + `
		[
		  {
			"file_path": "path/file-a.json",
			"file_name": "file-a.json",
			"content": {
			  "foo": "bar"
			}
		  },
		  {
			"file_path": "path/file-b.json",
			"file_name": "file-b.json",
			"content": [
			  {"x": "y"}
			]
		  }
		]
		` + "```",
	}
}

func fetchJSONData(ctx context.Context, params *plugin.RetrieveDataParams) (plugin.Data, diagnostics.Diag) {
	glob := params.Args.GetAttr("glob")
	path := params.Args.GetAttr("path")

	if !path.IsNull() && path.AsString() != "" {
		slog.Debug("Reading a file from a path", "path", path.AsString())
		data, err := readAndDecodeJSONFile(path.AsString())
		if err != nil {
			slog.Error(
				"Error while reading a JSON file",
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
		data, err := readJSONFiles(ctx, glob.AsString())
		if err != nil {
			slog.Error(
				"Error while reading the JSON files",
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

func readAndDecodeJSONFile(path string) (plugin.Data, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var content jsonData
	err = json.NewDecoder(file).Decode(&content)
	if err != nil {
		file.Close()
		return nil, err
	}
	return content.data, nil
}

func readJSONFiles(ctx context.Context, pattern string) (plugin.ListData, error) {
	paths, err := filepath.Glob(pattern)
	if err != nil {
		return nil, err
	}
	result := make(plugin.ListData, 0, len(paths))
	for _, path := range paths {
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		default:
			content, err := readAndDecodeJSONFile(path)
			if err != nil {
				return result, err
			}
			result = append(result, plugin.MapData{
				"file_path": plugin.StringData(path),
				"file_name": plugin.StringData(filepath.Base(path)),
				"content":   content,
			})
		}
	}
	return result, nil
}

type jsonData struct {
	data plugin.Data
}

func (d jsonData) toData(v any) (res plugin.Data, err error) {
	switch v := v.(type) {
	case nil:
		return nil, nil
	case float64:
		return plugin.NumberData(v), nil
	case string:
		return plugin.StringData(v), nil
	case bool:
		return plugin.BoolData(v), nil
	case map[string]any:
		m := make(plugin.MapData)
		for k, v := range v {
			m[k], err = d.toData(v)
			if err != nil {
				return nil, err
			}
		}
		return m, nil
	case []any:
		l := make(plugin.ListData, len(v))
		for i, v := range v {
			l[i], err = d.toData(v)
			if err != nil {
				return nil, err
			}
		}
		return l, nil
	default:
		return nil, fmt.Errorf("unsupported type %T", v)
	}
}

func (d *jsonData) UnmarshalJSON(b []byte) error {
	if !json.Valid(b) {
		return fmt.Errorf("invalid JSON data")
	}
	var result any
	err := json.Unmarshal(b, &result)
	if err != nil {
		return err
	}
	d.data, err = d.toData(result)
	return err
}
