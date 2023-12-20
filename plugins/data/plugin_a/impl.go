package plugin_a

import (
	"weave-cli/pkg/jsontools"
	"weave-cli/plugins/data"
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
