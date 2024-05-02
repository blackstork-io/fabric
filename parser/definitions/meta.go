package definitions

import "github.com/blackstork-io/fabric/plugin"

type MetaBlock struct {
	Name        string   `hcl:"name,optional"`
	Description string   `hcl:"description,optional"`
	Url         string   `hcl:"url,optional"`
	License     string   `hcl:"license,optional"`
	Authors     []string `hcl:"authors,optional"`
	Tags        []string `hcl:"tags,optional"`
	UpdatedAt   string   `hcl:"updated_at,optional"`
	Version     string   `hcl:"version,optional"`

	// TODO: ?store def range defRange hcl.Range
}

func (m *MetaBlock) AsJQData() plugin.Data {
	tags := make(plugin.ListData, len(m.Tags))
	authors := make(plugin.ListData, len(m.Authors))
	for i, tag := range m.Tags {
		tags[i] = plugin.StringData(tag)
	}
	for i, author := range m.Authors {
		authors[i] = plugin.StringData(author)
	}
	return plugin.ConvMapData{
		"authors": authors,
		"name":    plugin.StringData(m.Name),
		"tags":    tags,
		"version": plugin.StringData(m.Version),
	}
}
