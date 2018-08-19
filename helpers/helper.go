package helpers

import (
	"os"
	"sync"
)

// getEnv return env variable or default value provided
func GetEnv(name, defaultV string) string {
	if value, ok := os.LookupEnv(name); ok {
		return value
	}

	return defaultV
}

func SyncMapLen(m *sync.Map) int {
	l := 0
	m.Range(func(key, value interface{}) bool {
		l++
		return true
	})

	return l
}
