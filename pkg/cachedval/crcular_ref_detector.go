package cachedval

import (
	"github.com/hashicorp/hcl/v2"

	"github.com/blackstork-io/fabric/pkg/backtrace"
	"github.com/blackstork-io/fabric/pkg/utils"
)

var _ backtrace.ExecTracer = (*CircularRefDetector)(nil)

type CircularRefDetector struct {
	refs map[*backtrace.Backtracer]int
}

func NewCRD() *CircularRefDetector {
	return &CircularRefDetector{
		refs: map[*backtrace.Backtracer]int{},
	}
}

func (crd *CircularRefDetector) FrameEnter(bt *backtrace.Backtracer) (diag *hcl.Diagnostic) {
	_, found := crd.refs[bt]
	if !found {
		crd.refs[bt] = len(crd.refs)
		return
	}
	trace := make([]*backtrace.Backtracer, len(crd.refs))
	for btStep, frameNo := range crd.refs {
		// reverse the order of the backtrace
		trace[(len(trace)-1)-frameNo] = btStep
	}
	diag = (*bt).NewDiagnostic()
	for i, btStep := range trace {
		(*btStep).AppendBacktrace(diag)
		if btStep == bt && i != len(trace)-1 {
			diag.Detail += "\nReference loop entered because of:"
		}
	}
	return
}

func (crd *CircularRefDetector) FrameExit(bt *backtrace.Backtracer) {
	utils.Pop(crd.refs, bt)
}
