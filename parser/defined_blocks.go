package parser

import (
	"fmt"

	"github.com/hashicorp/hcl/v2"
	"github.com/zclconf/go-cty/cty"

	"github.com/blackstork-io/fabric/parser/definitions"
	"github.com/blackstork-io/fabric/pkg/diagnostics"
)

// Collection of defined blocks

type DefinedBlocks struct {
	GlobalConfig *definitions.GlobalConfig
	Config       map[definitions.Key]*definitions.Config
	Documents    map[string]*definitions.Document
	Sections     map[string]*definitions.Section
	Plugins      map[definitions.Key]*definitions.Plugin
}

func mapGetOrInit[K1, K2 comparable, V any](m map[K1]map[K2]V, key K1) (innerMap map[K2]V) {
	innerMap, found := m[key]
	if !found {
		innerMap = map[K2]V{}
		m[key] = innerMap
	}
	return
}

func mapToCty(m map[string]map[string]cty.Value) (res map[string]cty.Value) {
	res = make(map[string]cty.Value, len(m))
	for k, v := range m {
		if len(v) == 0 {
			continue
		}
		res[k] = cty.MapVal(v)
	}
	return
}

func PluginMapToCty[V definitions.FabricBlock](plugins map[definitions.Key]V) (content, data, publish cty.Value) {
	// [plugin_kind][plugin_name][block_name]*definitions.Plugin

	pluginMap := [3]map[string]map[string]cty.Value{
		{},
		{},
		{},
	}
	for k, v := range plugins {
		var idx int
		switch k.PluginKind {
		case definitions.BlockKindContent:
			idx = 0
		case definitions.BlockKindData:
			idx = 1
		case definitions.BlockKindPublish:
			idx = 2
		default:
			panic("must be exhaustive")
		}
		blockNameToVal := mapGetOrInit(pluginMap[idx], k.PluginName)
		blockNameToVal[k.BlockName] = definitions.ToCtyValue(v)
	}
	pluginKindToVal := [3]cty.Value{}

	for idx, pl := range pluginMap {
		if len(pl) == 0 {
			continue
		}
		pluginKindToVal[idx] = cty.MapVal(mapToCty(pl))
	}
	return pluginKindToVal[0], pluginKindToVal[1], pluginKindToVal[2]
}

func (db *DefinedBlocks) AsValueMap() map[string]cty.Value {
	content, data, publish := PluginMapToCty(db.Plugins)
	cfgContent, cfgData, cfgPublish := PluginMapToCty(db.Config)
	config := cty.ObjectVal(map[string]cty.Value{
		definitions.BlockKindContent: cfgContent,
		definitions.BlockKindData:    cfgData,
		definitions.BlockKindPublish: cfgPublish,
	})

	var sections cty.Value
	if len(db.Sections) == 0 {
		sections = cty.MapValEmpty(cty.Map((*definitions.Section)(nil).CtyType()))
	} else {
		sect := make(map[string]cty.Value, len(db.Sections))
		for k, v := range db.Sections {
			sect[k] = definitions.ToCtyValue(v)
		}
		sections = cty.MapVal(sect)
	}
	return map[string]cty.Value{
		definitions.BlockKindContent: content,
		definitions.BlockKindData:    data,
		definitions.BlockKindSection: sections,
		definitions.BlockKindConfig:  config,
		definitions.BlockKindPublish: publish,
	}
}

func (db *DefinedBlocks) DefaultConfigFor(plugin *definitions.Plugin) (config *definitions.Config) {
	return db.DefaultConfig(plugin.Kind(), plugin.Name())
}

func (db *DefinedBlocks) DefaultConfig(pluginKind, pluginName string) (config *definitions.Config) {
	return db.Config[definitions.Key{
		PluginKind: pluginKind,
		PluginName: pluginName,
		BlockName:  "",
	}]
}

func (db *DefinedBlocks) Merge(other *DefinedBlocks) (diags diagnostics.Diag) {
	if other.GlobalConfig != nil {
		if db.GlobalConfig != nil {
			diags.Add("Global config declared multiple times", "")
		}
		db.GlobalConfig = other.GlobalConfig
	}
	for k, v := range other.Config {
		diags.Append(AddIfMissing(db.Config, k, v))
	}
	for k, v := range other.Documents {
		diags.Append(AddIfMissing(db.Documents, k, v))
	}
	for k, v := range other.Sections {
		diags.Append(AddIfMissing(db.Sections, k, v))
	}
	for k, v := range other.Plugins {
		diags.Append(AddIfMissing(db.Plugins, k, v))
	}
	return
}

func AddIfMissing[M ~map[K]V, K comparable, V definitions.FabricBlock](m M, key K, newBlock V) *hcl.Diagnostic {
	if origBlock, found := m[key]; found {
		kind := origBlock.GetHCLBlock().Type
		origDefRange := origBlock.GetHCLBlock().DefRange
		return &hcl.Diagnostic{
			Severity: hcl.DiagError,
			Summary:  fmt.Sprintf("Duplicate '%s' declaration", kind),
			Detail:   fmt.Sprintf("'%s' with the same name originally defined at %s:%d", kind, origDefRange.Filename, origDefRange.Start.Line),
			Subject:  newBlock.GetHCLBlock().DefRange.Ptr(),
		}
	}
	m[key] = newBlock
	return nil
}

func NewDefinedBlocks() *DefinedBlocks {
	return &DefinedBlocks{
		Config:    map[definitions.Key]*definitions.Config{},
		Documents: map[string]*definitions.Document{},
		Sections:  map[string]*definitions.Section{},
		Plugins:   map[definitions.Key]*definitions.Plugin{},
	}
}
