package plugin

import (
	"encoding/json"
	"fmt"
)

type Data interface {
	Any() any
	data()
	ConvertableData
}

func (NumberData) data() {}
func (StringData) data() {}
func (BoolData) data()   {}
func (MapData) data()    {}
func (ListData) data()   {}

func (d NumberData) AsJQ() Data { return d }
func (d StringData) AsJQ() Data { return d }
func (d BoolData) AsJQ() Data   { return d }
func (d MapData) AsJQ() Data    { return d }
func (d ListData) AsJQ() Data   { return d }

type NumberData float64

func (d NumberData) Any() any {
	return float64(d)
}

type StringData string

func (d StringData) Any() any {
	return string(d)
}

type BoolData bool

func (d BoolData) Any() any {
	return bool(d)
}

type MapData map[string]Data

func (d MapData) Any() any {
	dst := make(map[string]any)
	for k, v := range d {
		if v == nil {
			dst[k] = nil
			continue
		}
		dst[k] = v.Any()
	}
	return dst
}

type ListData []Data

func (d ListData) Any() any {
	dst := make([]any, len(d))
	for i, v := range d {
		if v == nil {
			dst[i] = nil
			continue
		}
		dst[i] = v.Any()
	}
	return dst
}

func UnmarshalJSONData(data []byte) (Data, error) {
	var v any
	if err := json.Unmarshal(data, &v); err != nil {
		return nil, err
	}
	return ParseDataAny(v)
}

func ParseDataAny(v any) (Data, error) {
	switch v := v.(type) {
	case nil:
		return nil, nil
	case bool:
		return BoolData(v), nil
	case float64:
		return NumberData(v), nil
	case float32:
		return NumberData(v), nil
	case uint:
		return NumberData(v), nil
	case uint8:
		return NumberData(v), nil
	case uint16:
		return NumberData(v), nil
	case uint32:
		return NumberData(v), nil
	case uint64:
		return NumberData(v), nil
	case int:
		return NumberData(v), nil
	case int8:
		return NumberData(v), nil
	case int16:
		return NumberData(v), nil
	case int32:
		return NumberData(v), nil
	case int64:
		return NumberData(v), nil
	case uintptr:
		return NumberData(v), nil
	case string:
		return StringData(v), nil
	case []any:
		// TODO: potential bug:
		// this case would trigger only for []any, not, for example, []string
		// this can be worked around using reflection
		dst := make(ListData, len(v))
		for i, e := range v {
			d, err := ParseDataAny(e)
			if err != nil {
				return nil, err
			}
			dst[i] = d
		}
		return dst, nil
	case map[string]any:
		return ParseDataMapAny(v)
	default:
		return nil, fmt.Errorf("unsupported data type %T", v)
	}
}

func ParseDataMapAny(v map[string]any) (MapData, error) {
	if v == nil {
		return nil, nil
	}
	dst := make(MapData)
	for k, e := range v {
		d, err := ParseDataAny(e)
		if err != nil {
			return nil, err
		}
		dst[k] = d
	}
	return dst, nil
}

type ConvertableData interface {
	AsJQ() Data
}

type ConvMapData map[string]ConvertableData

func (d ConvMapData) AsJQ() Data {
	dst := make(MapData, len(d))
	for k, v := range d {
		if v == nil {
			dst[k] = nil
		} else {
			dst[k] = v.AsJQ()
		}
	}
	return dst
}

func (d ConvMapData) Any() any {
	dst := make(map[string]any, len(d))
	for k, v := range d {
		dst[k] = v.AsJQ().Any()
	}
	return dst
}
func (d ConvMapData) data() {}

type ConvListData []ConvertableData

func (d ConvListData) AsJQ() Data {
	dst := make(ListData, len(d))
	for k, v := range d {
		dst[k] = v.AsJQ()
	}
	return dst
}

func (d ConvListData) Any() any {
	dst := make([]any, len(d))
	for k, v := range d {
		dst[k] = v.AsJQ().Any()
	}
	return dst
}
func (d ConvListData) data() {}
