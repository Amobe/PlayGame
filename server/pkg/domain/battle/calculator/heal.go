package calculator

import (
	"github.com/shopspring/decimal"

	"github.com/Amobe/PlayGame/server/pkg/domain/vo"
)

type healerAttribute struct {
	matk      decimal.Decimal
	skillRate decimal.Decimal
	amp       decimal.Decimal
	healRate  decimal.Decimal
}

func BuildHealerAttribute(skillAttr, characterAttr vo.AttributeMap) healerAttribute {
	return healerAttribute{
		matk:      characterAttr.Get(vo.AttributeTypeMATK).Value,
		skillRate: skillAttr.Get(vo.AttributeTypeSDR).Value,
		amp:       characterAttr.Get(vo.AttributeTypeAMP).Value,
		healRate:  characterAttr.Get(vo.AttributeTypeHR).Value,
	}
}

func CalculateHeal(ha healerAttribute) decimal.Decimal {
	skillFactor := healSkillFactor(ha)
	increaseFactor := healIncreaseFactor(ha)
	return ha.matk.Mul(skillFactor).Mul(increaseFactor)
}

func healSkillFactor(ha healerAttribute) decimal.Decimal {
	return ha.skillRate.Add(ha.amp)
}

func healIncreaseFactor(ha healerAttribute) decimal.Decimal {
	factor := decimal.NewFromInt(1).Add(ha.healRate)
	if factor.IsNegative() {
		return decimal.Zero
	}
	return factor
}
