package utils

import (
	"context"
	"encoding/csv"
	"encoding/json"
	"io"

	"github.com/blackstork-io/fabric/plugin"
)

func ParseCSVContent(ctx context.Context, reader *csv.Reader) (plugin.ListData, error) {
	rowMaps := make(plugin.ListData, 0)

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
