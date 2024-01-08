package parser

import (
	"sync"

	"github.com/hashicorp/hcl/v2/hclsyntax"

	"github.com/blackstork-io/fabric/pkg/diagnostics"
)

type FabricBlock interface {
	Block() *hclsyntax.Block
}

type PluginBlock struct {
	block  *hclsyntax.Block
	Name   string
	Config *Config
}

var _ FabricBlock = (*PluginBlock)(nil)

func (b *PluginBlock) Block() *hclsyntax.Block {
	return b.block
}

type Key struct {
	PluginKind string
	PluginName string
	BlockName  string
}

type Config struct {
	block *hclsyntax.Block
	once  sync.Once
}

func (c *Config) GetKey() *Key {
	name := ""
	switch len(c.block.Labels) {
	case 0:
		// anonymous config block
		return nil
	case 3:
		// named config block
		name = c.block.Labels[2]
		fallthrough
	case 2:
		// default config block
		return &Key{
			PluginKind: c.block.Labels[0],
			PluginName: c.block.Labels[1],
			BlockName:  name,
		}
	default:
		panic("Invalid parsed config block")
	}
}

var _ FabricBlock = (*Config)(nil)

func (c *Config) Block() *hclsyntax.Block {
	return c.block
}

func parseConfigDefinition(block *hclsyntax.Block) (config *Config, diags diagnostics.Diag) {
	diags.Append(validatePluginKind(block, 0))
	diags.Append(validatePluginName(block, 1))
	diags.Append(validateBlockName(block, 2, true))
	diags.Append(validateLabelsLength(block, 3, "plugin_kind plugin_name block_name"))

	if diags.HasErrors() {
		return
	}
	config = &Config{
		block: block,
	}
	return
}

type Plugin struct {
	block *hclsyntax.Block
	once  sync.Once
}

func (p *Plugin) GetKey() *Key {
	if len(p.block.Labels) != 2 {
		return nil
	}
	return &Key{
		PluginKind: p.block.Type,
		PluginName: p.block.Labels[0],
		BlockName:  p.block.Labels[1],
	}
}

var _ FabricBlock = (*Plugin)(nil)

func (p *Plugin) Block() *hclsyntax.Block {
	return p.block
}

func parsePluginDefinition(block *hclsyntax.Block) (plugin *Plugin, diags diagnostics.Diag) {
	diags.Append(validatePluginName(block, 0))
	diags.Append(validateBlockName(block, 1, true))
	diags.Append(validateLabelsLength(block, 2, "plugin_name block_name"))

	if diags.HasErrors() {
		return
	}

	plugin = &Plugin{
		block: block,
	}
	return
}

type Document struct {
	block *hclsyntax.Block
	once  sync.Once
}

var _ FabricBlock = (*Document)(nil)

func (d *Document) Block() *hclsyntax.Block {
	return d.block
}

func (d *Document) Name() string {
	return d.block.Labels[0]
}

func parseDocumentDefinition(block *hclsyntax.Block) (doc *Document, diags diagnostics.Diag) {
	diags.Append(validateBlockName(block, 0, true))
	diags.Append(validateLabelsLength(block, 1, "block_name"))

	if diags.HasErrors() {
		return
	}
	doc = &Document{
		block: block,
	}
	return
}

type Section struct {
	block *hclsyntax.Block
	once  sync.Once
}

func (s *Section) Name() string {
	if len(s.block.Labels) == 0 {
		return ""
	}
	return s.block.Labels[0]
}

var _ FabricBlock = (*Section)(nil)

func (s *Section) Block() *hclsyntax.Block {
	return s.block
}

func parseSectionDefinition(block *hclsyntax.Block) (s *Section, diags diagnostics.Diag) {
	diags.Append(validateBlockName(block, 0, true))
	diags.Append(validateLabelsLength(block, 1, "block_name"))
	if diags.HasErrors() {
		return
	}
	s = &Section{
		block: block,
	}
	return
}

func parseBlockDefinitions(body *hclsyntax.Body) (res *DefinedBlocks, diags diagnostics.Diag) {
	res = NewDefinedBlocks()

	for _, block := range body.Blocks {
		switch block.Type {
		case BlockKindData, BlockKindContent:
			plugin, dgs := parsePluginDefinition(block)
			if diags.Extend(dgs) {
				continue
			}
			key := plugin.GetKey()
			if key == nil {
				panic("unable to get the key of the top-level block")
			}
			diags.Append(AddIfMissing(res.Plugins, *key, plugin))
		case BlockKindDocument:
			doc, dgs := parseDocumentDefinition(block)
			if diags.Extend(dgs) {
				continue
			}
			diags.Append(AddIfMissing(res.Documents, doc.Name(), doc))
		case BlockKindSection:
			section, dgs := parseSectionDefinition(block)
			if diags.Extend(dgs) {
				continue
			}
			key := section.Name()
			if key == "" {
				panic("unable to get the key of the top-level block")
			}
			diags.Append(AddIfMissing(res.Sections, key, section))
		case BlockKindConfig:
			cfg, dgs := parseConfigDefinition(block)
			if diags.Extend(dgs) {
				continue
			}
			key := cfg.GetKey()
			if key == nil {
				panic("unable to get the key of the top-level block")
			}
			diags.Append(AddIfMissing(res.Config, *key, cfg))
		default:
			diags.Append(newNestingDiag(block, body, []string{
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
