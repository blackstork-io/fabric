package definitions

type GlobalConfig struct {
	CacheDir       string            `hcl:"cache_dir,optional"`
	PluginRegistry *PluginRegistry   `hcl:"plugin_registry,block"`
	PluginVersions map[string]string `hcl:"plugin_versions,optional"`
}

type PluginRegistry struct {
	BaseURL   string `hcl:"base_url,optional"`
	MirrorDir string `hcl:"mirror_dir,optional"`
}

func (g *GlobalConfig) Merge(other *GlobalConfig) {
	if other.CacheDir != "" {
		g.CacheDir = other.CacheDir
	}
	if other.PluginRegistry != nil {
		if g.PluginRegistry == nil {
			g.PluginRegistry = other.PluginRegistry
		} else {
			if other.PluginRegistry.BaseURL != "" {
				g.PluginRegistry.BaseURL = other.PluginRegistry.BaseURL
			}
			if other.PluginRegistry.MirrorDir != "" {
				g.PluginRegistry.MirrorDir = other.PluginRegistry.MirrorDir
			}
		}
	}
	g.PluginVersions = other.PluginVersions
}
