package builtin

import (
	"github.com/blackstork-io/fabric/plugin"
	"github.com/blackstork-io/fabric/plugin/plugindata"
)

func countDeclarations(data *plugin.ContentSection, name string) int {
	count := 0
	for _, child := range data.Children {
		if section, ok := child.(*plugin.ContentSection); ok {
			count += countDeclarations(section, name)
			continue
		}
		if element, ok := child.(*plugin.ContentElement); ok {
			meta := element.Meta()
			if meta != nil && meta.Provider == name {
				count++
			}
		}
	}
	return count
}

func parseScope(datactx plugindata.Map) (document, section *plugin.ContentSection) {
	documentMap, ok := datactx["document"]
	if !ok {
		return
	}

	contentMap, ok := documentMap.(plugindata.Map)["content"]
	if !ok {
		return
	}

	content, err := plugin.ParseContentData(contentMap.(plugindata.Map))
	if err != nil {
		return
	}

	document, ok = content.(*plugin.ContentSection)
	if !ok {
		return
	}

	sectionMap, ok := datactx["section"]
	if !ok || sectionMap == nil {
		return
	}
	contentMap, ok = sectionMap.(plugindata.Map)["content"]
	if !ok {
		return
	}
	content, err = plugin.ParseContentData(contentMap.(plugindata.Map))
	if err != nil {
		return
	}
	section, ok = content.(*plugin.ContentSection)
	if !ok {
		return
	}
	return document, section
}

func findDepth(parent *plugin.ContentSection, id uint32, depth int) int {
	if parent.ID() == id {
		return depth
	}
	for _, child := range parent.Children {
		if child.ID() == id {
			return depth
		}
		if child, ok := child.(*plugin.ContentSection); ok {
			if d := findDepth(child, id, depth+1); d > -1 {
				return d
			}
		}
	}
	return -1
}
