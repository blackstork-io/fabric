package once

import (
	"sync/atomic"
)

type status int32

const (
	statusPanic      = status(-1)
	statusUnexecuted = status(0)
	statusSuccess    = status(1)
)

func (s status) IsSuccess() bool {
	return s > 0
}

func (s status) IsPanic() bool {
	return s < 0
}

func (s status) IsUnexecuted() bool {
	return s == 0
}

func (s *status) AtomicLoad() status {
	return status(atomic.LoadInt32((*int32)(s)))
}

func (s *status) AtomicStore(newStatus status) {
	atomic.StoreInt32((*int32)(s), (int32)(newStatus))
}
