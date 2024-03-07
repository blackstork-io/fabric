package lint

import "context"

type lintCtxKeyT struct{}

var lintCtxKey = lintCtxKeyT{}

func IsLinting(ctx context.Context) bool {
	return ctx.Value(lintCtxKey) != nil
}

func MakeLintContext(ctx context.Context) context.Context {
	return context.WithValue(ctx, lintCtxKey, struct{}{})
}
