package cachedval

import (
	"sync"

	"github.com/hashicorp/hcl/v2"

	"github.com/blackstork-io/fabric/pkg/diagnostics"
)


// Another calc request for this value already reported an error
var RepeatedError = &hcl.Diagnostic{
	Severity: hcl.DiagError,
	Extra:    diagnostics.HiddenError{},
}



// func (cv *CachedVal) Get(crd backtrace.ExecTracer, getter func() (V, diagnostics.Diag)) (val V, diags diagnostics.Diag) {
// 	if diags.Append(crd.FrameEnter(&cv.backtracer)) {
// 		return
// 	}
// 	defer crd.FrameExit(&cv.backtracer)
// 	v, diag := cv.get(getter)
// 	diags.Extend(diag)
// 	return v, diags
// }

// func main() {
// 	var t *backtrace.Backtracer = nil
// 	fmt.Println(unsafe.Sizeof(t))

// 	// d := diagnostics.Diag{
// 	// 	&hcl.Diagnostic{
// 	// 		Extra: CircularRefMarker2,
// 	// 	},
// 	// }
// 	// diagnostics.FindByExtra[T any](diags diagnostics.Diag)
// }
