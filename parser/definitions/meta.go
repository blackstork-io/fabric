package definitions

import "github.com/blackstork-io/fabric/plugin"

type MetaBlock struct {
	// XXX: is empty sting enougth or use a proper ptr-nil-if-missing?
	Author string   `hcl:"author,optional"`
	Tags   []string `hcl:"tags,optional"`

	// TODO: ?store def range defRange hcl.Range
}

func (m *MetaBlock) AsJQ() plugin.Data {
	tags := make(plugin.ListData, len(m.Tags))
	for i, tag := range m.Tags {
		tags[i] = plugin.StringData(tag)
	}
	return plugin.ConvMapData{
		"author": plugin.StringData(m.Author),
		"tags":   tags,
	}
}
