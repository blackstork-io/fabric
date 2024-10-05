package astv1

import (
	"errors"
)

// bubbleError is a wrapper for errors that are meant to be panic'ed
// and caught by the higher level caller.
// This is considered to be a bad practice, but return-error propagation
// is not practical in deeply nested encoding/decoding functions.
// Standard library uses this pattern in some places, i.e. encoding/json.
type bubbleError struct {
	err error
}

// bubbleWrap wraps an error in a bubbleError and returns it (not panics).
func bubbleWrap(err error) bubbleError {
	return bubbleError{err: err}
}

// bubbleUp wraps an error in a bubbleError and panics it.
func bubbleUp(err error) {
	panic(bubbleWrap(err))
}

// recoverBubbleError expects to be called in a defer statement with
// the current error and a recover(). If bubbleError is recovered, it
// would be [errors.Join]'ed with the current error (if non nil) and
// returned.
// If there was no panic, the passed in error would be returned.
// If there is any other panic, it would be re-panicked.
func recoverBubbleError(err error, recovered any) error {
	if recovered == nil {
		return err
	}
	if pErr, ok := recovered.(bubbleError); ok {
		if err != nil {
			return errors.Join(err, pErr.err)
		}
		return pErr.err
	}
	panic(recovered)
}
