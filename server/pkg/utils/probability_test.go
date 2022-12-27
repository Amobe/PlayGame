package utils_test

import (
	"math"
	"testing"

	"github.com/Amobe/PlayGame/server/pkg/utils"
	"github.com/stretchr/testify/assert"
)

func BenchmarkGetProbabilitySampling(b *testing.B) {
	rate := 0.5
	for n := 0; n < b.N; n++ {
		utils.GetProbabilitySampling(rate)
	}
}

func TestGetProbabilitySamplingAcurracy(t *testing.T) {
	rate := 0.555
	sample := 10000000
	successCnt := 0

	for i := 0; i < sample; i++ {
		if utils.GetProbabilitySampling(rate) {
			successCnt++
		}
	}

	successRate := float64(successCnt) / float64(sample)

	difference := math.Abs(rate - successRate)
	assert.Truef(t, difference < 0.0005, "difference %f must less than 0.0005", difference)
}
