package utils

import "strings"

func CorrectCase(s string) string {
	if strings.HasSuffix(s, "ID") {
		s = strings.TrimSuffix(s, "ID")
		return s + "Id"
	}
	return s
}
