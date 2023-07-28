package util

import "os"

func GetEnvString(key string, defaultVal string) string {
	value := os.Getenv(key)
	if len(value) == 0 {
		value = defaultVal
	}
	return value
}
