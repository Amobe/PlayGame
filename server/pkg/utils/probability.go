package utils

import (
	"math/rand"
	"time"
)

const sampleSize = 10000000000

func init() {
	rand.Seed(time.Now().UnixNano())
}

func GetProbabilitySampling(rate float64) bool {
	if rate <= 0 {
		return false
	}
	if rate >= 1 {
		return true
	}
	return rand.Intn(sampleSize) <= int(rate*sampleSize)
}

func GetRandFloatInRange(min, max float64) float64 {
	return min + rand.Float64()*(max-min)
}

func GetRandIntInRange(min, max int) int {
	return min + rand.Intn(max-min)
}

func GetRandIntInRangeN(min, max, n int) []int {
	res := make([]int, 0, n)
	for i := 0; i < n; i++ {
		res = append(res, GetRandIntInRange(min, max))
	}
	return res
}
