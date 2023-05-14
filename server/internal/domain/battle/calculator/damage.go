package calculator

import (
	"github.com/shopspring/decimal"

	"github.com/Amobe/PlayGame/server/internal/domain/vo"
	"github.com/Amobe/PlayGame/server/internal/utils"
)

type damageType string

const (
	damageTypePhysical damageType = "p"
	damageTypeMagical  damageType = "m"
	damageTypeHybrid   damageType = "h"
)

type damageAttackerAttribute struct {
	atk         decimal.Decimal
	atkB        decimal.Decimal
	matk        decimal.Decimal
	matkB       decimal.Decimal
	skillRate   decimal.Decimal
	amp         decimal.Decimal
	attackerCri decimal.Decimal
	skillCri    decimal.Decimal
	criD        decimal.Decimal
	di          decimal.Decimal
	hit         decimal.Decimal
	damageType  damageType
}

func BuildDamageAttackerAttribute(skillAttr, characterAttr vo.AttributeMap) damageAttackerAttribute {
	damageType := damageTypePhysical
	switch {
	case skillAttr.GetBool(vo.AttributeTypeMagicalDamage):
		damageType = damageTypeMagical
	case skillAttr.GetBool(vo.AttributeTypeHybridDamage):
		damageType = damageTypeHybrid
	}
	return damageAttackerAttribute{
		atk:         characterAttr.Get(vo.AttributeTypeATK).Value,
		atkB:        skillAttr.Get(vo.AttributeTypeATKB).Value,
		matk:        characterAttr.Get(vo.AttributeTypeMATK).Value,
		matkB:       skillAttr.Get(vo.AttributeTypeMATKB).Value,
		skillRate:   skillAttr.Get(vo.AttributeTypeSDR).Value,
		amp:         characterAttr.Get(vo.AttributeTypeAMP).Value,
		attackerCri: characterAttr.Get(vo.AttributeTypeCRI).Value,
		skillCri:    skillAttr.Get(vo.AttributeTypeCRI).Value,
		criD:        characterAttr.Get(vo.AttributeTypeCRID).Value,
		di:          characterAttr.Get(vo.AttributeTypeDI).Value,
		hit:         characterAttr.Get(vo.AttributeTypeHit).Value,
		damageType:  damageType,
	}
}

type damageTargetAttribute struct {
	def   decimal.Decimal
	mdef  decimal.Decimal
	ampR  decimal.Decimal
	criR  decimal.Decimal
	criDR decimal.Decimal
	dR    decimal.Decimal
	dodge decimal.Decimal
}

func BuildDamageTargetAttribute(characterAttr vo.AttributeMap) damageTargetAttribute {
	return damageTargetAttribute{
		def:   characterAttr.Get(vo.AttributeTypeDEF).Value,
		mdef:  characterAttr.Get(vo.AttributeTypeMDEF).Value,
		ampR:  characterAttr.Get(vo.AttributeTypeAMPR).Value,
		criR:  characterAttr.Get(vo.AttributeTypeCRIR).Value,
		criDR: characterAttr.Get(vo.AttributeTypeCRIDR).Value,
		dR:    characterAttr.Get(vo.AttributeTypeDR).Value,
		dodge: characterAttr.Get(vo.AttributeTypeDodge).Value,
	}
}

func CalculateDamage(daa damageAttackerAttribute, dta damageTargetAttribute) (damage decimal.Decimal, hit bool) {
	if !isHit(daa.hit, dta.dodge) {
		return decimal.Zero, false
	}

	randomFactor := randDamageFactor()

	baseFactor := baseDamageFactor(daa, dta)

	skillFactor := skillDamageFactor(daa, dta)

	criticalFactor := criticalDamageFactor(daa, dta)

	damageFactor := daa.di.Sub(dta.dR)

	// baseFactor * skillFactor * criticalFactor * randomFactor + damageFactor
	return baseFactor.Mul(skillFactor).Mul(criticalFactor).Mul(randomFactor).Add(damageFactor).Round(0), true
}

func baseDamageFactor(daa damageAttackerAttribute, dta damageTargetAttribute) decimal.Decimal {
	formula := func(atk, def decimal.Decimal) decimal.Decimal {
		numerator := atk.Mul(atk)
		denominator := atk.Add(def.Mul(decimal.NewFromInt(2)))
		if denominator.IsZero() {
			return decimal.Zero
		}
		// (atk * atk) / (atk + 2*def)
		return numerator.Div(denominator)
	}

	atk, def := daa.atk.Mul(daa.atkB), dta.def
	matk, mdef := daa.matk.Mul(daa.matkB), dta.mdef
	switch daa.damageType {
	case damageTypePhysical:
		return formula(atk, def)
	case damageTypeMagical:
		return formula(matk, mdef)
	case damageTypeHybrid:
		return formula(atk, def).Add(formula(matk, mdef))
	}
	return decimal.Zero
}

func skillDamageFactor(daa damageAttackerAttribute, dta damageTargetAttribute) decimal.Decimal {
	skillIncrease := daa.amp.Sub(dta.ampR)
	if skillIncrease.IsPositive() {
		skillIncrease = decimal.Zero
	}
	return daa.skillRate.Add(skillIncrease)
}

func randDamageFactor() decimal.Decimal {
	factor := utils.GetRandFloatInRange(0.8, 1.3)
	return decimal.NewFromFloat(factor)
}

func criticalDamageFactor(daa damageAttackerAttribute, dta damageTargetAttribute) decimal.Decimal {
	var (
		BaseFactor    = decimal.NewFromInt(1)
		MinimumFactor = decimal.NewFromFloat(1.25)
	)

	if !isCritical(daa.attackerCri, dta.criR, daa.skillCri) {
		return BaseFactor
	}
	criticalFactor := BaseFactor.Add(daa.criD).Sub(dta.criDR)
	if criticalFactor.LessThan(MinimumFactor) {
		return MinimumFactor
	}
	return criticalFactor
}

func isCritical(cri, criR, sCri decimal.Decimal) bool {
	criticalRate := cri.Add(sCri).Sub(criR)
	return utils.GetProbabilitySampling(criticalRate.InexactFloat64())
}

func isHit(attackerHit, targetDodge decimal.Decimal) bool {
	hitRate := calculateHitRate(attackerHit, targetDodge)
	return utils.GetProbabilitySampling(hitRate.InexactFloat64())
}

func calculateHitRate(attackerHit, targetDodge decimal.Decimal) decimal.Decimal {
	var (
		baseHit        = decimal.NewFromFloat(1.9)
		roundUpPlace   = int32(4)
		maximumHitRate = decimal.NewFromInt(1)
	)

	hit := utils.GetNonNegativeDecimal(attackerHit)
	dodge := utils.GetNonNegativeDecimal(targetDodge)

	if hit.Equal(dodge) {
		return baseHit.Div(decimal.NewFromInt(2))
	}
	numerator := hit
	denominator := hit.Add(dodge)
	hitRate := baseHit.Mul(numerator).Div(denominator).Round(roundUpPlace)
	return utils.GetDecimalWithUpperBound(hitRate, maximumHitRate)
}
