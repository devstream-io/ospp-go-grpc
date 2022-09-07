package util

import "strings"

func GenKey(keys ...string) string {
	return strings.Join(keys, ":")
}
