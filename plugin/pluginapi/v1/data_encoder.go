package pluginapiv1

import (
	"fmt"

	"github.com/blackstork-io/fabric/pkg/utils"
	"github.com/blackstork-io/fabric/plugin/plugindata"
)

func encodeData(d plugindata.Data) *Data {
	switch v := d.(type) {
	case nil:
		return nil
	case plugindata.Number:
		return &Data{
			Data: &Data_NumberVal{
				NumberVal: float64(v),
			},
		}
	case plugindata.String:
		return &Data{
			Data: &Data_StringVal{
				StringVal: string(v),
			},
		}
	case plugindata.Bool:
		return &Data{
			Data: &Data_BoolVal{
				BoolVal: bool(v),
			},
		}
	case plugindata.Map:
		return &Data{
			Data: &Data_MapVal{
				encodeMapData(v),
			},
		}
	case plugindata.List:
		return &Data{
			Data: &Data_ListVal{
				ListVal: &ListData{
					Value: utils.FnMap(v, encodeData),
				},
			},
		}
	default:
		if cd, ok := d.(plugindata.Convertible); ok {
			return encodeData(cd.AsPluginData())
		}
	}
	panic(fmt.Errorf("unexpected plugin data type: %T", d))
}

func encodeMapData(m plugindata.Map) *MapData {
	return &MapData{
		Value: utils.MapMap(m, encodeData),
	}
}
