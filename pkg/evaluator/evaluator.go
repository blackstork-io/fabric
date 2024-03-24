package evaluator

import (
	"context"
	"sync"

	"github.com/hashicorp/hcl/v2"

	"github.com/blackstork-io/fabric/pkg/diagnostics"
)

/*
Every request (render, data, lint, etc) has its own auto-incrementing atomic id

When request needs to perform a calculation, it stores its id into a field on the node.

Idea of the waiting-for:


When request encounters a node in the "calculating" state, it checks its ID:
	- if it matches current request ID - look through the frame stack to check for the ref loops
	(in SAAS)
	- if it doesn't match: use a global id->id map to check for ref loops between requests
		- if ref loop found – break it by cancelling the newest request in the loop, commanding it to rerun
		  after the oldest is finished


*/



type CachedValStatus uint32

const (
	statusUninit  = CachedValStatus(0)
	statusSuccess = CachedValStatus(1 << (iota - 1))
	statusFaliure
	statusCalculating
)

func (s CachedValStatus) IsCalculating() bool {
	return s&statusCalculating != 0
}

func (s CachedValStatus) ExclusiveStatuses() CachedValStatus {
	return s & (^statusCalculating)
}

type TracebackLinkedList struct {
	parent *TracebackLinkedList
	// using it to store arbitraty function pointers
	frame FrameIdentifier
}

func (ll *TracebackLinkedList) Contains(frame FrameIdentifier) bool {
	for ll != nil {
		if frame.Eq(ll.frame) {
			return true
		}
		ll = ll.parent
	}
	return false
}

func (ll *TracebackLinkedList) CircularRefDetector(frame FrameIdentifier) *hcl.Diagnostic {
	loopLength := 0
	// search for the looping node
	for curNode := ll; curNode != nil; curNode = curNode.parent {
		loopLength++
		if frame.Eq(curNode.frame) {
			diag := TracebackDiag{
				loopLength: loopLength + 1,
				traceback:  []FrameIdentifier{frame},
			}
			// dump traceback
			for curNode := ll; curNode != nil; curNode = curNode.parent {
				diag.traceback = append(diag.traceback, curNode.frame)
			}
			return &hcl.Diagnostic{
				Severity: hcl.DiagError,
				Summary:  "Circular reference detected",
				// TODO: format error in the printer
				// Detail:   "Looped back to this object through reference chain:",
				// Subject:  rng,
				Extra: &diag,
			}
		}
	}
	return nil
}

type FrameIdentifier interface {
	Eq(FrameIdentifier) bool
}

type CachedValInfo struct {
	// extra info (for example no-cache markers, etc)
	extra     any
	calcCount uint32
	status    CachedValStatus
}

type CachedVal[V any] struct {
	val V
	CachedValInfo
	cond sync.Cond
	_mu  sync.Mutex
}

func NewCachedVal[V any]() *CachedVal[V] {
	cv := &CachedVal[V]{}
	cv.cond = sync.Cond{L: &cv._mu}
	return cv
}

// Wrapper over sync.Cond.Wait that adds context cancellability for Wait operation
// Condition lock must be held before calling this function and will be held after it returns
func CondWait(ctx context.Context, cond *sync.Cond, until func() bool) (ctxCancelled bool) {
	if until() {
		return true
	}

	// could be an atomic, but context cancelation happens only once while condition
	// broadcasts can happen more often, so non-atomic read in the loop is preferable.
	// ctxCancelled changed under lock cond.L
	stop := context.AfterFunc(ctx, func() {
		cond.L.Lock()
		ctxCancelled = true
		cond.L.Unlock()
		cond.Broadcast()
	})

	for {
		cond.Wait()
		if ctxCancelled {
			return
		}
		if until() {
			// no need to stop in the ctxCancelled case
			// ctxCancelled will never be set, the lock is held by us even after the return
			stop()
			return
		}
	}
}

type TracebackDiag struct {
	traceback  []FrameIdentifier
	loopLength int
}

type Operation[K any, V any] interface {
	Get(e *eval, input *K) (val V, diags diagnostics.Diag)
}

type NoCache[K any, V any] struct {
	calculate func(eval, *K) (V, diagnostics.Diag)
}

type FrameId[K, V any] struct {
	key       *K
	operation Operation[K, V]
}

func (cv FrameId[K, V]) Eq(other FrameIdentifier) bool {
	otherCv, ok := other.(FrameId[K, V])
	return ok && otherCv == cv
}

func getRange(v any) (rng *hcl.Range) {
	switch t := v.(type) {
	case *hcl.Block:
		rng = &t.DefRange
		rng2 := getRange(t.Body)
		if rng2 != nil {
			r := hcl.RangeBetween(*rng, *rng2)
			rng = &r
		}
		return rng
	case *hcl.Attribute:
		return &t.Range
	default:
		// covers *hclsyntax.Block, *hclsyntax.Attribute, *hclsyntax.Body
		type Ranger interface {
			Range() hcl.Range
		}
		ranger, ok := t.(Ranger)
		if !ok {
			return nil
		}
		r := ranger.Range()
		return &r
	}
}

func (nc *NoCache[K, V]) Get(e *eval, input *K) (val V, diags diagnostics.Diag) {
	frame := FrameId[K, V]{
		key:       input,
		operation: nc,
	}

	if diags.Append(e.evalList.CircularRefDetector(frame)) {
		return
	}
	return nc.calculate(e.WithFrame(frame), input)
}

type OperationCache[K, V any] struct {
	m         map[*K]*CachedVal[V]
	calculate func(eval, *K) (V, diagnostics.Diag)
	mu        sync.Mutex
}

// functions getting *eval must not change it
func (tc *OperationCache[K, V]) Get(e *eval, input *K) (val V, diags diagnostics.Diag) {
	// get or create the cached val in the map
	tc.mu.Lock()
	cv, found := tc.m[input]
	if !found {
		cv = NewCachedVal[V]()
		tc.m[input] = cv
	}
	tc.mu.Unlock()

	// evaluate the cached value
	cv.cond.L.Lock()
	defer cv.cond.L.Unlock()
	initCalcCount := cv.calcCount
	for {
		if cv.status.IsCalculating() && initCalcCount == cv.calcCount {
			// Wait for the calculated value once, then return the most recent result regardless of the calc bit set
			frame := FrameId[K, V]{
				key:       input,
				operation: tc,
			}
			if diags.Append(e.evalList.CircularRefDetector(frame)) {
				return
			}

			lastCalcCount := cv.calcCount
			if CondWait(e.Context, &cv.cond, func() bool { return lastCalcCount != cv.calcCount }) {
				return // ctx cancelled
			}
			// else - we already waited for the calculation once, just use the most recent cache
		}
		switch cv.status.ExclusiveStatuses() {
		case statusSuccess:
			val = cv.val
			return
		case statusFaliure:
			diags.Append(diagnostics.RepeatedError)
			return
		case statusUninit:
			cv.status |= statusCalculating
			// setting success var explicitly because calculation might panic
			success := false
			defer func() {
				// unlocked by the initial defer
				cv.cond.L.Lock()
				if success {
					cv.status = statusSuccess
					cv.val = val
				} else {
					cv.status = statusFaliure
				}
			}()
			cv.cond.L.Unlock()
			val, diags = tc.calculate(
				e.WithFrame(FrameId[K, V]{
					key:       input,
					operation: tc,
				}),
				input,
			)
			success = !diags.HasErrors()
			return

		default:
			panic("Must be exhaustive")
		}
		// TODO: rerunable block logic goes here. Track last calculation date?
	}
}
