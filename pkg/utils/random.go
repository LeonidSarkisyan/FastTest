package utils

import (
	"math/rand/v2"
	"slices"
)

func GenerateSixDigitNumber(n int) []int64 {
	uniqueNumbers := make([]int64, n)

	var index int

	for n > 0 {
		r := rand.Int64N(900000) + 100000
		index = slices.Index(uniqueNumbers, r)
		if index == -1 {
			uniqueNumbers[n-1] = r
			n -= 1
		}
	}

	return uniqueNumbers
}
