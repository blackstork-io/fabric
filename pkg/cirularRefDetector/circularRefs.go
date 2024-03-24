package circularRefDetector

import (
	"fmt"
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

type extraMarker struct{}

// Marker to identify circular access diagnostic.
// Identifies the diagnostic that will be extended with the backtrace.
var ExtraMarker = extraMarker{}

// Adds pointer to circular reference detection.
// refRange is optional.
func Add[T any](ptr *T, refRange *hcl.Range) {
	detector.add(unsafe.Pointer(ptr), refRange)
}

// Checks if the passed in ptr was added previously.
func Check[T any](ptr *T) bool {
	return detector.check(unsafe.Pointer(ptr))
}

// Removes pointer from circular reference detection and
// appends the backtrace to the diagnostic marked by ExtraMarker
// (if it exists).
func Remove[T any](ptr *T, diags *diagnostics.Diag) {
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

	diag, _ := diagnostics.FindByExtra[extraMarker](*diags)
	if diag == nil {
		return
	}
	if rng != nil {
		diag.Detail = fmt.Sprintf(
			"%s\n  at %s:%d:%d",
			diag.Detail, rng.Filename, rng.Start.Line, rng.Start.Column,
		)
	} else {
		diag.Detail += "\n  at <missing location info>"
	}
}
