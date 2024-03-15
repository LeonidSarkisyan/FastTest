package utils

import (
	"math/big"
	"math/rand/v2"
	"time"
)

func RandomBigNumber(number int) int64 {
	// Преобразуем входное число в *big.Int
	max := big.NewInt(int64(number))

	// Генерируем случайное большое число в заданном диапазоне
	randomBigNumber := new(big.Int).Rand(rand.New(rand.NewS(time.Now().UnixNano())), max)

	return randomBigNumber
}

func Read(p []byte) {
	for i := 0; i < len(p); {
		val := rand.Uint64()
		for j := 0; j < 8 && i < len(p); j++ {
			p[i] = byte(val & 0xff)
			val >>= 8
			i++
		}
	}
}
