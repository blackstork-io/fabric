package pluginapiv1

import "github.com/blackstork-io/fabric/plugin"

func decodeContent(src *Content) *plugin.Content {
	if src == nil {
		return nil
	}
	return &plugin.Content{
		Markdown: src.GetMarkdown(),
	}
}
