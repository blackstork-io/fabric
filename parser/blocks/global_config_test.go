package blocks

import (
	"testing"

	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/hclsyntax"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDefineGlobalConfig(t *testing.T) {
	t.Parallel()
	tt := []struct {
		name  string
		block string
		want  *GlobalConfig
	}{
		{
			name:  "empty",
			block: `fabric {}`,
			want: &GlobalConfig{
				CacheDir:       "./.fabric",
				PluginVersions: map[string]string{},
			},
		},
		{
			name: "with_cache_dir",
			block: `fabric {
				  cache_dir = "./.other_cache"
			}`,
			want: &GlobalConfig{
				CacheDir:       "./.other_cache",
				PluginVersions: map[string]string{},
			},
		},
		{
			name: "with_plugin_registry",
			block: `fabric {
				  plugin_registry {
					mirror_dir = "./.other_mirror"
				  }
			}`,
			want: &GlobalConfig{
				CacheDir: "./.fabric",
				PluginRegistry: &PluginRegistry{
					MirrorDir: "./.other_mirror",
				},
				PluginVersions: map[string]string{},
			},
		},
		{
			name: "with_plugin_versions",
			block: `fabric {
				  plugin_versions = {
					"namespace/plugin1" = "1.0.0"
					"namespace/plugin2" = "2.0.0"
				  }
			}`,
			want: &GlobalConfig{
				CacheDir: "./.fabric",
				PluginVersions: map[string]string{
					"namespace/plugin1": "1.0.0",
					"namespace/plugin2": "2.0.0",
				},
			},
		},
		{
			name: "with_all",
			block: `fabric {
				  cache_dir = "./.other_cache"
				  plugin_registry {
					mirror_dir = "./.other_mirror"
				  }
				  plugin_versions = {
					"namespace/plugin1" = "1.0.0"
					"namespace/plugin2" = "2.0.0"
				  }
			}`,
			want: &GlobalConfig{
				CacheDir: "./.other_cache",
				PluginRegistry: &PluginRegistry{
					MirrorDir: "./.other_mirror",
				},
				PluginVersions: map[string]string{
					"namespace/plugin1": "1.0.0",
					"namespace/plugin2": "2.0.0",
				},
			},
		},
	}

	for _, tc := range tt {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			f, hcldiags := hclsyntax.ParseConfig([]byte(tc.block), "", hcl.Pos{})
			require.Len(t, hcldiags, 0)
			body, ok := f.Body.(*hclsyntax.Body)
			require.True(t, ok)
			block := body.Blocks[0]
			got, diags := DefineGlobalConfig(block)
			assert.Len(t, diags, 0)
			tc.want.block = block
			assert.Equal(t, tc.want, got)
		})
	}
}
