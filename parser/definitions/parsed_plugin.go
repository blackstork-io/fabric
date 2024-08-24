package definitions

import (
	"github.com/blackstork-io/fabric/parser/evaluation"
)

type ParsedPlugin struct {
	Source       *Plugin
	PluginName   string
	BlockName    string
	Meta         *MetaBlock
	Config       evaluation.Configuration
	Invocation   *evaluation.BlockInvocation
	Vars         *ParsedVars
	RequiredVars []string
}

type ParsedContent struct {
	Section *ParsedSection
	Plugin  *ParsedPlugin
}
