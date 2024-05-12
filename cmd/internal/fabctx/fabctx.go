package fabctx

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"os/signal"
	"time"
)

// FabCtx is a context that can be used to cancel the main context and trigger cleanup.
// It is used to handle graceful shutdowns for the fabric CLI.
type FabCtx struct {
	mainCtx    context.Context
	cleanupCtx context.Context
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
	if fc, ok := fc.(*FabCtx); ok {
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

type fabCtxOpts struct {
	signals bool
}

type Option func(*fabCtxOpts)

func NoSignals(opts *fabCtxOpts) {
	opts.signals = false
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
				slog.WarnContext(&ctx, "Received os.Interrupt")
				mainCancel(fmt.Errorf("got termination request (gentle)"))
			case 1:
				slog.ErrorContext(&ctx, "Received second os.Interrupt")
				cleanupCancel(fmt.Errorf("got termination request (forceful)"))
			default:
				slog.ErrorContext(&ctx, "Rough exit (3 interrupts received, probably deadlocked)")
			}
			caught++
		}
	}()
	return &ctx
}
