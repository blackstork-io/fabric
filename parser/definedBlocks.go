package parser

import (
	"fmt"

	"github.com/hashicorp/hcl/v2"

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

func (db *DefinedBlocks) GetSection(expr hcl.Expression) (section *definitions.Section, diags diagnostics.Diag) {
	res, diags := db.Traverse(expr)
	if diags.HasErrors() {
		return
	}
	section, ok := res.(*definitions.Section)
	if !ok {
		diags.Append(&hcl.Diagnostic{
			Severity: hcl.DiagError,
			Summary:  "Invalid path",
			Detail:   "This path is not referencing a section block",
			Subject:  expr.Range().Ptr(),
		})
	}
	return
}

func (db *DefinedBlocks) GetPlugin(expr hcl.Expression) (plugin *definitions.Plugin, diags diagnostics.Diag) {
	res, diags := db.Traverse(expr)
	if diags.HasErrors() {
		return
	}
	plugin, ok := res.(*definitions.Plugin)
	if !ok {
		diags.Append(&hcl.Diagnostic{
			Severity: hcl.DiagError,
			Summary:  "Invalid path",
			Detail:   "This path is not referencing a plugin block",
			Subject:  expr.Range().Ptr(),
		})
	}
	return
}

func (db *DefinedBlocks) GetConfig(expr hcl.Expression) (cfg *definitions.Config, diags diagnostics.Diag) {
	res, diags := db.Traverse(expr)
	if diags.HasErrors() {
		return
	}
	cfg, ok := res.(*definitions.Config)
	if !ok {
		diags.Append(&hcl.Diagnostic{
			Severity: hcl.DiagError,
			Summary:  "Invalid path",
			Detail:   "This path is not referencing a config block",
			Subject:  expr.Range().Ptr(),
		})
	}
	return
}

func (db *DefinedBlocks) DefaultConfigFor(plugin *definitions.Plugin) (config *definitions.Config) {
	return db.Config[definitions.Key{
		PluginKind: plugin.Kind(),
		PluginName: plugin.Name(),
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
