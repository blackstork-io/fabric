package utils

func Contains[K comparable, V any](m map[K]V, key K) bool {
	_, found := m[key]
	return found
}

func SliceToSet[K comparable](slice []K) map[K]struct{} {
	res := make(map[K]struct{}, len(slice))
	for _, v := range slice {
		res[v] = struct{}{}
	}
	return res
}
