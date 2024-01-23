package definitions

import (
	"sync"

	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/hclsyntax"

	"github.com/blackstork-io/fabric/parser/evaluation"
	"github.com/blackstork-io/fabric/pkg/diagnostics"
)

type Plugin struct {
	Block *hclsyntax.Block

	// Current plugin name. For unevaluated refs is "ref",
	// after evaluation will change to the referenced plugin name.
	Once       sync.Once
	Evaluated  bool
	EvalResult *evaluation.Plugin
}

func (p *Plugin) DefRange() hcl.Range {
	return p.Block.DefRange()
}

func (p *Plugin) Kind() string {
	return p.Block.Type
}

func (p *Plugin) Name() string {
	if p.Evaluated {
		// resolved plugin name (in case of ref)
		return p.EvalResult.PluginName
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

var _ FabricBlock = (*Plugin)(nil)

func (p *Plugin) GetHCLBlock() *hcl.Block {
	return p.Block.AsHCLBlock()
}

func DefinePlugin(block *hclsyntax.Block, atTopLevel bool) (plugin *Plugin, diags diagnostics.Diag) {
	nameRequired := atTopLevel || (block.Type == BlockKindData)

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

	if diags.HasErrors() {
		return
	}

	plugin = &Plugin{
		Block: block,
	}

	return
}
