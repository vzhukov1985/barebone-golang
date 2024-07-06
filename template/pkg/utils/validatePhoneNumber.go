package utils

import "regexp"

func ValidatePhoneNumber(phoneNumber string) bool {
	var validPhone = regexp.MustCompile(`(^8|7|\+7)((\d{10})|(\s\(\d{3}\)\s\d{3}\s\d{2}\s\d{2}))`)

	return validPhone.MatchString(phoneNumber)
}
