package builtin

import (
	"context"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/zclconf/go-cty/cty"

	"github.com/blackstork-io/fabric/plugin"
)

func Test_makeLocalFilePublisher(t *testing.T) {
	schema := makeLocalFilePublisher()
	assert.NotNil(t, schema.Doc)
	assert.NotNil(t, schema.Tags)
	assert.NotNil(t, schema.Args)
	assert.Equal(t, []plugin.OutputFormat{
		plugin.OutputFormatMD,
		plugin.OutputFormatHTML,
		plugin.OutputFormatPDF,
	}, schema.AllowedFormats)
	assert.NotNil(t, schema.PublishFunc)
}

func Test_publishLocalFileMD(t *testing.T) {
	dir := t.TempDir()
	titleMeta := plugin.MapData{
		"provider": plugin.StringData("title"),
		"plugin":   plugin.StringData("blackstork/builtin"),
	}
	params := &plugin.PublishParams{
		Format: plugin.OutputFormatMD,
		Args: cty.ObjectVal(map[string]cty.Value{
			"path": cty.StringVal(filepath.Join(dir, "{{.document.meta.name}}.{{.format}}")),
		}),
		DataContext: plugin.MapData{
			"document": plugin.MapData{
				"meta": plugin.MapData{
					"name": plugin.StringData("test_document"),
				},
				"content": plugin.MapData{
					"type": plugin.StringData("section"),
					"children": plugin.ListData{
						plugin.MapData{
							"type":     plugin.StringData("element"),
							"markdown": plugin.StringData("# Header 1"),
							"meta":     titleMeta,
						},
						plugin.MapData{
							"type":     plugin.StringData("element"),
							"markdown": plugin.StringData("Lorem ipsum dolor sit amet, consectetur adipiscing elit."),
						},
					},
				},
			},
		},
	}
	diags := publishLocalFile(context.Background(), params)
	require.Empty(t, diags)
	bytes, err := os.ReadFile(filepath.Join(dir, "test_document.md"))
	require.NoError(t, err)
	got := string(bytes)
	assert.Equal(t, got, "# Header 1\n\nLorem ipsum dolor sit amet, consectetur adipiscing elit.")
}

func Test_publishLocalFileHTML(t *testing.T) {
	dir := t.TempDir()
	titleMeta := plugin.MapData{
		"provider": plugin.StringData("title"),
		"plugin":   plugin.StringData("blackstork/builtin"),
	}
	params := &plugin.PublishParams{
		Format: plugin.OutputFormatHTML,
		Args: cty.ObjectVal(map[string]cty.Value{
			"path": cty.StringVal(filepath.Join(dir, "{{.document.meta.name}}.{{.format}}")),
		}),
		DataContext: plugin.MapData{
			"document": plugin.MapData{
				"meta": plugin.MapData{
					"name": plugin.StringData("test_document"),
				},
				"content": plugin.MapData{
					"type": plugin.StringData("section"),
					"children": plugin.ListData{
						plugin.MapData{
							"type":     plugin.StringData("element"),
							"markdown": plugin.StringData("# Header 1"),
							"meta":     titleMeta,
						},
						plugin.MapData{
							"type":     plugin.StringData("element"),
							"markdown": plugin.StringData("Lorem ipsum dolor sit amet, consectetur adipiscing elit."),
						},
					},
				},
			},
		},
	}
	diags := publishLocalFile(context.Background(), params)
	require.Empty(t, diags)
	bytes, err := os.ReadFile(filepath.Join(dir, "test_document.html"))
	require.NoError(t, err)
	got := string(bytes)
	assert.Contains(t, got, "<h1 id=\"header-1\">Header 1</h1>")
	assert.Contains(t, got, "<p>Lorem ipsum dolor sit amet, consectetur adipiscing elit.</p>")
}

func Test_publishLocalFile_invalidPath(t *testing.T) {
	params := &plugin.PublishParams{
		Format: plugin.OutputFormatMD,
		Args: cty.ObjectVal(map[string]cty.Value{
			"path": cty.StringVal(""),
		}),
		DataContext: plugin.MapData{
			"document": plugin.MapData{
				"content": plugin.MapData{
					"type": plugin.StringData("section"),
					"children": plugin.ListData{
						plugin.MapData{
							"type":     plugin.StringData("element"),
							"markdown": plugin.StringData("# Header 1"),
						},
					},
				},
			},
		},
	}
	diags := publishLocalFile(context.Background(), params)
	require.NotEmpty(t, diags)
}
