package jsontools

import (
	"encoding/json"
	"fmt"
)

func MapSet(m map[string]any, keys []string, val any) (map[string]any, error) {
	curMap := m
	if len(keys) == 0 {
		return m, fmt.Errorf("MapSet: keys empty")
	}
	if len(keys) > 1 {
		for _, k := range keys[:len(keys)-1] {
			v, found := curMap[k]
			if found {
				var ok bool
				curMap, ok = v.(map[string]any)
				if !ok {
					return m, fmt.Errorf("MapSet: failed to cast to map[string]any")
				}
			} else {
				nextMap := map[string]any{}
				curMap[k] = nextMap
				curMap = nextMap
			}
		}
	}
	curMap[keys[len(keys)-1]] = val
	return m, nil
}

func MapGet(m any, keys []string) (val any, err error) {
	if len(keys) == 0 {
		err = fmt.Errorf("MapGet: keys empty")
		return
	}
	val = m
	for _, k := range keys {
		asMap, ok := val.(map[string]any)
		if !ok {
			err = fmt.Errorf("MapGet: failed to cast to map[string]any")
			return
		}
		val, ok = asMap[k]
		if !ok {
			err = fmt.Errorf("can't find a key")
			return
		}
	}
	return
}

func Dump(obj any) string {
	objBytes, err := json.Marshal(obj)
	if err != nil {
		objBytes = []byte(fmt.Sprintf("Failed to dump the object as json: %s", err))
	}
	return string(objBytes)
}

func UnmarshalBytes(bytes, value any) error {
	data, ok := bytes.([]byte)
	if !ok {
		return fmt.Errorf("expected array of bytes")
	}
	return json.Unmarshal(data, value)
}
