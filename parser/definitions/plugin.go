package definitions

import (
	"sync"

	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/hclsyntax"
	"github.com/zclconf/go-cty/cty"

	"github.com/blackstork-io/fabric/parser/blocks/internal/tree"
	"github.com/blackstork-io/fabric/pkg/diagnostics"
)

type OnceParser struct {
	Block *hclsyntax.Block

	Once        sync.Once
	Parsed      bool
	ParseResult *ParsedPlugin
}

func (db *DefinedBlocks) ParsePlugin(plugin *definitions.Plugin) (res *definitions.ParsedPlugin, diags diagnostics.Diag) {
	if circularRefDetector.Check(plugin) {
		// This produces a bit of an incorrect error and shouldn't trigger in normal operation
		// but I re-check for the circular refs here out of abundance of caution:
		// deadlocks or infinite loops may occur, and are hard to debug
		diags.Append(&hcl.Diagnostic{
			Severity: hcl.DiagError,
			Summary:  "Circular reference detected",
			Detail:   "Looped back to this block through reference chain:",
			Subject:  plugin.DefRange().Ptr(),
			Extra:    circularRefDetector.ExtraMarker,
		})
		return
	}
	plugin.Once.Do(func() {
		res, diags = db.parsePlugin(plugin)
		if diags.HasErrors() {
			return
		}
		plugin.ParseResult = res
		plugin.Parsed = true
	})
	if !plugin.Parsed {
		if diags == nil {
			diags.Append(diagnostics.RepeatedError)
		}
		return
	}
	res = plugin.ParseResult
	return
}

type Plugin struct {
	Block *hclsyntax.Block

	Once        sync.Once
	Parsed      bool
	ParseResult *ParsedPlugin
}

func (p *Plugin) DefRange() hcl.Range {
	return p.Block.DefRange()
}

func (p *Plugin) Kind() string {
	return p.Block.Type
}

func (p *Plugin) Name() string {
	if p.Parsed {
		// resolved plugin name (in case of ref)
		return p.ParseResult.PluginName
	}
	return p.Block.Labels[0]
}

// Whether or not the original block is a ref.
func (p *Plugin) IsRef() bool {
	return p.Block.Labels[0] == PluginTypeRef
}

func (p *Plugin) BlockName() string {
	if len(p.Block.Labels) < 2 {
		return ""
	}
	return p.Block.Labels[1]
}

func (p *Plugin) GetKey() *Key {
	blockName := p.BlockName()
	if blockName == "" {
		return nil
	}
	return &Key{
		PluginKind: p.Kind(),
		PluginName: p.Block.Labels[0],
		BlockName:  blockName,
	}
}

func (p *Plugin) GetHCLBlock() *hcl.Block {
	return p.Block.AsHCLBlock()
}

func (p *Plugin) GetGenericPlugin() *Plugin {
	return p
}

func checkPlugin(block *hclsyntax.Block, atTopLevel bool) (diags diagnostics.Diag) {
	nameRequired := atTopLevel

	diags.Append(validatePluginKind(block, block.Type, block.TypeRange))
	diags.Append(validatePluginName(block, 0))
	diags.Append(validateBlockName(block, 1, nameRequired))
	var usage string
	if nameRequired {
		usage = "plugin_name block_name"
	} else {
		usage = "plugin_name <block_name>"
	}
	diags.Append(validateLabelsLength(block, 2, usage))
	return
}

func DefineContentPlugin(block *hclsyntax.Block, atTopLevel bool) (plugin *ContentPlugin, diags diagnostics.Diag) {
	diags = checkPlugin(block, atTopLevel)
	if diags.HasErrors() {
		return
	}
	plugin = &ContentPlugin{
		Plugin: Plugin{
			Block: block,
		},
	}
	return
}

func DefineDataPlugin(block *hclsyntax.Block, atTopLevel bool) (plugin *DataPlugin, diags diagnostics.Diag) {
	diags = checkPlugin(block, atTopLevel)
	if diags.HasErrors() {
		return
	}
	plugin = &DataPlugin{
		Plugin: Plugin{
			Block: block,
		},
	}
	return
}

type ContentPlugin struct {
	tree.NodeSigil
	Plugin
}

func (*ContentPlugin) FriendlyName() string {
	return "content plugin"
}

var ctyContentPluginType = capsuleTypeFor[ContentPlugin]()

func (p *ContentPlugin) CtyType() cty.Type {
	return ctyContentPluginType
}

func (p *ContentPlugin) AsCtyValue() cty.Value {
	return cty.CapsuleVal(p.CtyType(), p)
}

type DataPlugin struct {
	tree.NodeSigil
	Plugin
}

func (*DataPlugin) FriendlyName() string {
	return "data plugin"
}

var ctyDataPluginType = capsuleTypeFor[DataPlugin]()

func (p *DataPlugin) CtyType() cty.Type {
	return ctyDataPluginType
}

func (p *DataPlugin) AsCtyValue() cty.Value {
	return cty.CapsuleVal(p.CtyType(), p)
}

var (
	_ PluginIface = &ContentPlugin{}
	_ PluginIface = &DataPlugin{}
)

type PluginIface interface {
	tree.CtyAble
	tree.Node
	FabricBlock
	BlockName() string
	DefRange() hcl.Range
	GetHCLBlock() *hcl.Block
	GetKey() *Key
	IsRef() bool
	Kind() string
	Name() string
	GetGenericPlugin() *Plugin
}
