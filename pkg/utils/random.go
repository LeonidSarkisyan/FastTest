package utils

import (
	"crypto/sha1"
	"encoding/hex"
	"math/rand/v2"
	"path/filepath"
	"slices"
	"time"
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

func GenerateUniqueFilename(originalFilename string) string {
	hash := sha1.New()
	hash.Write([]byte(time.Now().String() + originalFilename))
	sha1Hash := hex.EncodeToString(hash.Sum(nil))
	return sha1Hash + filepath.Ext(originalFilename)
}
