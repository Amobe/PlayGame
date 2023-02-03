package calculator

import (
	"github.com/shopspring/decimal"

	"github.com/Amobe/PlayGame/server/pkg/domain/vo"
	"github.com/Amobe/PlayGame/server/pkg/utils"
)

type damageType string

const (
	damageTypePhysical damageType = "p"
	damageTypeMagical  damageType = "m"
)

type damageAttackerAttribute struct {
	atk         decimal.Decimal
	matk        decimal.Decimal
	skillRate   decimal.Decimal
	amp         decimal.Decimal
	attackerCri decimal.Decimal
	skillCri    decimal.Decimal
	criD        decimal.Decimal
	di          decimal.Decimal
	hit         decimal.Decimal
}

func BuildDamageAttackerAttribute(skillAttr, characterAttr vo.AttributeMap) damageAttackerAttribute {
	return damageAttackerAttribute{
		atk:         characterAttr.Get(vo.AttributeTypeATK).Value,
		matk:        characterAttr.Get(vo.AttributeTypeMATK).Value,
		skillRate:   skillAttr.Get(vo.AttributeTypeSDR).Value,
		amp:         characterAttr.Get(vo.AttributeTypeAMP).Value,
		attackerCri: characterAttr.Get(vo.AttributeTypeCRI).Value,
		skillCri:    skillAttr.Get(vo.AttributeTypeCRI).Value,
		criD:        characterAttr.Get(vo.AttributeTypeCRID).Value,
		di:          characterAttr.Get(vo.AttributeTypeDI).Value,
		hit:         characterAttr.Get(vo.AttributeTypeHit).Value,
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

func CalculateDamage(damageType damageType, daa damageAttackerAttribute, dta damageTargetAttribute) (damage decimal.Decimal, targetDodge bool) {
	if isTargetDodge(daa.hit, dta.dodge) {
		return decimal.Zero, true
	}

	randomFactor := randDamageFactor()

	baseFactor := baseDamageFactor(damageType, daa, dta)

	skillFactor := skillDamageFactor(daa, dta)

	criticalFactor := criticalDamageFactor(daa, dta)

	damageFactor := daa.di.Sub(dta.dR)

	// baseFactor * skillFactor * criticalFactor * randomFactor + damageFactor
	return baseFactor.Mul(skillFactor).Mul(criticalFactor).Mul(randomFactor).Add(damageFactor).Round(0), false
}

func baseDamageFactor(damageType damageType, daa damageAttackerAttribute, dta damageTargetAttribute) decimal.Decimal {
	atk, def := daa.atk, dta.def
	if damageType == damageTypeMagical {
		atk, def = daa.matk, dta.mdef
	}
	numerator := atk.Mul(atk)
	denominator := atk.Add(def.Mul(decimal.NewFromInt(2)))
	if denominator.IsZero() {
		return decimal.Zero
	}
	// (atk * atk) / (atk + 2*def)
	return numerator.Div(denominator)
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
		MininumFactor = decimal.NewFromFloat(1.25)
	)

	if !isCritical(daa.attackerCri, dta.criR, daa.skillCri) {
		return BaseFactor
	}
	criticalFactor := BaseFactor.Add(daa.criD).Sub(dta.criDR)
	if criticalFactor.LessThan(MininumFactor) {
		return MininumFactor
	}
	return criticalFactor
}

func isCritical(cri, criR, sCri decimal.Decimal) bool {
	criticalRate := cri.Add(sCri).Sub(criR)
	return utils.GetProbabilitySampling(criticalRate.InexactFloat64())
}

func isTargetDodge(attackerHit, targetDodge decimal.Decimal) bool {
	hitRate := calculateHitRate(attackerHit, targetDodge)
	isHit := utils.GetProbabilitySampling(hitRate.InexactFloat64())
	return !isHit
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
