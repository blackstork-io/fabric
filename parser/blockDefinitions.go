package parser

import (
	"sync"

	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/hcldec"
	"github.com/hashicorp/hcl/v2/hclsyntax"
	"github.com/zclconf/go-cty/cty"

	"github.com/blackstork-io/fabric/pkg/diagnostics"
)

type FabricBlock interface {
	GetHCLBlock() *hcl.Block
}

type PluginBlock struct {
	block  *hclsyntax.Block
	Name   string
	Config *Config
}

var _ FabricBlock = (*PluginBlock)(nil)

func (b *PluginBlock) GetHCLBlock() *hcl.Block {
	return b.block.AsHCLBlock()
}

type Key struct {
	PluginKind string
	PluginName string
	BlockName  string
}

type Config struct {
	*hcl.Block
	BlockRange hcl.Range
	once       sync.Once
	val        cty.Value
}

// Parse implements ConfigurationObject.
func (c *Config) Parse(spec hcldec.Spec) (val cty.Value, diags diagnostics.Diag) {
	c.once.Do(func() {
		var diag hcl.Diagnostics
		c.val, diag = hcldec.Decode(c.Body, spec, nil)
		if diags.ExtendHcl(diag) {
			// don't let partially-decoded values live
			c.val = cty.NilVal
		}
	})
	val = c.val
	if val.IsNull() && diags == nil {
		diags.Append(diagnostics.RepeatedError)
	}
	return
}

// Range implements ConfigurationObject.
func (c *Config) Range() hcl.Range {
	return c.BlockRange
}

func (c *Config) GetKey() *Key {
	name := ""
	switch len(c.Block.Labels) {
	case 0:
		// anonymous config block
		return nil
	case 3:
		// named config block
		name = c.Block.Labels[2]
		fallthrough
	case 2:
		// default config block
		return &Key{
			PluginKind: c.Block.Labels[0],
			PluginName: c.Block.Labels[1],
			BlockName:  name,
		}
	default:
		panic("Invalid parsed config block")
	}
}

var _ FabricBlock = (*Config)(nil)

func (c *Config) GetHCLBlock() *hcl.Block {
	return c.Block
}

func DefineConfig(block *hclsyntax.Block) (config *Config, diags diagnostics.Diag) {
	diags.Append(validatePluginKindLabel(block, 0))
	diags.Append(validatePluginName(block, 1))
	diags.Append(validateBlockName(block, 2, true))
	diags.Append(validateLabelsLength(block, 3, "plugin_kind plugin_name block_name"))

	if diags.HasErrors() {
		return
	}
	config = &Config{
		Block:      block.AsHCLBlock(),
		BlockRange: block.Range(),
	}
	return
}

type Plugin struct {
	block *hclsyntax.Block

	pluginName string
	once       sync.Once
	isValid    bool
	invoke     *blockInvocation
	config     Configuration
}

func (p *Plugin) DefRange() hcl.Range {
	return p.block.DefRange()
}

func (p *Plugin) Kind() string {
	return p.block.Type
}

// Current plugin name. For unevaluated refs is "ref",
// after evaluation will change to the referenced plugin name.
func (p *Plugin) PluginName() string {
	return p.pluginName
}

// Whether or not the original block was a ref
func (p *Plugin) IsRef() bool {
	return p.block.Labels[0] == PluginTypeRef
}

func (p *Plugin) BlockName() string {
	if len(p.block.Labels) < 2 {
		return ""
	}
	return p.block.Labels[1]
}

func (p *Plugin) GetKey() *Key {
	blockName := p.BlockName()
	if blockName == "" {
		return nil
	}
	return &Key{
		PluginKind: p.Kind(),
		PluginName: p.PluginName(),
		BlockName:  blockName,
	}
}

var _ FabricBlock = (*Plugin)(nil)

func (p *Plugin) GetHCLBlock() *hcl.Block {
	return p.block.AsHCLBlock()
}

func DefinePlugin(block *hclsyntax.Block, atTopLevel bool) (plugin *Plugin, diags diagnostics.Diag) {
	nameRequired := atTopLevel || (block.Type == BlockKindData)

	diags.Append(validatePluginKind(block, block.Type, block.TypeRange))
	diags.Append(validatePluginName(block, 0))
	if nameRequired {
		diags.Append(validateBlockName(block, 1, true))
		diags.Append(validateLabelsLength(block, 2, "plugin_name block_name"))
	} else {
		diags.Append(validateBlockName(block, 1, false))
		diags.Append(validateLabelsLength(block, 2, "plugin_name <block_name>"))
	}

	if diags.HasErrors() {
		return
	}

	plugin = &Plugin{
		block:      block,
		pluginName: block.Labels[0], // always required, so no bounds checking
	}

	return
}

// Document and section are very similar conceptually
type DocumentOrSection struct {
	block *hclsyntax.Block
	once  sync.Once
	meta  MetaBlock
}

func (d *DocumentOrSection) IsDocument() bool {
	return d.block.Type == BlockKindDocument
}

var _ FabricBlock = (*DocumentOrSection)(nil)

func (d *DocumentOrSection) GetHCLBlock() *hcl.Block {
	return d.block.AsHCLBlock()
}

func (d *DocumentOrSection) Name() string {
	return d.block.Labels[0]
}

func DefineSectionOrDocument(block *hclsyntax.Block, atTopLevel bool) (doc *DocumentOrSection, diags diagnostics.Diag) {
	nameRequired := atTopLevel || block.Type == BlockKindDocument

	if nameRequired {
		diags.Append(validateBlockName(block, 0, true))
		diags.Append(validateLabelsLength(block, 1, "block_name"))
	} else {
		diags.Append(validateBlockName(block, 0, false))
		diags.Append(validateLabelsLength(block, 1, "<block_name>"))
	}

	if diags.HasErrors() {
		return
	}
	doc = &DocumentOrSection{
		block: block,
	}
	return
}

func parseBlockDefinitions(body *hclsyntax.Body) (res *DefinedBlocks, diags diagnostics.Diag) {
	res = NewDefinedBlocks()

	for _, block := range body.Blocks {
		switch block.Type {
		case BlockKindData, BlockKindContent:
			plugin, dgs := DefinePlugin(block, true)
			if diags.Extend(dgs) {
				continue
			}
			key := plugin.GetKey()
			if key == nil {
				panic("unable to get the key of the top-level block")
			}
			diags.Append(AddIfMissing(res.Plugins, *key, plugin))
		case BlockKindDocument, BlockKindSection:
			blk, dgs := DefineSectionOrDocument(block, true)
			if diags.Extend(dgs) {
				continue
			}
			if blk.IsDocument() {
				diags.Append(AddIfMissing(res.Documents, blk.Name(), blk))
			} else {
				diags.Append(AddIfMissing(res.Sections, blk.Name(), blk))
			}
		case BlockKindConfig:
			cfg, dgs := DefineConfig(block)
			if diags.Extend(dgs) {
				continue
			}
			key := cfg.GetKey()
			if key == nil {
				panic("unable to get the key of the top-level block")
			}
			diags.Append(AddIfMissing(res.Config, *key, cfg))
		default:
			diags.Append(newNestingDiag(
				"Top level of fabric document",
				block,
				body,
				[]string{
					BlockKindData,
					BlockKindContent,
					BlockKindDocument,
					BlockKindSection,
					BlockKindConfig,
				}))
		}
	}
	return
}
