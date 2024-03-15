package utils

import (
	"math/rand/v2"
)

func GenerateSixDigitNumber() int64 {
	return rand.Int64N(900000) + 100000
}
