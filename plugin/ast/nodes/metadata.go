package nodes

import "github.com/blackstork-io/fabric/plugin/plugindata"

type ContentMeta struct {
	Provider string
	Plugin   string
	Version  string
}

func (meta *ContentMeta) AsData() plugindata.Data {
	if meta == nil {
		return nil
	}
	return plugindata.Map{
		"provider": plugindata.String(meta.Provider),
		"plugin":   plugindata.String(meta.Plugin),
		"version":  plugindata.String(meta.Version),
	}
}
