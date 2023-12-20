package plugin_b

import (
	"weave-cli/pkg/jsontools"
	"weave-cli/plugins/data"
)

type Impl struct{}

var _ data.Plugin = (*Impl)(nil)

func (Impl) Execute(input any) (result any, err error) {
	var inputParsed struct {
		ParameterZ any `json:"parameter_z"`
	}
	err = jsontools.UnmarshalBytes(input, &inputParsed)
	if err != nil {
		return
	}
	return inputParsed.ParameterZ, nil
}
