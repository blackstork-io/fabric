package pluginapiv1

import "github.com/blackstork-io/fabric/plugin"

func encodeContent(src *plugin.Content) *Content {
	if src == nil {
		return nil
	}
	return &Content{
		Markdown: src.Markdown,
		Location: encodeLocation(src.Location),
	}
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
