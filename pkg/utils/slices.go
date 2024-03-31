package utils

// Sets slice[idx] = val, growing the slice if needed, and returns the updated slice.
func SetAt[T any](slice []T, idx int, val T) []T {
	needToAlloc := idx - len(slice)
	switch {
	case needToAlloc > 0:
		slice = append(slice, make([]T, needToAlloc)...)

		fallthrough
	case needToAlloc == 0:
		slice = append(slice, val)
	default:
		slice[idx] = val
	}
	return slice
}

// Produce a new slice by applying function fn to items of the slice s
func FnMap[I, O any](s []I, fn func(I) O) []O {
	out := make([]O, len(s))
	for i, v := range s {
		out[i] = fn(v)
	}
	return out
}
