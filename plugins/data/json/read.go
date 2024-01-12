package json

import (
	"encoding/json"
	"io/fs"
)

// JSONDocument represents a JSON document that was read from the filesystem
type JSONDocument struct {
	Filename string          `json:"filename"`
	Contents json.RawMessage `json:"contents"`
}

func (doc JSONDocument) Map() map[string]any {
	var result any
	_ = json.Unmarshal(doc.Contents, &result)
	return map[string]any{
		"filename": doc.Filename,
		"contents": result,
	}
}

// readFS reads all JSON documents from the filesystem that match the given glob pattern
// The pattern is relative to the root of the filesystem
func readFS(filesystem fs.FS, pattern string) ([]JSONDocument, error) {
	matchers, err := fs.Glob(filesystem, pattern)
	if err != nil {
		return nil, err
	}
	result := []JSONDocument{}
	for _, matcher := range matchers {
		file, err := filesystem.Open(matcher)
		if err != nil {
			return nil, err
		}
		var contents json.RawMessage
		err = json.NewDecoder(file).Decode(&contents)
		if err != nil {
			return nil, err
		}
		result = append(result, JSONDocument{
			Filename: matcher,
			Contents: contents,
		})
	}
	return result, nil
}
