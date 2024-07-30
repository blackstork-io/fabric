package circularRefDetector

import (
	"sync"
	"unsafe"

	"github.com/hashicorp/hcl/v2"

	"github.com/blackstork-io/fabric/pkg/diagnostics"
)

type circularRefDetector struct {
	mu   sync.Mutex
	refs map[unsafe.Pointer]*hcl.Range
}

var detector = circularRefDetector{
	refs: make(map[unsafe.Pointer]*hcl.Range),
}

// Adds pointer to circular reference detection.
// refRange is optional.
func Add[T any](ptr *T, refRange *hcl.Range) {
	//nolint:gosec // We are just using pointers as unique ids, not accessing the data.
	detector.add(unsafe.Pointer(ptr), refRange)
}

// Checks if the passed in ptr was added previously.
func Check[T any](ptr *T) bool {
	//nolint:gosec // We are just using pointers as unique ids, not accessing the data.
	return detector.check(unsafe.Pointer(ptr))
}

// Removes pointer from circular reference detection and
// appends the backtrace to the diagnostic marked by ExtraMarker
// (if it exists).
func Remove[T any](ptr *T, diags *diagnostics.Diag) {
	//nolint:gosec // We are just using pointers as unique ids, not accessing the data.
	detector.remove(unsafe.Pointer(ptr), diags)
}

func (d *circularRefDetector) add(ptr unsafe.Pointer, refRange *hcl.Range) {
	d.mu.Lock()
	d.refs[ptr] = refRange
	d.mu.Unlock()
}

func (d *circularRefDetector) check(ptr unsafe.Pointer) (found bool) {
	d.mu.Lock()
	_, found = d.refs[ptr]
	d.mu.Unlock()
	return
}

func (d *circularRefDetector) remove(ptr unsafe.Pointer, diags *diagnostics.Diag) {
	d.mu.Lock()
	rng, found := d.refs[ptr]
	delete(d.refs, ptr)
	d.mu.Unlock()
	if !found || diags == nil {
		return
	}

	tb, _ := diagnostics.DiagnosticsGetExtra[diagnostics.TracebackExtra](*diags)
	if tb == nil {
		return
	}
	tb.Traceback = append(tb.Traceback, rng)
}
