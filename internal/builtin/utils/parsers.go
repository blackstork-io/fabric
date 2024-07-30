package utils

import (
	"context"
	"encoding/csv"
	"encoding/json"
	"io"

	"github.com/blackstork-io/fabric/plugin/plugindata"
)

func ParseCSVContent(ctx context.Context, reader *csv.Reader) (plugindata.List, error) {
	rowMaps := make(plugindata.List, 0)
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
			rowMap := make(plugindata.Map, len(headers))
			for j, header := range headers {
				if header == "" {
					continue
				}
				if j >= len(row) {
					rowMap[header] = nil
					continue
				}
				if row[j] == "true" {
					rowMap[header] = plugindata.Bool(true)
				} else if row[j] == "false" {
					rowMap[header] = plugindata.Bool(false)
				} else {
					n := json.Number(row[j])
					if f, err := n.Float64(); err == nil {
						rowMap[header] = plugindata.Number(f)
					} else {
						rowMap[header] = plugindata.String(row[j])
					}
				}
			}
			rowMaps = append(rowMaps, rowMap)
		}
	}
}
