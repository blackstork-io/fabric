package fabctx

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"os/signal"
	"time"
)

type FabCtx struct {
	mainCtx    context.Context
	cleanupCtx context.Context
	linting    bool
}

var _ context.Context = (*FabCtx)(nil)

func (ctx *FabCtx) Deadline() (deadline time.Time, ok bool) {
	return ctx.mainCtx.Deadline()
}

func (ctx *FabCtx) Done() <-chan struct{} {
	return ctx.mainCtx.Done()
}

func (ctx *FabCtx) Err() error {
	return ctx.mainCtx.Err()
}

type getFabCtxT struct{}

var getFabCtx = getFabCtxT{}

func Get(ctx context.Context) *FabCtx {
	if ctx == nil {
		return nil
	}
	if fc, ok := ctx.(*FabCtx); ok {
		return fc
	}
	fc := ctx.Value(getFabCtx)
	if fc == nil {
		return nil
	}
	if fc, ok := ctx.(*FabCtx); ok {
		return fc
	}
	return nil
}

func (ctx *FabCtx) Value(v any) any {
	switch v.(type) {
	case getFabCtxT:
		return ctx
	default:
		// fabCtx is the root context
		return nil
	}
}

func (ctx *FabCtx) CleanupCtx() context.Context {
	if ctx == nil {
		slog.Warn("CleanupCtx was called on a nil ptr!")
		return context.Background()
	}
	return ctx.cleanupCtx
}

func (ctx *FabCtx) IsLinting() bool {
	if ctx == nil {
		slog.Warn("IsLinting was called on a nil ptr!")
		return false
	}
	return ctx.linting
}

type fabCtxOpts struct {
	signals bool
}

type Option func(*fabCtxOpts)

func NoSignals(opts *fabCtxOpts) {
	opts.signals = false
}

func WithLinting(parent *FabCtx) *FabCtx {
	ctx := *parent
	ctx.linting = true
	return &ctx
}

// Returns a cli-appropriate context (cancelable by ctrl+c).
func New(options ...Option) *FabCtx {
	opts := fabCtxOpts{
		signals: true,
	}
	for _, opt := range options {
		opt(&opts)
	}

	if !opts.signals {
		return &FabCtx{
			mainCtx:    context.Background(),
			cleanupCtx: context.Background(),
		}
	}

	var (
		ctx                       FabCtx
		mainCancel, cleanupCancel context.CancelCauseFunc
	)
	ctx.cleanupCtx, cleanupCancel = context.WithCancelCause(context.Background())
	ctx.mainCtx, mainCancel = context.WithCancelCause(ctx.cleanupCtx)

	c := make(chan os.Signal, 2)
	signal.Notify(c, os.Interrupt)

	go func() {
		caught := 0
		for range c {
			switch caught {
			case 0:
				slog.Warn("Received os.Interrupt")
				mainCancel(fmt.Errorf("got termination request (gentle)"))
			case 1:
				slog.Error("Received second os.Interrupt")
				cleanupCancel(fmt.Errorf("got termination request (forceful)"))
			default:
				slog.Error("Rough exit (3 interrupts received, probably deadlocked)")
				os.Exit(1)
			}
			caught++
		}
	}()
	return &ctx
}
