package battle

import (
	"fmt"

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
	AllyMinions  vo.Camp
	EnemyMinions vo.Camp
	Status       MinionSlotStatus

	calculateAttackDamageFn
	targetPickerFn
}

var (
	defaultCalculateAttackDamageFn = calculateAttackDamage
	defaultTargetPickerFn          = targetPickerFromFirst
)

func NewMinionSlot(allyMinions, enemyMinions vo.Camp) *MinionSlot {
	return &MinionSlot{
		AllyMinions:  allyMinions,
		EnemyMinions: enemyMinions,
		Status:       MinionSlotStatusStarted,
	}
}

func (s *MinionSlot) PlayOneRound() ([]vo.Affect, error) {
	var affects []vo.Affect
	actionOrder, err := s.getActionOrder()
	if err != nil {
		return nil, fmt.Errorf("get action order error: %w", err)
	}
	for _, actorIdx := range actionOrder {
		affect, err := s.act(actorIdx)
		if err != nil {
			return nil, fmt.Errorf("act error: %w", err)
		}
		affects = append(affects, affect...)
		if s.Status != MinionSlotStatusStarted {
			break
		}
	}
	return affects, nil
}

func (s *MinionSlot) act(actorIdx vo.GroundIdx) ([]vo.Affect, error) {
	enemySummoner, err := s.enemySummoner()
	if err != nil {
		return nil, fmt.Errorf("get enemy summoner error: %w", err)
	}
	allySummoner, err := s.allySummoner()
	if err != nil {
		return nil, fmt.Errorf("get ally summoner error: %w", err)
	}
	attacker, targets, err := s.getAttackerAndTargets(actorIdx)
	if err != nil {
		return nil, fmt.Errorf("get attacker and targets error: %w", err)
	}
	if attacker.IsDead() {
		return nil, nil
	}
	var affects []vo.Affect
	for _, target := range targets {
		affect := s.attack(attacker, target)
		affects = append(affects, affect)
	}
	if enemySummoner.IsDead() {
		s.Status = MinionSlotStatusAllyWon
	}
	if allySummoner.IsDead() {
		s.Status = MinionSlotStatusEnemyWon
	}
	return affects, nil
}

func (s *MinionSlot) getActionOrder() ([]vo.GroundIdx, error) {
	// assume ally summoner is faster than enemy summoner, action idx is starting from 1 to 5
	startActionIdx := 1
	enemySummoner, err := s.enemySummoner()
	if err != nil {
		return nil, fmt.Errorf("get enemy summoner error: %w", err)
	}
	allySummoner, err := s.allySummoner()
	if err != nil {
		return nil, fmt.Errorf("get ally summoner error: %w", err)
	}
	if enemySummoner.GetAgi() > allySummoner.GetAgi() {
		// if enemy summoner is faster than ally summoner, action idx is starting from 7 to 11
		startActionIdx = 7
	}
	actionOrder := make([]vo.GroundIdx, 0, 10)
	for i := startActionIdx; i < startActionIdx+5; i++ {
		firstActorIdx := vo.GroundIdx(i)
		actionOrder = append(actionOrder, firstActorIdx, firstActorIdx.GetOppositeIdx())
	}
	return actionOrder, nil
}

func (s *MinionSlot) getAttackerAndTargets(idx vo.GroundIdx) (attacker vo.Character, targets []vo.Character, err error) {
	getAttackerFn := s.AllyMinions.Get
	getTargetFn := s.getTargetsFn(s.EnemyMinions)
	if idx.IsEnemy() {
		getAttackerFn = s.EnemyMinions.Get
		getTargetFn = s.getTargetsFn(s.AllyMinions)
	}

	attacker, err = getAttackerFn(idx.ToCampIdx())
	if err != nil {
		return vo.EmptyCharacter, nil, fmt.Errorf("get attacker error: %w", err)
	}
	if attacker.IsDead() {
		return attacker, nil, nil
	}
	targetNumber := int(attacker.GetSkill().AttributeMap.Get(vo.AttributeTypeTarget).Value.IntPart())
	targets, err = getTargetFn(targetNumber)
	return attacker, targets, nil
}

func (s *MinionSlot) getTargetsFn(camp vo.Camp) func(number int) (targets []vo.Character, err error) {
	return func(number int) (targets []vo.Character, err error) {
		targets = make([]vo.Character, 0, number)
		randIdx := s.getTargetPickerFn()(1, 5, number)
		for _, idx := range randIdx {
			target, err := camp.Get(vo.CampIdx(idx))
			if err != nil {
				return nil, fmt.Errorf("get target error: %w", err)
			}
			if target.IsDead() {
				target, err = camp.GetSummoner()
				if err != nil {
					return nil, fmt.Errorf("get summoner error: %w", err)
				}
			}
			targets = append(targets, target)
		}
		return targets, nil
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

func (s *MinionSlot) allySummoner() (vo.Character, error) {
	return s.AllyMinions.GetSummoner()
}

func (s *MinionSlot) enemySummoner() (vo.Character, error) {
	return s.EnemyMinions.GetSummoner()
}
