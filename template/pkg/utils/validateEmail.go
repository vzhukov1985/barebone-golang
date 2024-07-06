package utils

import (
	"errors"
	"github.com/badoux/checkmail"
	"net"
	"strings"
)

var blackList = []string{
	"hi2.in",
	"@free.fr",
}

func getHost(email string) string {
	i := strings.LastIndexByte(email, '@')
	return email[i+1:]
}

func validateHost(host string) error {
	_, err := net.LookupMX(host)
	if err != nil {
		return errors.New("unresolvable host")
	}
	return nil
}

func ValidateEmail(email string) bool {

	err := checkmail.ValidateFormat(email)
	if err != nil {
		return false
	}

	host := getHost(email)

	err = validateHost(host)
	if err != nil {
		return false
	}

	for _, v := range blackList {
		v := v
		if strings.EqualFold(host, v) {
			return false
		}
	}

	return true
}
