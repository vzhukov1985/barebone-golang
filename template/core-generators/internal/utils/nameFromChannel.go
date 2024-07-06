package utils

import "strings"

func NameFromChannel(ch string) string {
	spStrs := strings.Split(ch, "/")

	res := ""
	for _, v := range spStrs {
		if len(v) > 0 {
			res += UpperCaseFirstLetter(v)
		}
	}
	return res
}
