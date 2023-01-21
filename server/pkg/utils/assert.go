package utils

import (
	"testing"

	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
)

func AssertDecimal(t *testing.T, expected, actual decimal.Decimal) bool {
	res := expected.Equal(actual)
	return assert.Truef(t, res, "decimal %s != %s", expected.String(), actual.String())
}
