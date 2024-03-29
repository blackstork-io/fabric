package pluginapiv1

import "github.com/blackstork-io/fabric/plugin"

func encodeContentResult(src *plugin.ContentResult) *ContentResult {
	if src == nil {
		return nil
	}
	return &ContentResult{
		Content:  encodeContent(src.Content),
		Location: encodeLocation(src.Location),
	}
}

func encodeContent(src plugin.Content) *Content {
	if src == nil {
		return nil
	}
	switch val := src.(type) {
	case *plugin.ContentElement:
		return &Content{
			Value: &Content_Element{
				Element: encodeContentElement(val),
			},
		}
	case *plugin.ContentSection:
		return &Content{
			Value: &Content_Section{
				Section: encodeContentSection(val),
			},
		}
	case *plugin.ContentEmpty:
		return &Content{
			Value: &Content_Empty{
				Empty: encodeContentEmpty(val),
			},
		}
	default:
		return nil
	}
}

func encodeContentSection(src *plugin.ContentSection) *ContentSection {
	if src == nil {
		return nil
	}
	children := make([]*Content, len(src.Children))
	for i, child := range src.Children {
		children[i] = encodeContent(child)
	}
	return &ContentSection{
		Children: children,
	}
}

func encodeContentElement(src *plugin.ContentElement) *ContentElement {
	if src == nil {
		return nil
	}
	return &ContentElement{
		Markdown: src.Markdown,
	}
}

func encodeContentEmpty(src *plugin.ContentEmpty) *ContentEmpty {
	if src == nil {
		return nil
	}
	return &ContentEmpty{}
}

func encodeLocation(src *plugin.Location) *Location {
	if src == nil {
		return nil
	}
	return &Location{
		Index:  uint32(src.Index),
		Effect: encodeLocationEffect(src.Effect),
	}
}

func encodeLocationEffect(src plugin.LocationEffect) LocationEffect {
	switch src {
	case plugin.LocationEffectBefore:
		return LocationEffect_LOCATION_EFFECT_BEFORE
	case plugin.LocationEffectAfter:
		return LocationEffect_LOCATION_EFFECT_AFTER
	default:
		return LocationEffect_LOCATION_EFFECT_UNSPECIFIED
	}
}
