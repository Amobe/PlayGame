package character

import (
	"strconv"

	"github.com/Amobe/PlayGame/server/pkg/utils"
)

type Skill interface {
	Use(am, dm AttributeTypeMap) (aa, ta []Attribute)
}

type skill struct {
	SkillID      string
	Name         string
	AttributeMap AttributeTypeMap
}

type SkillPoisonHit struct {
	skill
}

func NewSkillPoisonHit() SkillPoisonHit {
	am := NewAttributeTypeMap()
	am.Insert(Attribute{Type: AttributeTypeSDR, Value: "1.1"})
	am.Insert(Attribute{Type: AttributeTypeStatusH, Value: "0.5"})
	return SkillPoisonHit{
		skill: skill{
			SkillID:      "poisonHit",
			Name:         "poisonHit",
			AttributeMap: am,
		},
	}
}

func (h SkillPoisonHit) Use(am, dm AttributeTypeMap) (aa, ta []Attribute) {
	atk := am[AttributeTypeATK].GetFloat()
	def := dm[AttributeTypeDEF].GetFloat()
	sdr := h.AttributeMap[AttributeTypeSDR].GetFloat()
	amp := am[AttributeTypeAMP].GetFloat()
	ampR := dm[AttributeTypeAMPR].GetFloat()
	dI := am[AttributeTypeDI].GetFloat()
	dR := dm[AttributeTypeDR].GetFloat()
	cri := am[AttributeTypeCRI].GetFloat()
	sCri := h.AttributeMap[AttributeTypeCRI].GetFloat()
	criR := dm[AttributeTypeCRIR].GetFloat()
	criD := am[AttributeTypeCRID].GetFloat()
	criDR := dm[AttributeTypeCRIDR].GetFloat()

	// physical damage
	damage := damageCal(atk, def, sdr, amp, ampR, cri, criR, sCri, criD, criDR, dI, dR)
	value := strconv.Itoa(damage * -1)
	ta = append(ta, Attribute{Type: AttributeTypeHP, Value: value})

	// poison hit
	sh := am[AttributeTypeStatusH].GetFloat()
	shR := dm[AttributeTypeStatusHR].GetFloat()
	sSh := h.AttributeMap[AttributeTypeStatusH].GetFloat()
	if isStatusHit(sh, shR, sSh) {
		ta = append(ta, Attribute{Type: AttributeTypePoisoned, Value: "1"})
	}

	return nil, ta
}

var SkillMap = map[string]Skill{
	"poisonHit": NewSkillPoisonHit(),
}

func isCritical(cri, criR, sCri float64) bool {
	criticalRate := cri + sCri - criR
	return utils.GetProbabilitySampling(criticalRate)
}

func isStatusHit(sh, shR, sSh float64) bool {
	hitRate := sh + sSh - shR
	return utils.GetProbabilitySampling(hitRate)
}

func damageCal(atk, def, skillRate, amp, ampR, cri, criR, sCri, criD, criDR, dI, dR float64) int {
	randomFactor := randDamageFactor()

	baseFactor := (atk * atk) / (atk + 2*def)

	skillFactor := skillDamageFactor(skillRate, amp, ampR)

	criticalFactor := criticalDamageFactor(cri, criR, sCri, criD, criDR)

	damangeFactor := dI - dR

	return int(baseFactor*skillFactor*criticalFactor*randomFactor + damangeFactor)
}

func skillDamageFactor(damageRate, amp, ampR float64) float64 {
	skillIncrease := amp - ampR
	if skillIncrease > 0 {
		skillIncrease = 0
	}
	return damageRate + skillIncrease
}

func randDamageFactor() float64 {
	return utils.GetRandFloatInRange(0.8, 1.3)
}

func criticalDamageFactor(cri, criR, sCri, criD, criDR float64) float64 {
	if !isCritical(cri, criR, sCri) {
		return 1
	}
	criticalFactor := 1 + criD - criDR
	if criticalFactor > 1.25 {
		return 1.25
	}
	if criticalFactor < 1 {
		return 1
	}
	return criticalFactor
}
