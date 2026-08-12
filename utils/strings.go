package utils

import "strings"

// ContainsAny 判断 subs 中任一子串是否出现在 src 中。
func ContainsAny(src string, subs []string) bool {
	for _, s := range subs {
		if strings.Contains(src, s) {
			return true
		}
	}
	return false
}
