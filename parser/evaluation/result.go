package evaluation

import (
	"github.com/blackstork-io/fabric/plugin"
)

type Result plugin.ListData

func (d *Result) AsJQData() plugin.Data {
	return plugin.ListData(*d)
}

func (d Result) AsGoType() (result []string) {
	result = make([]string, len(d))
	for i, s := range d {
		// The only way to modify resultsList is through append, so this is always correct
		content := plugin.ParseContentData(s.(plugin.MapData))
		if content == nil {
			continue
		}
		result[i] = content.Markdown
	}
	return
}

func (d *Result) Append(s *plugin.Content) {
	if s == nil {
		return
	}
	if s.Location != nil {
		idx := s.Location.Index
		if idx < 0 {
			idx = 0
		}
		if idx >= len(*d) {
			*d = append(*d, s.AsData())
			return
		}
		switch s.Location.Effect {
		case plugin.LocationEffectBefore:
			*d = append((*d)[:idx], append(plugin.ListData{s.AsData()}, (*d)[idx:]...)...)
			return
		case plugin.LocationEffectAfter:
			*d = append((*d)[:idx+1], append(plugin.ListData{s.AsData()}, (*d)[idx+1:]...)...)
			return
		default:
			(*d)[idx] = s.AsData()
			return
		}
	}
	*d = append(*d, s.AsData())
}
