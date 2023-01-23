package character

import (
	"github.com/shopspring/decimal"

	"github.com/Amobe/PlayGame/server/pkg/domain/vo"
	"github.com/Amobe/PlayGame/server/pkg/utils"
)

type Skill interface {
	Use(am, dm vo.AttributeMap) (aa, ta []vo.Attribute)
	Name() string
}

type skill struct {
	SkillID      string
	name         string
	AttributeMap vo.AttributeMap
}

func (s skill) Name() string {
	return s.name
}

type SkillEmpty struct {
	skill
}

func NewSkillEmpty() SkillEmpty {
	return SkillEmpty{
		skill{
			SkillID:      "empty",
			name:         "empty",
			AttributeMap: nil,
		},
	}
}

func (s SkillEmpty) Use(am, dm vo.AttributeMap) (aa, ta []vo.Attribute) {
	return
}

type SkillPoisonHit struct {
	skill
}

func NewSkillPoisonHit() SkillPoisonHit {
	am := vo.NewAttributeMap(
		vo.NewAttribute(vo.AttributeTypeSDR, decimal.NewFromFloat(1.1)),
		vo.NewAttribute(vo.AttributeTypeStatusH, decimal.NewFromFloat(0.5)),
	)
	return SkillPoisonHit{
		skill: skill{
			SkillID:      "poisonHit",
			name:         "poisonHit",
			AttributeMap: am,
		},
	}
}

func (h SkillPoisonHit) Use(am, dm vo.AttributeMap) (aa, ta []vo.Attribute) {
	atk := am.Get(vo.AttributeTypeATK).Value
	def := dm.Get(vo.AttributeTypeDEF).Value
	sdr := h.AttributeMap.Get(vo.AttributeTypeSDR).Value
	amp := am.Get(vo.AttributeTypeAMP).Value
	ampR := dm.Get(vo.AttributeTypeAMPR).Value
	dI := am.Get(vo.AttributeTypeDI).Value
	dR := dm.Get(vo.AttributeTypeDR).Value
	cri := am.Get(vo.AttributeTypeCRI).Value
	sCri := h.AttributeMap.Get(vo.AttributeTypeCRI).Value
	criR := dm.Get(vo.AttributeTypeCRIR).Value
	criD := am.Get(vo.AttributeTypeCRID).Value
	criDR := dm.Get(vo.AttributeTypeCRIDR).Value

	// physical damage
	damage := damageCal(atk, def, sdr, amp, ampR, cri, criR, sCri, criD, criDR, dI, dR)
	ta = append(ta, vo.NewAttribute(vo.AttributeTypeHP, damage.Neg()))

	// poison hit
	sh := am.Get(vo.AttributeTypeStatusH).Value
	shR := dm.Get(vo.AttributeTypeStatusHR).Value
	sSh := h.AttributeMap.Get(vo.AttributeTypeStatusH).Value
	if isStatusHit(sh, shR, sSh) {
		ta = append(ta, vo.NewAttribute(vo.AttributeTypePoisoned, decimal.NewFromInt(1)))
	}

	return nil, ta
}

var SkillMap = map[string]Skill{
	"poisonHit": NewSkillPoisonHit(),
}

func isCritical(cri, criR, sCri decimal.Decimal) bool {
	criticalRate := cri.Add(sCri).Sub(criR)
	return utils.GetProbabilitySampling(criticalRate.InexactFloat64())
}

func isStatusHit(sh, shR, sSh decimal.Decimal) bool {
	hitRate := sh.Add(sSh).Sub(shR)
	return utils.GetProbabilitySampling(hitRate.InexactFloat64())
}

func damageCal(atk, def, skillRate, amp, ampR, cri, criR, sCri, criD, criDR, dI, dR decimal.Decimal) decimal.Decimal {
	randomFactor := randDamageFactor()

	// baseFactor = (atk * atk) / (atk + 2*def)
	baseFactor := atk.Mul(atk).Div(atk.Add(def.Mul(decimal.NewFromInt(2))))

	skillFactor := skillDamageFactor(skillRate, amp, ampR)

	criticalFactor := criticalDamageFactor(cri, criR, sCri, criD, criDR)

	damageFactor := dI.Sub(dR)

	return baseFactor.Mul(skillFactor).Mul(criticalFactor).Mul(randomFactor).Add(damageFactor).Round(0)
}

func skillDamageFactor(damageRate, amp, ampR decimal.Decimal) decimal.Decimal {
	skillIncrease := amp.Sub(ampR)
	if skillIncrease.IsPositive() {
		skillIncrease = decimal.Zero
	}
	return damageRate.Add(skillIncrease)
}

func randDamageFactor() decimal.Decimal {
	factor := utils.GetRandFloatInRange(0.8, 1.3)
	return decimal.NewFromFloat(factor)
}

func criticalDamageFactor(cri, criR, sCri, criD, criDR decimal.Decimal) decimal.Decimal {
	var (
		BaseFactor    = decimal.NewFromInt(1)
		MininumFactor = decimal.NewFromFloat(1.25)
	)

	if !isCritical(cri, criR, sCri) {
		return BaseFactor
	}
	criticalFactor := BaseFactor.Add(criD).Sub(criDR)
	if criticalFactor.LessThan(MininumFactor) {
		return MininumFactor
	}
	return criticalFactor
}
