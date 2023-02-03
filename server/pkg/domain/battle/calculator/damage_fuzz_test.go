package calculator

import (
	"testing"

	"github.com/shopspring/decimal"
)

func FuzzBaseDamageFactor_Physical(f *testing.F) {
	f.Fuzz(func(t *testing.T, atk, def uint32) {
		aa := damageAttackerAttribute{
			atk: decimal.NewFromInt(int64(atk)),
		}
		ta := damageTargetAttribute{
			def: decimal.NewFromInt(int64(def)),
		}
		_ = baseDamageFactor(damageTypePhysical, aa, ta)
	})
}

func FuzzBaseDamageFactor_Magical(f *testing.F) {
	f.Fuzz(func(t *testing.T, matk, mdef uint32) {
		aa := damageAttackerAttribute{
			matk: decimal.NewFromInt(int64(matk)),
		}
		ta := damageTargetAttribute{
			mdef: decimal.NewFromInt(int64(mdef)),
		}
		_ = baseDamageFactor(damageTypeMagical, aa, ta)
	})
}

func FuzzSkillDamageFactor(f *testing.F) {
	f.Fuzz(func(t *testing.T, amp, sr, ampR float64) {
		aa := damageAttackerAttribute{
			amp:       decimal.NewFromFloat(amp),
			skillRate: decimal.NewFromFloat(sr),
		}
		ta := damageTargetAttribute{
			ampR: decimal.NewFromFloat(ampR),
		}
		_ = skillDamageFactor(aa, ta)
	})
}

func FuzzCalculateHitRate(f *testing.F) {
	f.Fuzz(func(t *testing.T, hit, dodge float64) {
		attackerHit := decimal.NewFromFloat(hit)
		targetDodge := decimal.NewFromFloat(dodge)
		_ = calculateHitRate(attackerHit, targetDodge)
	})
}
