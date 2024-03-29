package pluginapiv1

import "github.com/blackstork-io/fabric/plugin"

func decodeContentResult(src *ContentResult) *plugin.ContentResult {
	if src == nil {
		return nil
	}
	return &plugin.ContentResult{
		Content:  decodeContent(src.Content),
		Location: decodeLocation(src.Location),
	}
}

func decodeContent(src *Content) plugin.Content {
	if src == nil {
		return nil
	}
	switch val := src.Value.(type) {
	case *Content_Element:
		return decodeContentElement(val.Element)
	case *Content_Section:
		return decodeContentSection(val.Section)
	case *Content_Empty:
		return decodeContentEmpty(val.Empty)
	default:
		return nil
	}
}

func decodeContentElement(src *ContentElement) *plugin.ContentElement {
	if src == nil {
		return nil
	}
	return &plugin.ContentElement{
		Markdown: src.GetMarkdown(),
	}
}

func decodeContentEmpty(src *ContentEmpty) *plugin.ContentEmpty {
	if src == nil {
		return nil
	}
	return &plugin.ContentEmpty{}
}

func decodeContentSection(src *ContentSection) *plugin.ContentSection {
	if src == nil {
		return nil
	}
	children := make([]plugin.Content, len(src.GetChildren()))
	for i, child := range src.GetChildren() {
		children[i] = decodeContent(child)
	}
	return &plugin.ContentSection{
		Children: children,
	}
}

func decodeLocation(src *Location) *plugin.Location {
	if src == nil {
		return nil
	}
	return &plugin.Location{
		Index:  src.GetIndex(),
		Effect: decodeLocationEffect(src.GetEffect()),
	}
}

func decodeLocationEffect(src LocationEffect) plugin.LocationEffect {
	switch src {
	case LocationEffect_LOCATION_EFFECT_BEFORE:
		return plugin.LocationEffectBefore
	case LocationEffect_LOCATION_EFFECT_AFTER:
		return plugin.LocationEffectAfter
	default:
		return plugin.LocationEffectUnspecified
	}
}
