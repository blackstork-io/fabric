package builtin

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/hashicorp/hcl/v2"
	"github.com/zclconf/go-cty/cty"

	"github.com/blackstork-io/fabric/plugin"
	"github.com/blackstork-io/fabric/plugin/dataspec"
)

func makeJSONDataSource() *plugin.DataSource {
	return &plugin.DataSource{
		DataFunc: fetchJSONData,
		Args: dataspec.ObjectSpec{
			&dataspec.AttrSpec{
				Name:     "glob",
				Type:     cty.String,
				Required: true,
			},
		},
	}
}

func fetchJSONData(ctx context.Context, params *plugin.RetrieveDataParams) (plugin.Data, hcl.Diagnostics) {
	glob := params.Args.GetAttr("glob")
	if glob.IsNull() || glob.AsString() == "" {
		return nil, hcl.Diagnostics{{
			Severity: hcl.DiagError,
			Summary:  "Failed to parse arguments",
			Detail:   "glob is required",
		}}
	}
	data, err := readJSONFiles(ctx, glob.AsString())
	if err != nil {
		return nil, hcl.Diagnostics{{
			Severity: hcl.DiagError,
			Summary:  "Failed to read json files",
			Detail:   err.Error(),
		}}
	}
	return data, nil
}

func readJSONFiles(ctx context.Context, pattern string) (plugin.ListData, error) {
	matchers, err := filepath.Glob(pattern)
	if err != nil {
		return nil, err
	}
	result := make(plugin.ListData, 0, len(matchers))
	for _, matcher := range matchers {
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		default:
			file, err := os.Open(matcher)
			if err != nil {
				return nil, err
			}
			var contents jsonData
			err = json.NewDecoder(file).Decode(&contents)
			if err != nil {
				file.Close()
				return nil, err
			}
			result = append(result, plugin.MapData{
				"filename": plugin.StringData(matcher),
				"contents": contents.data,
			})

			err = file.Close()
			if err != nil {
				return nil, err
			}
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
