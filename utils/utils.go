package utils

import "os"

func IsFileExist(filePath string) bool {
	_, err := os.Stat(filePath)
	return err == nil
}

func GetOrDefault[K comparable, V any](m map[K]V, key K, defaultValue V) V {
	if v, ok := m[key]; ok {
		return v
	}
	return defaultValue
}
