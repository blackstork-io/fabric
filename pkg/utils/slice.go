package utils

func ApplyFn[T any](slice []T, fn func(T)) {
	for _, it := range slice {
		fn(it)
	}
}

func ApplyFnRef[T any](slice []T, fn func(*T)) {
	for _, it := range slice {
		fn(&it)
	}
}
