package blocks

import (
	"fmt"

	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/hcldec"
	"github.com/hashicorp/hcl/v2/hclsyntax"
	"github.com/zclconf/go-cty/cty"
	"github.com/zclconf/go-cty/cty/convert"

	"github.com/blackstork-io/fabric/pkg/diagnostics"
)

var globalConfigSpec = &hcldec.ObjectSpec{
	"cache_dir": &hcldec.AttrSpec{
		Name:     "cache_dir",
		Type:     cty.String,
		Required: false,
	},
	"plugin_registry": &hcldec.BlockSpec{
		TypeName: "plugin_registry",
		Nested: hcldec.ObjectSpec{
			"mirror_dir": &hcldec.AttrSpec{
				Name:     "mirror_dir",
				Type:     cty.String,
				Required: false,
			},
		},
	},
	"plugin_versions": &hcldec.AttrSpec{
		Name:     "plugin_versions",
		Type:     cty.Map(cty.String),
		Required: false,
	},
}

type GlobalConfig struct {
	block          *hclsyntax.Block
	CacheDir       string
	PluginRegistry *PluginRegistry
	PluginVersions map[string]string
}

type PluginRegistry struct {
	MirrorDir string
}

func DefineGlobalConfig(block *hclsyntax.Block) (cfg *GlobalConfig, diags diagnostics.Diag) {
	if len(block.Labels) > 0 {
		return nil, diagnostics.Diag{{
			Severity: hcl.DiagError,
			Summary:  "Invalid global config",
			Detail:   "Global config should not have labels",
		}}
	}
	value, hclDiags := hcldec.Decode(block.Body, globalConfigSpec, nil)
	if diags.ExtendHcl(hclDiags) {
		return
	}
	typ := hcldec.ImpliedType(globalConfigSpec)
	errs := value.Type().TestConformance(typ)
	if len(errs) > 0 {
		var err error
		value, err = convert.Convert(value, typ)
		if err != nil {
			diags.AppendErr(err, "Error while serializing global config")
			return
		}
	}
	cfg = &GlobalConfig{
		block:          block,
		CacheDir:       "./.fabric",
		PluginVersions: make(map[string]string),
	}
	cacheDir := value.GetAttr("cache_dir")
	if !cacheDir.IsNull() && cacheDir.AsString() != "" {
		cfg.CacheDir = cacheDir.AsString()
	}
	pluginRegistry := value.GetAttr("plugin_registry")
	if !pluginRegistry.IsNull() {
		mirrorDir := pluginRegistry.GetAttr("mirror_dir")
		if !mirrorDir.IsNull() || mirrorDir.AsString() != "" {
			cfg.PluginRegistry = &PluginRegistry{
				MirrorDir: mirrorDir.AsString(),
			}
		}
	}
	pluginVersions := value.GetAttr("plugin_versions")
	if !pluginVersions.IsNull() {
		versionMap := pluginVersions.AsValueMap()
		for k, v := range versionMap {
			if v.Type() != cty.String {
				diags.Append(&hcl.Diagnostic{
					Severity: hcl.DiagError,
					Summary:  "Invalid plugin version",
					Detail:   fmt.Sprintf("Version of plugin '%s' should be a string", k),
				})
				continue
			}
			cfg.PluginVersions[k] = v.AsString()
		}
	}
	return cfg, nil
}
