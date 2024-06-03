package pluginapiv1

import (
	"github.com/blackstork-io/fabric/pkg/utils"
	"github.com/blackstork-io/fabric/plugin"
)

func encodeData(d plugin.Data) *Data {
	switch v := d.(type) {
	case nil:
		return nil
	case plugin.NumberData:
		return &Data{
			Data: &Data_NumberVal{
				NumberVal: float64(v),
			},
		}
	case plugin.StringData:
		return &Data{
			Data: &Data_StringVal{
				StringVal: string(v),
			},
		}
	case plugin.BoolData:
		return &Data{
			Data: &Data_BoolVal{
				BoolVal: bool(v),
			},
		}
	case plugin.MapData:
		return &Data{
			Data: &Data_MapVal{
				encodeMapData(v),
			},
		}
	case plugin.ListData:
		return &Data{
			Data: &Data_ListVal{
				ListVal: &ListData{
					Value: utils.FnMap(v, encodeData),
				},
			},
		}
	}
	panic("unreachable")
}

func encodeMapData(m plugin.MapData) *MapData {
	return &MapData{
		Value: utils.MapMap(m, encodeData),
	}
}
