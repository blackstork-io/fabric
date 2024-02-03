package pluginapiv1

import "github.com/blackstork-io/fabric/plugin"

func encodeData(d plugin.Data) *Data {
	switch v := d.(type) {
	case nil:
		return nil
	case plugin.NumberData:
		return &Data{
			Data: &Data_NumberVal{
				NumberVal: &NumberData{
					Value: float64(v),
				},
			},
		}
	case plugin.StringData:
		return &Data{
			Data: &Data_StringVal{
				StringVal: &StringData{
					Value: string(v),
				},
			},
		}
	case plugin.BoolData:
		return &Data{
			Data: &Data_BoolVal{
				BoolVal: &BoolData{
					Value: bool(v),
				},
			},
		}
	case plugin.MapData:
		return &Data{
			Data: &Data_MapVal{
				encodeMapData(v),
			},
		}
	case plugin.ListData:
		l := make([]*Data, len(v))
		for i, d := range v {
			l[i] = encodeData(d)
		}
		return &Data{
			Data: &Data_ListVal{
				ListVal: &ListData{
					Value: l,
				},
			},
		}
	}
	panic("unreachable")
}

func encodeMapData(m plugin.MapData) *MapData {
	dst := make(map[string]*Data)
	for k, v := range m {
		dst[k] = encodeData(v)
	}
	return &MapData{
		Value: dst,
	}
}
