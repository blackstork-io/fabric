package json

import (
	"encoding/json"
	"path/filepath"
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_readFS_Structures(t *testing.T) {
	// test cases
	tt := []struct {
		name     string
		files    []string
		dirs     []string
		pattern  string
		expected []JSONDocument
	}{
		{
			name:     "empty",
			files:    []string{},
			dirs:     []string{},
			pattern:  "*.json",
			expected: []JSONDocument{},
		},
		{
			name:    "one_file",
			files:   []string{"a.json"},
			dirs:    []string{},
			pattern: "*.json",
			expected: []JSONDocument{
				{
					Filename: "a.json",
					Contents: testJSON(map[string]any{
						"property_for": "a.json",
					}),
				},
			},
		},
		{
			name:    "two_files",
			files:   []string{"a.json", "b.json"},
			dirs:    []string{},
			pattern: "*.json",
			expected: []JSONDocument{
				{
					Filename: "a.json",
					Contents: testJSON(map[string]any{
						"property_for": "a.json",
					}),
				},
				{
					Filename: "b.json",
					Contents: testJSON(map[string]any{
						"property_for": "b.json",
					}),
				},
			},
		},
		{
			name:    "one_file_in_one_dir",
			files:   []string{"dir/a.json"},
			dirs:    []string{"dir"},
			pattern: "dir/*.json",
			expected: []JSONDocument{
				{
					Filename: "dir/a.json",
					Contents: testJSON(map[string]any{
						"property_for": "dir/a.json",
					}),
				},
			},
		},
		{
			name:    "one_file_in_two_dirs",
			files:   []string{"dir1/a.json", "dir2/a.json"},
			dirs:    []string{"dir1", "dir2"},
			pattern: "*/a.json",
			expected: []JSONDocument{
				{
					Filename: "dir1/a.json",
					Contents: testJSON(map[string]any{
						"property_for": "dir1/a.json",
					}),
				},
				{
					Filename: "dir2/a.json",
					Contents: testJSON(map[string]any{
						"property_for": "dir2/a.json",
					}),
				},
			},
		},
		{
			name:    "two_files_in_one_dir",
			files:   []string{"dir/a.json", "dir/b.json"},
			dirs:    []string{"dir"},
			pattern: "dir/*.json",
			expected: []JSONDocument{
				{
					Filename: "dir/a.json",
					Contents: testJSON(map[string]any{
						"property_for": "dir/a.json",
					}),
				},
				{
					Filename: "dir/b.json",
					Contents: testJSON(map[string]any{
						"property_for": "dir/b.json",
					}),
				},
			},
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			tmpfs := makeTestFS(t)
			for _, file := range tc.files {
				data := testJSON(map[string]any{
					"property_for": file,
				})
				assert.NoError(t, tmpfs.MkdirAll(filepath.Dir(file), 0o764))
				assert.NoError(t, tmpfs.WriteFile(file, []byte(data), 0o654))
			}
			for _, dir := range tc.dirs {
				assert.NoError(t, tmpfs.MkdirAll(dir, 0o764))
			}
			result, err := readFS(tmpfs, tc.pattern)
			assert.NoError(t, err)
			assert.Equal(t, tc.expected, result)
		})
	}
}

func Test_readFS_Errors(t *testing.T) {
	tt := []struct {
		name      string
		files     []string
		dirs      []string
		pattern   string
		preparefn func(filename string) json.RawMessage
	}{
		{
			name:    "pattern_error",
			files:   []string{},
			dirs:    []string{},
			pattern: "[",
		},
		{
			name:    "file_error",
			files:   []string{"a.json"},
			dirs:    []string{},
			pattern: "*.json",
			preparefn: func(filename string) json.RawMessage {
				return []byte("not json")
			},
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			tmpfs := makeTestFS(t)
			for _, file := range tc.files {
				data := tc.preparefn(file)
				assert.NoError(t, tmpfs.MkdirAll(filepath.Dir(file), 0o764))
				assert.NoError(t, tmpfs.WriteFile(file, []byte(data), 0o654))
			}
			for _, dir := range tc.dirs {
				assert.NoError(t, tmpfs.MkdirAll(dir, 0o764))
			}
			_, err := readFS(tmpfs, tc.pattern)
			assert.Error(t, err)
		})
	}
}

func TestJSONDocument_Map(t *testing.T) {
	type fields struct {
		Filename string
		Contents json.RawMessage
	}
	tt := []struct {
		name   string
		fields fields
		want   map[string]any
	}{
		{
			name: "simple",
			fields: fields{
				Filename: "a.json",
				Contents: testJSON(map[string]any{
					"property_for": "a.json",
				}),
			},
			want: map[string]any{
				"filename": "a.json",
				"contents": map[string]any{
					"property_for": "a.json",
				},
			},
		},
		{
			name: "complex",
			fields: fields{
				Filename: "a.json",
				Contents: testJSON(map[string]any{
					"property_for": "a.json",
					"nested": map[string]any{
						"property_for": "a.json",
						"nested": map[string]any{
							"property_for": "a.json",
						},
					},
				}),
			},
			want: map[string]any{
				"filename": "a.json",
				"contents": map[string]any{
					"property_for": "a.json",
					"nested": map[string]any{
						"property_for": "a.json",
						"nested": map[string]any{
							"property_for": "a.json",
						},
					},
				},
			},
		},
		{
			name: "array",
			fields: fields{
				Filename: "a.json",
				Contents: testJSON([]any{
					map[string]any{
						"id":           float64(0),
						"property_for": "a.json",
					},
					map[string]any{
						"id":           float64(1),
						"property_for": "a.json",
					},
				}),
			},
			want: map[string]any{
				"filename": "a.json",
				"contents": []any{
					map[string]any{
						"id":           float64(0),
						"property_for": "a.json",
					},
					map[string]any{
						"id":           float64(1),
						"property_for": "a.json",
					},
				},
			},
		},
	}
	for _, tt := range tt {
		t.Run(tt.name, func(t *testing.T) {
			doc := JSONDocument{
				Filename: tt.fields.Filename,
				Contents: tt.fields.Contents,
			}
			if got := doc.Map(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("JSONDocument.Map() = %v, want %v", got, tt.want)
			}
		})
	}
}
