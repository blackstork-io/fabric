package plugin

type LocationEffect int

const (
	LocationEffectUnspecified LocationEffect = iota
	LocationEffectBefore
	LocationEffectAfter
)

type Location struct {
	Index  int
	Effect LocationEffect
}

func (loc *Location) AsData() Data {
	if loc == nil {
		return nil
	}
	return MapData{
		"index":  NumberData(loc.Index),
		"effect": NumberData(loc.Effect),
	}
}

type Content struct {
	Markdown string
	Location *Location
	meta     *ContentMeta
}

func (c *Content) Meta() *ContentMeta {
	if c == nil {
		return nil
	}
	return c.meta
}

func (c *Content) AsData() Data {
	if c == nil {
		return nil
	}
	return MapData{
		"markdown": StringData(c.Markdown),
		"location": c.Location.AsData(),
		"meta":     c.meta.AsData(),
	}
}

type ContentMeta struct {
	Provider string
	Plugin   string
	Version  string
}

func (meta *ContentMeta) AsData() Data {
	if meta == nil {
		return nil
	}
	return MapData{
		"provider": StringData(meta.Provider),
	}
}

func ParseContentData(data MapData) *Content {
	if data == nil {
		return nil
	}
	md, _ := data["markdown"].(StringData)
	return &Content{
		Markdown: string(md),
		Location: ParseLocationData(data["location"]),
		meta:     ParseContentMeta(data["meta"]),
	}
}

func ParseLocationData(data Data) *Location {
	if data == nil {
		return nil
	}
	loc := data.(MapData)
	index, _ := loc["index"].(NumberData)
	effect, _ := loc["effect"].(NumberData)
	return &Location{
		Index:  int(index),
		Effect: LocationEffect(effect),
	}
}

func ParseContentMeta(data Data) *ContentMeta {
	if data == nil {
		return nil
	}
	meta := data.(MapData)
	provider, _ := meta["provider"].(StringData)
	return &ContentMeta{
		Provider: string(provider),
	}
}
