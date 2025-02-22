package plugindata

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"time"

	"gopkg.in/yaml.v3"
)

type Data interface {
	Any() any
	data()
	Convertible
}

func (Number) data() {}
func (String) data() {}
func (Bool) data()   {}
func (Map) data()    {}
func (List) data()   {}
func (Time) data()   {}

func (d Number) AsPluginData() Data { return d }
func (d String) AsPluginData() Data { return d }
func (d Bool) AsPluginData() Data   { return d }
func (d Map) AsPluginData() Data    { return d }
func (d List) AsPluginData() Data   { return d }
func (d Time) AsPluginData() Data   { return d }

type Number float64

func (d Number) Any() any {
	return float64(d)
}

type String string

func (d String) Any() any {
	return string(d)
}

type Bool bool

func (d Bool) Any() any {
	return bool(d)
}

type Map map[string]Data

func (d Map) Any() any {
	dst := make(map[string]any, len(d))
	for k, v := range d {
		if v == nil {
			dst[k] = nil
			continue
		}
		dst[k] = v.Any()
	}
	return dst
}

func (d Map) Clone() Map {
	dst := make(Map, len(d))
	for k, v := range d {
		switch v := v.(type) {
		case Map:
			dst[k] = v.Clone()
		case List:
			dst[k] = v.Clone()
		default:
			dst[k] = v
		}
	}
	return dst
}

type List []Data

func (d List) Any() any {
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

type Time time.Time

func (d Time) Any() any {
	return (time.Time)(d)
}

func (d List) Clone() List {
	dst := make(List, len(d))
	for i, v := range d {
		switch v := v.(type) {
		case Map:
			dst[i] = v.Clone()
		case List:
			dst[i] = v.Clone()
		default:
			dst[i] = v
		}
	}
	return dst
}

func UnmarshalJSON(data []byte) (Data, error) {
	var v any
	if err := json.Unmarshal(data, &v); err != nil {
		return nil, err
	}
	return ParseAny(v)
}

func UnmarshalYAML(data []byte) (Data, error) {
	var v any
	if err := yaml.Unmarshal(data, &v); err != nil {
		return nil, err
	}
	return ParseAny(v)
}

func ParseAny(v any) (Data, error) {
	switch v := v.(type) {
	case nil:
		return nil, nil
	case bool:
		return Bool(v), nil
	case float64:
		return Number(v), nil
	case float32:
		return Number(v), nil
	case uint:
		return Number(v), nil
	case uint8:
		return Number(v), nil
	case uint16:
		return Number(v), nil
	case uint32:
		return Number(v), nil
	case uint64:
		return Number(v), nil
	case int:
		return Number(v), nil
	case int8:
		return Number(v), nil
	case int16:
		return Number(v), nil
	case int32:
		return Number(v), nil
	case int64:
		return Number(v), nil
	case uintptr:
		return Number(v), nil
	case string:
		return String(v), nil
	case time.Time:
		return Time(v), nil
	case []any:
		// TODO: potential bug:
		// this case would trigger only for []any, not, for example, []string
		// this can be worked around using reflection
		dst := make(List, len(v))
		for i, e := range v {
			d, err := ParseAny(e)
			if err != nil {
				return nil, err
			}
			dst[i] = d
		}
		return dst, nil
	case map[string]any:
		return ParseMapAny(v)
	default:
		return nil, fmt.Errorf("unsupported data type %T", v)
	}
}

func ParseMapAny(v map[string]any) (Map, error) {
	if v == nil {
		return nil, nil
	}
	dst := make(Map)
	for k, e := range v {
		d, err := ParseAny(e)
		if err != nil {
			return nil, err
		}
		dst[k] = d
	}
	return dst, nil
}

type Convertible interface {
	AsPluginData() Data
}

func IsTruthy(d Data) bool {
	switch d := d.(type) {
	case Bool:
		return bool(d)
	case Number:
		return float64(d) != 0
	case String:
		return string(d) != ""
	case List:
		return len(d) > 0
	case Map:
		return len(d) > 0
	case nil:
		return false
	default:
		slog.Debug("unsupported data type", "dt", d)
		return false
	}
}
