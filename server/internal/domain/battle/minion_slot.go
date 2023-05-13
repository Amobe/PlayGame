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

type MinionSlot struct {
	Ground vo.Ground
	Status MinionSlotStatus

	calculateAttackDamageFn
}

var (
	defaultCalculateAttackDamageFn = calculateAttackDamage
)

func NewMinionSlot(ground vo.Ground) *MinionSlot {
	return &MinionSlot{
		Ground: ground,
		Status: MinionSlotStatusStarted,
	}
}

func (s *MinionSlot) PlayOneRound() ([]vo.Affect, error) {
	var affects []vo.Affect
	matchOrder, err := s.Ground.GetMatchOrder()
	if err != nil {
		return nil, fmt.Errorf("get action order: %w", err)
	}
	for _, actorIdx := range matchOrder {
		affect, err := s.act(actorIdx)
		if err != nil {
			return nil, fmt.Errorf("act: %w", err)
		}
		affects = append(affects, affect...)
		if s.Status != MinionSlotStatusStarted {
			break
		}
	}
	return affects, nil
}

func (s *MinionSlot) act(actorIdx vo.GroundIdx) ([]vo.Affect, error) {
	enemySummoner, err := s.Ground.GetEnemySummoner()
	if err != nil {
		return nil, fmt.Errorf("get enemy summoner: %w", err)
	}
	allySummoner, err := s.Ground.GetAllySummoner()
	if err != nil {
		return nil, fmt.Errorf("get ally summoner: %w", err)
	}
	match, err := s.Ground.GetMatch(actorIdx)
	if err != nil {
		return nil, fmt.Errorf("get match: %w", err)
	}
	if match.Actor.IsDead() {
		return nil, nil
	}
	var affects []vo.Affect
	for _, target := range match.Targets {
		affect, err := s.attack(match.Actor, target)
		if err != nil {
			return nil, fmt.Errorf("attack: %w", err)
		}
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

func (s *MinionSlot) attack(attacker, target vo.Character) (vo.Affect, error) {
	skill := attacker.GetSkill()
	damage, isHit := s.getCalculateAttackDamageFn()(attacker, target, skill)
	if !isHit {
		return vo.NewMissAffect(attacker.GetGroundIdx(), target.GetGroundIdx()), nil
	}
	affects := vo.NewAttributeMap(
		vo.Attribute{
			Type:  vo.AttributeTypeDamage,
			Value: damage,
		},
		vo.Attribute{
			Type:  vo.AttributeTypeHP,
			Value: damage.Neg(),
		},
	)
	if err := s.unitTakeAffect(target, affects); err != nil {
		return vo.Affect{}, fmt.Errorf("unit take affect: %w", err)
	}
	return vo.NewAffect(attacker.GetGroundIdx(), target.GetGroundIdx(), skill.Name, affects), nil
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

func (s *MinionSlot) unitTakeAffect(unit vo.Character, affects vo.AttributeMap) error {
	idx := unit.GetGroundIdx()
	newGround, err := s.Ground.UpdateCharacter(idx, unit.TakeAffect(affects))
	if err != nil {
		return fmt.Errorf("update ground: %w", err)
	}
	s.Ground = newGround
	return nil
}

func (s *MinionSlot) ToString() string {
	return fmt.Sprintf("ground: %s, status: %s", s.Ground.ToString(), s.Status)
}
