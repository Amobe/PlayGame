package skill

import (
	"strconv"

	"github.com/Amobe/PlayGame/server/pkg/domain/character"
	"github.com/Amobe/PlayGame/server/pkg/utils"
	"github.com/Amobe/PlayGame/server/pkg/utils/domain"
)

var _ domain.Aggregator = &skill{}

type coreAggregator = domain.CoreAggregator

type iSkill interface {
	Use(am, dm character.AttributeTypeMap) (aa, ta []character.Attribute)
	Name() string
}

type skill struct {
	coreAggregator
	SkillID      string
	name         string
	AttributeMap character.AttributeTypeMap
}

func (s skill) ID() string {
	return s.SkillID
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

func (s SkillEmpty) Use(am, dm character.AttributeTypeMap) (aa, ta []character.Attribute) {
	return
}

type SkillPoisonHit struct {
	skill
}

func NewSkillPoisonHit() SkillPoisonHit {
	am := character.NewAttributeTypeMap()
	am.Insert(character.Attribute{Type: character.AttributeTypeSDR, Value: "1.1"})
	am.Insert(character.Attribute{Type: character.AttributeTypeStatusH, Value: "0.5"})
	return SkillPoisonHit{
		skill: skill{
			SkillID:      "poisonHit",
			name:         "poisonHit",
			AttributeMap: am,
		},
	}
}

func (h SkillPoisonHit) Use(am, dm character.AttributeTypeMap) (aa, ta []character.Attribute) {
	atk := am[character.AttributeTypeATK].GetFloat()
	def := dm[character.AttributeTypeDEF].GetFloat()
	sdr := h.AttributeMap[character.AttributeTypeSDR].GetFloat()
	amp := am[character.AttributeTypeAMP].GetFloat()
	ampR := dm[character.AttributeTypeAMPR].GetFloat()
	dI := am[character.AttributeTypeDI].GetFloat()
	dR := dm[character.AttributeTypeDR].GetFloat()
	cri := am[character.AttributeTypeCRI].GetFloat()
	sCri := h.AttributeMap[character.AttributeTypeCRI].GetFloat()
	criR := dm[character.AttributeTypeCRIR].GetFloat()
	criD := am[character.AttributeTypeCRID].GetFloat()
	criDR := dm[character.AttributeTypeCRIDR].GetFloat()

	// physical damage
	damage := damageCal(atk, def, sdr, amp, ampR, cri, criR, sCri, criD, criDR, dI, dR)
	value := strconv.Itoa(damage * -1)
	ta = append(ta, character.Attribute{Type: character.AttributeTypeHP, Value: value})

	// poison hit
	sh := am[character.AttributeTypeStatusH].GetFloat()
	shR := dm[character.AttributeTypeStatusHR].GetFloat()
	sSh := h.AttributeMap[character.AttributeTypeStatusH].GetFloat()
	if isStatusHit(sh, shR, sSh) {
		ta = append(ta, character.Attribute{Type: character.AttributeTypePoisoned, Value: "1"})
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
	if criticalFactor < 1.25 {
		return 1.25
	}
	return criticalFactor
}
