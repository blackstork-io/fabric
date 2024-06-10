package plugin

import (
	"github.com/blackstork-io/fabric/pkg/encapsulator"
)

var EncapsulatedData = encapsulator.NewCodec[Data]("data", nil)
