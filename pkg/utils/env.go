package utils

import "os"

func GetEnv(key string, dflt string) string {
	val, ok := os.LookupEnv(key)
	if !ok {
		return dflt
	}
	return val
}
