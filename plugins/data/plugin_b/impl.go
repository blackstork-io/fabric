package plugin_b //nolint: revive,stylecheck // temporary name

import (
	"github.com/blackstork-io/fabric/pkg/jsontools"
	"github.com/blackstork-io/fabric/plugins/data"
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
