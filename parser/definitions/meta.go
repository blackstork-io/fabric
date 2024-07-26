package definitions

import "github.com/blackstork-io/fabric/plugin/plugindata"

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

func (m *MetaBlock) AsPluginData() plugindata.Data {
	tags := make(plugindata.List, len(m.Tags))
	authors := make(plugindata.List, len(m.Authors))
	for i, tag := range m.Tags {
		tags[i] = plugindata.String(tag)
	}
	for i, author := range m.Authors {
		authors[i] = plugindata.String(author)
	}
	return plugindata.Map{
		"authors": authors,
		"name":    plugindata.String(m.Name),
		"tags":    tags,
		"version": plugindata.String(m.Version),
	}
}
