package utils

import "github.com/shopspring/decimal"

// GetNonNegativeDecimal returns a non negative decimal. If ori is negative, the return decimal is zero.
func GetNonNegativeDecimal(ori decimal.Decimal) decimal.Decimal {
	if ori.IsNegative() {
		return decimal.Zero
	}
	return ori
}

// GetDecimalWithUpperBound returns a decimal which will not greater than upperBound.
func GetDecimalWithUpperBound(ori decimal.Decimal, upperBound decimal.Decimal) decimal.Decimal {
	if ori.GreaterThan(upperBound) {
		return upperBound
	}
	return ori
}
