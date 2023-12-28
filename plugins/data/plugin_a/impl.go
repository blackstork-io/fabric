package plugin_a //nolint: revive,stylecheck // temporary name

import (
	"github.com/blackstork-io/fabric/pkg/jsontools"
	"github.com/blackstork-io/fabric/plugins/data"
)

// Actual implementation of the plugin

type Impl struct{}

var _ data.Plugin = (*Impl)(nil)

func (Impl) Execute(input any) (result any, err error) {
	var inputParsed struct {
		ParameterX int64 `json:"parameter_x"`
		ParameterY int64 `json:"parameter_y"`
	}
	err = jsontools.UnmarshalBytes(input, &inputParsed)
	if err != nil {
		return
	}

	return inputParsed.ParameterX + inputParsed.ParameterY, nil
}
