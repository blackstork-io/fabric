package pluginapiv1

import "github.com/blackstork-io/fabric/plugin"

func decodeData(src *Data) plugin.Data {
	switch {
	case src == nil || src.GetData() == nil:
		return nil
	case src.GetNumberVal() != nil:
		return plugin.NumberData(src.GetNumberVal().GetValue())
	case src.GetStringVal() != nil:
		return plugin.StringData(src.GetStringVal().GetValue())
	case src.GetBoolVal() != nil:
		return plugin.BoolData(src.GetBoolVal().GetValue())
	case src.GetMapVal() != nil:
		return decodeMapData(src.GetMapVal())
	case src.GetListVal() != nil:
		dst := make(plugin.ListData, len(src.GetListVal().GetValue()))
		for i, v := range src.GetListVal().GetValue() {
			dst[i] = decodeData(v)
		}
		return dst
	}
	panic("unreachable")
}

func decodeMapData(src *MapData) plugin.MapData {
	dst := make(plugin.MapData)
	for k, v := range src.GetValue() {
		dst[k] = decodeData(v)
	}
	return dst
}
