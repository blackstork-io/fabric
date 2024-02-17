package evaluation

import (
	"maps"

	"github.com/blackstork-io/fabric/plugin"
)

type DataContext struct {
	plugin.ConvMapData
	mapOwned bool
}

// passed-in map must not be modified afterwards
func NewDataContext(m plugin.ConvMapData) DataContext {
	return DataContext{
		ConvMapData: m,
		mapOwned:    true,
	}
}

func (dc *DataContext) Share() DataContext {
	dc.mapOwned = false
	return *dc
}

func (dc *DataContext) makeOwned() {
	dc.ConvMapData = maps.Clone(dc.ConvMapData)
	dc.mapOwned = true
}

func (dc *DataContext) Delete(key string) {
	if _, found := dc.ConvMapData[key]; !found {
		return
	}
	if !dc.mapOwned {
		dc.makeOwned()
	}
	delete(dc.ConvMapData, key)
}

func (dc *DataContext) Set(key string, val plugin.ConvertableData) {
	if !dc.mapOwned {
		dc.makeOwned()
	}
	dc.ConvMapData[key] = val
}
