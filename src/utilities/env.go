package utilities

import (
	"os"
)

// GetEnv function returns the value of an environment variable if it exists,
// otherwise it returns a fallback value.
func GetEnv(key string, fallback string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	} else {
		return fallback
	}
}
