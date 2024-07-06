package utils

import "strings"

func UpperCaseFirstLetter(s string) string {
	return strings.ToUpper(s[:1]) + s[1:]
}
