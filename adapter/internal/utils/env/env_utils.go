package envutils

import "os"

// GetEnv lookup environment variable with key,
// if not defined returns default value
func GetEnv(key, defaultValue string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return defaultValue
}
