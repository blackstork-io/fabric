package definitions

import "github.com/blackstork-io/fabric/plugin"

type MetaBlock struct {
	Name        string   `hcl:"name,optional"`
	Description string   `hcl:"description,optional"`
	Url         string   `hcl:"url,optional"`
	License     string   `hcl:"license,optional"`
	Author      string   `hcl:"author,optional"`
	Tags        []string `hcl:"tags,optional"`
	UpdatedAt   string   `hcl:"updated_at,optional"`

	// TODO: ?store def range defRange hcl.Range
}

func (m *MetaBlock) AsJQData() plugin.Data {
	tags := make(plugin.ListData, len(m.Tags))
	for i, tag := range m.Tags {
		tags[i] = plugin.StringData(tag)
	}
	return plugin.ConvMapData{
		"author": plugin.StringData(m.Author),
		"tags":   tags,
	}
}
