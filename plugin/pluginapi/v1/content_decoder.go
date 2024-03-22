package pluginapiv1

import "github.com/blackstork-io/fabric/plugin"

func decodeContent(src *Content) *plugin.Content {
	if src == nil {
		return nil
	}
	return &plugin.Content{
		Markdown: src.GetMarkdown(),
		Location: decodeLocation(src.GetLocation()),
	}
}

func decodeLocation(src *Location) *plugin.Location {
	if src == nil {
		return nil
	}
	return &plugin.Location{
		Index:  int(src.GetIndex()),
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
