package evaluation

type Plugin struct {
	PluginName string
	BlockName  string
	Config     Configuration
	Invocation Invocation
}

func (pe *Plugin) AsBlockInvocation() *BlockInvocation {
	res, ok := pe.Invocation.(*BlockInvocation)
	if !ok {
		panic("This Plugin does not store a BlockInvocation!")
	}
	return res
}
