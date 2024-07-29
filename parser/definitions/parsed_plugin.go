package definitions

import (
	"github.com/blackstork-io/fabric/parser/evaluation"
)

type ParsedPlugin struct {
	PluginName string
	BlockName  string
	Meta       *MetaBlock
	Config     evaluation.Configuration
	Invocation *evaluation.BlockInvocation
	Vars       *ParsedVars
}

type ParsedContent struct {
	Section *ParsedSection
	Plugin  *ParsedPlugin
}
