package blocks

import (
	"fmt"

	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/hclsyntax"
	"github.com/zclconf/go-cty/cty"

	"github.com/blackstork-io/fabric/parser/blocks/internal/tree"
	"github.com/blackstork-io/fabric/parser/definitions"
	"github.com/blackstork-io/fabric/pkg/diagnostics"
)

// Collection of defined blocks
// TMP
type (
	Document = Config
	Section  = Config
)

type DefinedBlocks struct {
	tree.NodeSigil
	GlobalConfig *GlobalConfig
	Config       nodeMap[pluginKindT, nodeMap[pluginNameT, nodeMap[blockNameT, *Config]]]
	Documents    nodeMap[documentNameT, *Document]
	Sections     nodeMap[sectionNameT, *Section]
	Plugins      nodeMap[pluginKindT, nodeMap[pluginNameT, nodeMap[blockNameT, Plugin]]]
}

// FriendlyName implements tree.Node.
func (db *DefinedBlocks) FriendlyName() string {
	return "root"
}

func (db *DefinedBlocks) AsCtyValue() cty.Value {
	return cty.ObjectVal(map[string]cty.Value{
		definitions.BlockKindContent: db.Plugins.Map[definitions.BlockKindContent].AsCtyValue(),
		definitions.BlockKindData:    db.Plugins.Map[definitions.BlockKindData].AsCtyValue(),
		definitions.BlockKindSection: db.Sections.AsCtyValue(),
		definitions.BlockKindConfig:  db.Config.AsCtyValue(),
	})
}


func (db *DefinedBlocks) DefaultConfig(pluginKind, pluginName string) (config *Config) {
	return db.Config.Get(pluginKind).Get(pluginName).Get("")
}

func (db *DefinedBlocks) Merge(other *DefinedBlocks) (diags diagnostics.Diag) {
	if other.GlobalConfig != nil {
		if db.GlobalConfig != nil {
			diags.Add("Global config declared multiple times", "")
		}
		db.GlobalConfig = other.GlobalConfig
	}
	diags.Extend(MergeNestedMap(db.Config, other.Config))
	diags.Extend(MergeMap(db.Documents, other.Documents))
	diags.Extend(MergeMap(db.Sections, other.Sections))
	diags.Extend(MergeNestedMap(db.Plugins, other.Plugins))
	return
}

func AddIfMissing[M ~map[K]V, K comparable, V FabricBlock](m M, key K, newBlock V) *hcl.Diagnostic {
	if origBlock, found := m[key]; found {
		kind := origBlock.HCLBlock().Type
		origDefRange := origBlock.HCLBlock().DefRange()
		return &hcl.Diagnostic{
			Severity: hcl.DiagError,
			Summary:  fmt.Sprintf("Duplicate '%s' declaration", kind),
			Detail:   fmt.Sprintf("'%s' with the same name originally defined at %s:%d", kind, origDefRange.Filename, origDefRange.Start.Line),
			Subject:  newBlock.HCLBlock().DefRange().Ptr(),
		}
	}
	m[key] = newBlock
	return nil
}

func NewDefinedBlocks() *DefinedBlocks {
	return &DefinedBlocks{
		Config:    NewMap[pluginKindT, nodeMap[pluginNameT, nodeMap[blockNameT, *Config]]](),
		Documents: NewMap[documentNameT, *Document](),
		Sections:  NewMap[sectionNameT, *Section](),
		Plugins:   NewMap[pluginKindT, nodeMap[pluginNameT, nodeMap[blockNameT, Plugin]]](), // map[definitions.Key]*definitions.Plugin{},
	}
}

func NewMap[K key, V val]() nodeMap[K, V] {
	return nodeMap[K, V]{
		Map: map[string]V{},
	}
}

type (
	documentNameT string
	pluginKindT   string
	pluginNameT   string
	blockNameT    string
	sectionNameT  string
)

func (documentNameT) FriendlyName() string {
	return "document name"
}

func (pluginKindT) FriendlyName() string {
	return "plugin kind"
}

func (pluginNameT) FriendlyName() string {
	return "plugin name"
}

func (blockNameT) FriendlyName() string {
	return "block name"
}

func (sectionNameT) FriendlyName() string {
	return "section name"
}

type key interface {
	~string
	tree.Namer
}

type val interface {
	tree.Node
	tree.CtyAble
}

type nodeMap[K key, V val] struct {
	tree.NodeSigil
	Map map[string]V
}

func (m nodeMap[K, V]) CtyType() cty.Type {
	var v V
	return cty.Map(v.CtyType())
}

func (m nodeMap[K, V]) AsCtyValue() cty.Value {
	if len(m.Map) == 0 {
		var v V
		return cty.MapValEmpty(v.CtyType())
	}
	mp := make(map[string]cty.Value, len(m.Map))
	for k, v := range m.Map {
		mp[k] = v.AsCtyValue()
	}
	return cty.MapVal(mp)
}

func (m nodeMap[K, V]) FriendlyName() string {
	var k K
	var v V
	return fmt.Sprintf(
		"map of string (%s) to %s", k.FriendlyName(), v.FriendlyName(),
	)
}

func (m nodeMap[K, V]) IndexStr(idx string) tree.Node {
	if v, found := m.Map[idx]; found {
		return v
	}
	return nil
}

func (m nodeMap[K, V]) Get(idx string) V {
	if m.Map != nil {
		if v, found := m.Map[idx]; found {
			return v
		}
	}
	var v V
	return v
}

func SetIfMissing[V interface {
	FabricBlock
	val
}](m nodeMap[pluginKindT, nodeMap[pluginNameT, nodeMap[blockNameT, V]]], pluginKind, pluginName, blockName string, val V) *hcl.Diagnostic {
	m1, found := m.Map[pluginKind]
	if !found {
		m1 = NewMap[pluginNameT, nodeMap[blockNameT, V]]()
		m2 := NewMap[blockNameT, V]()
		m2.Map[blockName] = val
		m1.Map[pluginName] = m2
		m.Map[pluginKind] = m1
		return nil
	}
	m2, found := m1.Map[pluginName]
	if !found {
		m2 = NewMap[blockNameT, V]()
		m2.Map[blockName] = val
		m1.Map[pluginName] = m2
		return nil
	}
	return AddIfMissing(m2.Map, blockName, val)
}

type FabricBlock interface {
	HCLBlock() *hclsyntax.Block
	tree.Node
	tree.CtyAble
}

func MergeNestedMap[V FabricBlock](dst, src nodeMap[pluginKindT, nodeMap[pluginNameT, nodeMap[blockNameT, V]]]) (diags diagnostics.Diag) {
	for k1, src1 := range src.Map {
		dst1, found := dst.Map[k1]
		if !found {
			dst.Map[k1] = src1
			continue
		}
		for k2, src2 := range src1.Map {
			dst2, found := dst1.Map[k2]
			if !found {
				dst1.Map[k2] = src2
				continue
			}
			diags.Extend(MergeMap(dst2, src2))
		}
	}
	return nil
}

func MergeMap[K key, V FabricBlock](dst, src nodeMap[K, V]) (diags diagnostics.Diag) {
	for k, v := range src.Map {
		diags.Append(AddIfMissing(dst.Map, k, v))
	}
	return
}
