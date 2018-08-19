package helpers

import (
	"os"
)

// getEnv return env variable or default value provided
func GetEnv(name, defaultV string) string {
	if value, ok := os.LookupEnv(name); ok {
		return value
	}

	return defaultV
}
