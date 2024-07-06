package utils

import (
	"crypto/md5"
	"fmt"
	"github.com/goccy/go-json"
)

func Md5Hash(v interface{}) (string, error) {
	s, err := json.Marshal(v)
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("%x", md5.Sum(s)), nil
}
