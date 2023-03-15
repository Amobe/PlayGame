package battle

import (
	"github.com/shopspring/decimal"

	"github.com/Amobe/PlayGame/server/pkg/domain/battle/calculator"
	"github.com/Amobe/PlayGame/server/pkg/domain/vo"
)

//go:generate mockery --name Unit --inpackage
type Unit interface {
	GetGroundIdx() GroundIdx
	GetAttributeMap() vo.AttributeMap
	GetAgi() int
	IsDead() bool
	TakeAffect(affects []vo.Attribute) Unit
	GetSkill() vo.Skill
}

// GroundIdx is the index on the battleground.
// The battleground contains the ally and enemy minions.
// The GroundIdx for ally minions is between 1 and 5.
// The GroundIdx for ally summoner is 6.
// The GroundIdx for enemy minions is between 7 and 11.
// The GroundIdx for enemy summoner is 12.
type GroundIdx int

var (
	AllySummonerGroundIdx  = GroundIdx(6)
	EnemySummonerGroundIdx = GroundIdx(12)
)

func (g GroundIdx) ToCampIdx() CampIdx {
	if g > AllySummonerGroundIdx {
		return CampIdx(g - 6)
	}
	return CampIdx(g)
}

func (g GroundIdx) IsEnemy() bool {
	return g > 6
}

func (g GroundIdx) GetOppositeIdx() GroundIdx {
	if g.IsEnemy() {
		return g - 6
	}
	return g + 6
}

func (g GroundIdx) ToInt32() int32 {
	return int32(g)
}

// CampIdx is the index of the camp slots.
// The camp slots only contain the minions with same camp.
// The CampIdx of the minions is between 1 and 5.
// The CampIdx of the summoner is 6.
type CampIdx int

var SummonerCampIdx = CampIdx(6)

func (c CampIdx) IsSummoner() bool {
	return c == SummonerCampIdx
}

type Minions [6]Unit

func (m *Minions) Get(idx CampIdx) Unit {
	return m[idx-1]
}

func (m *Minions) Set(idx CampIdx, unit Unit) {
	m[idx-1] = unit
}

func (m *Minions) GetSummoner() Unit {
	return m.Get(SummonerCampIdx)
}

func (m *Minions) SetSummoner(summoner Unit) {
	m.Set(SummonerCampIdx, summoner)
}

type MinionSlotStatus string

const (
	MinionSlotStatusStarted  MinionSlotStatus = "started"
	MinionSlotStatusAllyWon  MinionSlotStatus = "ally_won"
	MinionSlotStatusEnemyWon MinionSlotStatus = "enemy_won"
)

type calculateAttackDamageFn func(attacker, target Unit, skill vo.Skill) (damage decimal.Decimal, hit bool)

type targetPickerFn func(min, max, number int) []int

type MinionSlot struct {
	AllyMinions  Minions
	EnemyMinions Minions
	Status       MinionSlotStatus

	calculateAttackDamageFn
	targetPickerFn
}

var (
	defaultCalculateAttackDamageFn = calculateAttackDamage
	defaultTargetPickerFn          = targetPickerFromFirst
)

func NewMinionSlot(allyMinions, enemyMinions Minions) *MinionSlot {
	return &MinionSlot{
		AllyMinions:  allyMinions,
		EnemyMinions: enemyMinions,
		Status:       MinionSlotStatusStarted,

		calculateAttackDamageFn: defaultCalculateAttackDamageFn,
		targetPickerFn:          defaultTargetPickerFn,
	}
}

func (s *MinionSlot) PlayOneRound() []Affect {
	var affects []Affect
	actionOrder := s.getActionOrder()
	for _, actorIdx := range actionOrder {
		affects = append(affects, s.act(actorIdx)...)
		if s.Status != MinionSlotStatusStarted {
			break
		}
	}
	return affects
}

func (s *MinionSlot) act(actorIdx GroundIdx) []Affect {
	attacker, targets := s.getAttackerAndTargets(actorIdx)
	if attacker.IsDead() {
		return nil
	}
	var affects []Affect
	for _, target := range targets {
		affect := s.attack(attacker, target)
		affects = append(affects, affect)
	}
	if s.enemySummoner().IsDead() {
		s.Status = MinionSlotStatusAllyWon
	}
	if s.allySummoner().IsDead() {
		s.Status = MinionSlotStatusEnemyWon
	}
	return affects
}

func (s *MinionSlot) getActionOrder() []GroundIdx {
	// assume ally summoner is faster than enemy summoner, action idx is starting from 1 to 5
	startActionIdx := 1
	if s.enemySummoner().GetAgi() > s.allySummoner().GetAgi() {
		// if enemy summoner is faster than ally summoner, action idx is starting from 7 to 11
		startActionIdx = 7
	}
	actionOrder := make([]GroundIdx, 0, 10)
	for i := startActionIdx; i < startActionIdx+5; i++ {
		firstActorIdx := GroundIdx(i)
		actionOrder = append(actionOrder, firstActorIdx, firstActorIdx.GetOppositeIdx())
	}
	return actionOrder
}

func (s *MinionSlot) getAttackerAndTargets(idx GroundIdx) (attacker Unit, targets []Unit) {
	getAttackerFn := s.AllyMinions.Get
	getTargetFn := s.getTargetsFn(s.EnemyMinions)
	if idx.IsEnemy() {
		getAttackerFn = s.EnemyMinions.Get
		getTargetFn = s.getTargetsFn(s.AllyMinions)
	}

	attacker = getAttackerFn(idx.ToCampIdx())
	if attacker.IsDead() {
		return attacker, nil
	}
	targetNumber := int(attacker.GetSkill().AttributeMap.Get(vo.AttributeTypeTarget).Value.IntPart())
	targets = getTargetFn(targetNumber)
	return attacker, targets
}

func (s *MinionSlot) getTargetsFn(minions Minions) func(number int) (targets []Unit) {
	return func(number int) (targets []Unit) {
		targets = make([]Unit, 0, number)
		randIdx := s.targetPickerFn(1, 5, number)
		for _, idx := range randIdx {
			target := minions.Get(CampIdx(idx))
			if target.IsDead() {
				target = minions.GetSummoner()
			}
			targets = append(targets, target)
		}
		return targets
	}
}

func targetPickerFromFirst(min, max, number int) []int {
	var res []int
	for i := min; i <= max && len(res) < number; i++ {
		res = append(res, i)
	}
	return res
}

func (s *MinionSlot) attack(attacker Unit, target Unit) Affect {
	skill := attacker.GetSkill()
	damage, isHit := s.calculateAttackDamageFn(attacker, target, skill)
	if !isHit {
		return NewMissAffect(attacker.GetGroundIdx(), target.GetGroundIdx())
	}
	affects := []vo.Attribute{
		{
			Type:  vo.AttributeTypeDamage,
			Value: damage,
		},
	}
	s.unitTakeAffect(target, affects)
	return NewAffect(attacker.GetGroundIdx(), target.GetGroundIdx(), skill.Name, affects)
}

func calculateAttackDamage(attacker, target Unit, skill vo.Skill) (damage decimal.Decimal, isHit bool) {
	attackerAttribute := attacker.GetAttributeMap()
	skillAttribute := skill.AttributeMap
	targetAttribute := target.GetAttributeMap()
	da := calculator.BuildDamageAttackerAttribute(skillAttribute, attackerAttribute)
	dt := calculator.BuildDamageTargetAttribute(targetAttribute)
	return calculator.CalculateDamage(da, dt)
}

func (s *MinionSlot) getUnit(groundIdx GroundIdx) Unit {
	campIdx := groundIdx.ToCampIdx()
	if groundIdx.IsEnemy() {
		return s.EnemyMinions.Get(campIdx)
	}
	return s.AllyMinions.Get(campIdx)
}

func (s *MinionSlot) unitTakeAffect(unit Unit, affects []vo.Attribute) {
	groundIdx := unit.GetGroundIdx()
	campIDx := groundIdx.ToCampIdx()
	minionSetFn := s.AllyMinions.Set
	if groundIdx.IsEnemy() {
		minionSetFn = s.EnemyMinions.Set
	}
	minionSetFn(campIDx, unit.TakeAffect(affects))
}

func (s *MinionSlot) allySummoner() Unit {
	return s.AllyMinions.GetSummoner()
}

func (s *MinionSlot) enemySummoner() Unit {
	return s.EnemyMinions.GetSummoner()
}
