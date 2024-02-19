package clicontext

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"os/signal"
)

type cleanupCtxKeyT struct{}

var cleanupCtxKey = cleanupCtxKeyT{}

// Returns the cleanup context given the main context.
// Use cleanup context in, for example, deferred statements
func GetCleanupCtx(ctx context.Context) context.Context {
	cleanupCtx := ctx.Value(cleanupCtxKey)
	if cleanupCtx == nil {
		return context.Background()
	}
	return cleanupCtx.(context.Context)
}

// Returns a cli-appropriate context (cancelable by ctrl+c).
func New() context.Context {
	cleanupCtx, cleanupCancel := context.WithCancelCause(context.Background())
	valCtx := context.WithValue(cleanupCtx, cleanupCtxKey, cleanupCtx)
	mainCtx, mainCancel := context.WithCancelCause(valCtx)

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
	return mainCtx
}
