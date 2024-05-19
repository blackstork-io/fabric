package eval

import (
	"github.com/zclconf/go-cty/cty"

	"github.com/blackstork-io/fabric/parser/definitions"
)

type PluginAction struct {
	PluginName string
	BlockName  string
	Meta       *definitions.MetaBlock
	Config     cty.Value
	Args       cty.Value
}
