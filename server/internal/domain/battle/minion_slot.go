package battle

import (
	"github.com/shopspring/decimal"

	"github.com/Amobe/PlayGame/server/internal/domain/battle/calculator"
	"github.com/Amobe/PlayGame/server/internal/domain/vo"
)

type MinionSlotStatus string

const (
	MinionSlotStatusStarted  MinionSlotStatus = "started"
	MinionSlotStatusAllyWon  MinionSlotStatus = "ally_won"
	MinionSlotStatusEnemyWon MinionSlotStatus = "enemy_won"
)

type calculateAttackDamageFn func(attacker, target vo.Character, skill vo.Skill) (damage decimal.Decimal, hit bool)

type targetPickerFn func(min, max, number int) []int

type MinionSlot struct {
	AllyMinions  *Minions
	EnemyMinions *Minions
	Status       MinionSlotStatus

	calculateAttackDamageFn
	targetPickerFn
}

var (
	defaultCalculateAttackDamageFn = calculateAttackDamage
	defaultTargetPickerFn          = targetPickerFromFirst
)

func NewMinionSlot(allyMinions, enemyMinions *Minions) *MinionSlot {
	return &MinionSlot{
		AllyMinions:  allyMinions,
		EnemyMinions: enemyMinions,
		Status:       MinionSlotStatusStarted,
	}
}

func (s *MinionSlot) PlayOneRound() []vo.Affect {
	var affects []vo.Affect
	actionOrder := s.getActionOrder()
	for _, actorIdx := range actionOrder {
		affects = append(affects, s.act(actorIdx)...)
		if s.Status != MinionSlotStatusStarted {
			break
		}
	}
	return affects
}

func (s *MinionSlot) act(actorIdx vo.GroundIdx) []vo.Affect {
	attacker, targets := s.getAttackerAndTargets(actorIdx)
	if attacker.IsDead() {
		return nil
	}
	var affects []vo.Affect
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

func (s *MinionSlot) getActionOrder() []vo.GroundIdx {
	// assume ally summoner is faster than enemy summoner, action idx is starting from 1 to 5
	startActionIdx := 1
	if s.enemySummoner().GetAgi() > s.allySummoner().GetAgi() {
		// if enemy summoner is faster than ally summoner, action idx is starting from 7 to 11
		startActionIdx = 7
	}
	actionOrder := make([]vo.GroundIdx, 0, 10)
	for i := startActionIdx; i < startActionIdx+5; i++ {
		firstActorIdx := vo.GroundIdx(i)
		actionOrder = append(actionOrder, firstActorIdx, firstActorIdx.GetOppositeIdx())
	}
	return actionOrder
}

func (s *MinionSlot) getAttackerAndTargets(idx vo.GroundIdx) (attacker vo.Character, targets []vo.Character) {
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

func (s *MinionSlot) getTargetsFn(minions *Minions) func(number int) (targets []vo.Character) {
	return func(number int) (targets []vo.Character) {
		targets = make([]vo.Character, 0, number)
		randIdx := s.getTargetPickerFn()(1, 5, number)
		for _, idx := range randIdx {
			target := minions.Get(vo.CampIdx(idx))
			if target.IsDead() {
				target = minions.GetSummoner()
			}
			targets = append(targets, target)
		}
		return targets
	}
}

func (s *MinionSlot) getTargetPickerFn() targetPickerFn {
	if s.targetPickerFn != nil {
		return s.targetPickerFn
	}
	return defaultTargetPickerFn
}

func targetPickerFromFirst(min, max, number int) []int {
	var res []int
	for i := min; i <= max && len(res) < number; i++ {
		res = append(res, i)
	}
	return res
}

func (s *MinionSlot) attack(attacker, target vo.Character) vo.Affect {
	skill := attacker.GetSkill()
	damage, isHit := s.getCalculateAttackDamageFn()(attacker, target, skill)
	if !isHit {
		return vo.NewMissAffect(attacker.GetGroundIdx(), target.GetGroundIdx())
	}
	affects := vo.NewAttributeMap(
		vo.Attribute{
			Type:  vo.AttributeTypeDamage,
			Value: damage,
		},
	)
	s.unitTakeAffect(target, affects)
	return vo.NewAffect(attacker.GetGroundIdx(), target.GetGroundIdx(), skill.Name, affects)
}

func (s *MinionSlot) getCalculateAttackDamageFn() calculateAttackDamageFn {
	if s.calculateAttackDamageFn != nil {
		return s.calculateAttackDamageFn
	}
	return defaultCalculateAttackDamageFn
}

func calculateAttackDamage(attacker, target vo.Character, skill vo.Skill) (damage decimal.Decimal, isHit bool) {
	attackerAttribute := attacker.GetAttributeMap()
	skillAttribute := skill.AttributeMap
	targetAttribute := target.GetAttributeMap()
	da := calculator.BuildDamageAttackerAttribute(skillAttribute, attackerAttribute)
	dt := calculator.BuildDamageTargetAttribute(targetAttribute)
	return calculator.CalculateDamage(da, dt)
}

func (s *MinionSlot) unitTakeAffect(unit vo.Character, affects vo.AttributeMap) {
	groundIdx := unit.GetGroundIdx()
	campIDx := groundIdx.ToCampIdx()
	minionSetFn := s.AllyMinions.Set
	if groundIdx.IsEnemy() {
		minionSetFn = s.EnemyMinions.Set
	}
	minionSetFn(campIDx, unit.TakeAffect(affects))
}

func (s *MinionSlot) allySummoner() vo.Character {
	return s.AllyMinions.GetSummoner()
}

func (s *MinionSlot) enemySummoner() vo.Character {
	return s.EnemyMinions.GetSummoner()
}
