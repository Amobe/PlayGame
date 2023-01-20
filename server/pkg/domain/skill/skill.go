package skill

import (
	"strconv"

	"github.com/Amobe/PlayGame/server/pkg/domain/valueobject"
	"github.com/Amobe/PlayGame/server/pkg/utils"
	"github.com/Amobe/PlayGame/server/pkg/utils/domain"
)

var _ domain.Aggregator = &Skill{}

type coreAggregator = domain.CoreAggregator

type iSkill interface {
	Use(am, dm valueobject.AttributeTypeMap) (aa, ta []valueobject.Attribute)
	Name() string
}

type Skill struct {
	coreAggregator
	SkillID      string
	name         string
	AttributeMap valueobject.AttributeTypeMap
}

func (s Skill) ID() string {
	return s.SkillID
}

func (s Skill) Name() string {
	return s.name
}

type SkillEmpty struct {
	Skill
}

func NewSkillEmpty() SkillEmpty {
	return SkillEmpty{
		Skill{
			SkillID:      "empty",
			name:         "empty",
			AttributeMap: nil,
		},
	}
}

func (s SkillEmpty) Use(am, dm valueobject.AttributeTypeMap) (aa, ta []valueobject.Attribute) {
	return
}

type SkillPoisonHit struct {
	Skill
}

func NewSkillPoisonHit() SkillPoisonHit {
	am := valueobject.NewAttributeTypeMap()
	am.Insert(valueobject.Attribute{Type: valueobject.AttributeTypeSDR, Value: "1.1"})
	am.Insert(valueobject.Attribute{Type: valueobject.AttributeTypeStatusH, Value: "0.5"})
	return SkillPoisonHit{
		Skill: Skill{
			SkillID:      "poisonHit",
			name:         "poisonHit",
			AttributeMap: am,
		},
	}
}

func (h SkillPoisonHit) Use(am, dm valueobject.AttributeTypeMap) (aa, ta []valueobject.Attribute) {
	atk := am[valueobject.AttributeTypeATK].GetFloat()
	def := dm[valueobject.AttributeTypeDEF].GetFloat()
	sdr := h.AttributeMap[valueobject.AttributeTypeSDR].GetFloat()
	amp := am[valueobject.AttributeTypeAMP].GetFloat()
	ampR := dm[valueobject.AttributeTypeAMPR].GetFloat()
	dI := am[valueobject.AttributeTypeDI].GetFloat()
	dR := dm[valueobject.AttributeTypeDR].GetFloat()
	cri := am[valueobject.AttributeTypeCRI].GetFloat()
	sCri := h.AttributeMap[valueobject.AttributeTypeCRI].GetFloat()
	criR := dm[valueobject.AttributeTypeCRIR].GetFloat()
	criD := am[valueobject.AttributeTypeCRID].GetFloat()
	criDR := dm[valueobject.AttributeTypeCRIDR].GetFloat()

	// physical damage
	damage := damageCal(atk, 1, def, sdr, amp, ampR, cri, criR, sCri, criD, criDR, dI, dR)
	value := strconv.Itoa(damage * -1)
	ta = append(ta, valueobject.Attribute{Type: valueobject.AttributeTypeHP, Value: value})

	// poison hit
	sh := am[valueobject.AttributeTypeStatusH].GetFloat()
	shR := dm[valueobject.AttributeTypeStatusHR].GetFloat()
	sSh := h.AttributeMap[valueobject.AttributeTypeStatusH].GetFloat()
	if isStatusHit(sh, shR, sSh) {
		ta = append(ta, valueobject.Attribute{Type: valueobject.AttributeTypePoisoned, Value: "1"})
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

func damageCal(atk, atkBonus, def, skillRate, amp, ampR, cri, criR, sCri, criD, criDR, dI, dR float64) int {
	randomFactor := randDamageFactor()

	finalATK := atk * atkBonus

	baseFactor := (finalATK * finalATK) / (finalATK + 2*def)

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
