package cachedval

import (
	"github.com/hashicorp/hcl/v2"

	"github.com/blackstork-io/fabric/pkg/backtrace"
)

type Evaluator interface {
	DefineConfig(hcl.Block)
	DefineData()
	DefineContext()
	DefineDocument()
	DefineSection()
	DefineGlobalConfig()

	// Produce config value
	EvaluateConfig()
	EvaluateConfigAttr()
	// Produce ref base block
	EvaluateBaseAttr()

	backtrace.ExecTracer
}
