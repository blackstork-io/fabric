package builtin

import (
	"bytes"
	"fmt"
	"strings"
	"text/template"

	"github.com/Masterminds/sprig/v3"

	"github.com/blackstork-io/fabric/plugin"
	"github.com/blackstork-io/fabric/plugin/plugindata"
)

func countDeclarations(data *plugin.ContentSectionOrDoc, name string) int {
	count := 0
	for _, child := range data.Children {
		if section, ok := child.(*plugin.ContentSectionOrDoc); ok {
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

func parseScope(datactx plugindata.Map) (document, section *plugin.ContentSectionOrDoc) {
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

func findDepth(parent *plugin.ContentSectionOrDoc, id uint32, depth int) int {
	if parent.ID() == id {
		return depth
	}
	for _, child := range parent.Children {
		if child.ID() == id {
			return depth
		}
		if child, ok := child.(*plugin.ContentSectionOrDoc); ok {
			if d := findDepth(child, id, depth+1); d > -1 {
				return d
			}
		}
	}
	return -1
}

func firstTitle(el plugin.Content) (string, bool) {
	switch el := el.(type) {
	case *plugin.ContentSectionOrDoc:
		for _, c := range el.Children {
			if title, ok := firstTitle(c); ok {
				return title, true
			}
		}
	case *plugin.ContentElement:
		meta := el.Meta()
		if meta.Plugin == "blackstork/builtin" && meta.Provider == "title" {
			return string(bytes.TrimSpace(
				bytes.TrimPrefix(el.AsMarkdownSrc(), []byte("#")),
			)), true
		}
	}
	return "", false
}

func templateString(str string, datactx plugindata.Map) (string, error) {
	tmpl, err := template.New("pattern").Funcs(sprig.FuncMap()).Parse(str)
	if err != nil {
		return "", fmt.Errorf("failed to parse a text template: %w", err)
	}

	var buf bytes.Buffer
	err = tmpl.Execute(&buf, datactx.Any())
	if err != nil {
		return "", fmt.Errorf("failed to execute a text template: %w", err)
	}
	return strings.TrimSpace(buf.String()), nil
}
