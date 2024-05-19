package definitions

import (
	"github.com/blackstork-io/fabric/parser/evaluation"
)

type ParsedPlugin struct {
	PluginName string
	BlockName  string
	Meta       *MetaBlock
	Config     evaluation.Configuration
	Invocation evaluation.Invocation
}

func (pe *ParsedPlugin) GetBlockInvocation() *evaluation.BlockInvocation {
	res, ok := pe.Invocation.(*evaluation.BlockInvocation)
	if !ok {
		panic("This Plugin does not store a BlockInvocation!")
	}
	return res
}

type ParsedContent struct {
	Section *ParsedSection
	Plugin  *ParsedPlugin
}
