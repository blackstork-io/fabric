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

func FnMap[I, O any](fn func(I) O, s []I) []O {
	out := make([]O, len(s))
	for i, v := range s {
		out[i] = fn(v)
	}
	return out
}
