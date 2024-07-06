package utils

import (
	"fmt"
	"math/rand"
	"time"
)

func GetRandomInt(min, max int64) int {

	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	maxRandom := int(max - min)

	if maxRandom <= 0 {
		return int(min)
	}

	return r.Intn(maxRandom) + int(min)
}

func GetRandomIntWithSalt(min, max, salt int64) int {

	r := rand.New(rand.NewSource(salt + time.Now().UnixNano()))
	maxRandom := int(max - min)

	if maxRandom <= 0 {
		return int(min)
	}

	return r.Intn(maxRandom) + int(min)
}

func RandomCode(digits int) string {
	var retval string
	for k := 0; k < digits; k++ {
		retval += fmt.Sprintf("%d", GetRandomIntWithSalt(0, 9, int64(GetRandomInt(0, 1e9))))
	}

	return retval
}
