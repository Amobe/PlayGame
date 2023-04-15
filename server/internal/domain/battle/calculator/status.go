package calculator

import (
	"github.com/shopspring/decimal"

	"github.com/Amobe/PlayGame/server/internal/utils"
)

type statusAttackerAttribute struct {
	statusH      decimal.Decimal
	skillStatusH decimal.Decimal
}

type statusTargetAttribute struct {
	statusR decimal.Decimal
}

func CalculateStatusHit(sa statusAttackerAttribute, st statusTargetAttribute) bool {
	hitRate := calculateStatusHitRate(sa, st)
	isHit := utils.GetProbabilitySampling(hitRate.InexactFloat64())
	return isHit
}

func calculateStatusHitRate(sa statusAttackerAttribute, st statusTargetAttribute) decimal.Decimal {
	statusHitRate := sa.statusH.Add(sa.skillStatusH).Sub(st.statusR)
	if statusHitRate.IsNegative() {
		return decimal.Zero
	}
	return utils.GetDecimalWithUpperBound(statusHitRate, decimal.NewFromInt(1))
}
