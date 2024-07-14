package pluginapiv1

import (
	"fmt"

	"github.com/blackstork-io/fabric/pkg/utils"
	"github.com/blackstork-io/fabric/plugin"
)

func decodeData(src *Data) plugin.Data {
	switch src.GetData().(type) {
	case nil:
		return nil
	case *Data_NumberVal:
		return plugin.NumberData(src.GetNumberVal())
	case *Data_StringVal:
		return plugin.StringData(src.GetStringVal())
	case *Data_BoolVal:
		return plugin.BoolData(src.GetBoolVal())
	case *Data_MapVal:
		return decodeMapData(src.GetMapVal().GetValue())
	case *Data_ListVal:
		return plugin.ListData(utils.FnMap(src.GetListVal().GetValue(), decodeData))
	}
	panic(fmt.Sprintf("Unexpected src data type: %T", src.GetData()))
}

func decodeMapData(src map[string]*Data) plugin.MapData {
	return plugin.MapData(utils.MapMap(src, decodeData))
}
