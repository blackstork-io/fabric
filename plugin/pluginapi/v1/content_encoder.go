package pluginapiv1

import "github.com/blackstork-io/fabric/plugin"

func encodeContent(src *plugin.Content) *Content {
	if src == nil {
		return nil
	}
	return &Content{
		Markdown: src.Markdown,
	}
}
