package utils

import "strings"

func LowerCaseFirstLetter(s string) string {
	return strings.ToLower(s[:1]) + s[1:]
}
