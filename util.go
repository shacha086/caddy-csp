package csp

func SetOrAppend[K comparable, V any](m map[K][]V, k K, elems ...V) bool {
	s, ok := m[k]
	m[k] = append(s, elems...)
	return ok
}
