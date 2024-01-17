package csv

import (
	"encoding/csv"
	"encoding/json"
	"io/fs"
)

func readFS(filesystem fs.FS, path string, sep rune) ([]map[string]any, error) {
	f, err := filesystem.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	r := csv.NewReader(f)
	r.Comma = sep
	records, err := r.ReadAll()
	if err != nil {
		return nil, err
	}
	if len(records) == 0 {
		return []map[string]any{}, nil
	}
	result := make([]map[string]any, len(records)-1)
	headers := records[0]
	for i, record := range records[1:] {
		result[i] = make(map[string]any, len(headers))
		for j, header := range headers {
			if header == "" {
				continue
			}
			if j >= len(record) {
				result[i][header] = nil
				continue
			}
			if record[j] == "true" {
				result[i][header] = true
			} else if record[j] == "false" {
				result[i][header] = false
			} else {
				n := json.Number(record[j])
				if e, err := n.Int64(); err == nil {
					result[i][header] = e
				} else if f, err := n.Float64(); err == nil {
					result[i][header] = f
				} else {
					result[i][header] = record[j]
				}
			}
		}
	}
	return result, nil
}
