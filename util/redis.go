package util

import "fmt"

// GetKeyRedis builds a namespaced redis key in the form "key:value".
func GetKeyRedis(key string, value string) string {
	return fmt.Sprintf("%s:%s", key, value)
}
