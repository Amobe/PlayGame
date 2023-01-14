package utils

func CopyMap[K comparable, V any](dst map[K]V, src map[K]V) {
	for k, v := range src {
		dst[k] = v
	}
}
