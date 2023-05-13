package vo

import (
	"fmt"

	"github.com/Amobe/PlayGame/server/internal/utils"
)

type Ground struct {
	Ally  Camp
	Enemy Camp
}

func NewGround(ally, enemy Camp) Ground {
	return Ground{
		Ally:  ally,
		Enemy: enemy,
	}
}

// GetMatchOrder returns the order of GroundIdx for each Match.
func (g Ground) GetMatchOrder() ([]GroundIdx, error) {
	// assume ally summoner is faster than enemy summoner, action idx is starting from 1 to 5
	startActionIdx := 1
	enemySummoner, err := g.GetEnemySummoner()
	if err != nil {
		return nil, fmt.Errorf("get enemy summoner error: %w", err)
	}
	allySummoner, err := g.GetAllySummoner()
	if err != nil {
		return nil, fmt.Errorf("get ally summoner error: %w", err)
	}
	if enemySummoner.GetAgi() > allySummoner.GetAgi() {
		// if enemy summoner is faster than ally summoner, action idx is starting from 7 to 11
		startActionIdx = 7
	}
	actionOrder := make([]GroundIdx, 0, 10)
	for i := startActionIdx; i < startActionIdx+5; i++ {
		firstActorIdx := GroundIdx(i)
		actionOrder = append(actionOrder, firstActorIdx, firstActorIdx.GetOppositeIdx())
	}
	return actionOrder, nil
}

func (g Ground) GetAllySummoner() (Character, error) {
	return g.Ally.GetSummoner()
}

func (g Ground) GetEnemySummoner() (Character, error) {
	return g.Enemy.GetSummoner()
}

func (g Ground) UpdateCharacter(idx GroundIdx, character Character) (Ground, error) {
	if idx.IsEnemy() {
		enemy, err := g.Enemy.Set(idx.ToCampIdx(), character)
		if err != nil {
			return Ground{}, fmt.Errorf("set enemy error: %w", err)
		}
		return Ground{
			Ally:  g.Ally,
			Enemy: enemy,
		}, nil
	}

	ally, err := g.Ally.Set(idx.ToCampIdx(), character)
	if err != nil {
		return Ground{}, fmt.Errorf("set ally error: %w", err)
	}
	return Ground{
		Ally:  ally,
		Enemy: g.Enemy,
	}, nil
}

type Match struct {
	Actor   Character
	Targets []Character
}

// GetMatch returns the match by attacker GroundIdx.
func (g Ground) GetMatch(idx GroundIdx) (match Match, err error) {
	getAttackerFn := g.Ally.Get
	getTargetFn := getTargetsFn(g.Enemy)
	if idx.IsEnemy() {
		getAttackerFn = g.Enemy.Get
		getTargetFn = getTargetsFn(g.Ally)
	}

	attacker, err := getAttackerFn(idx.ToCampIdx())
	if err != nil {
		return Match{}, fmt.Errorf("get attacker error: %w", err)
	}
	if attacker.IsDead() {
		return Match{
			Actor: attacker,
		}, nil
	}
	targetNumber := int(attacker.GetSkill().AttributeMap.Get(AttributeTypeTarget).Value.IntPart())
	targets, err := getTargetFn(targetNumber)
	if err != nil {
		return Match{}, fmt.Errorf("get targets error: %w", err)
	}
	return Match{
		Actor:   attacker,
		Targets: targets,
	}, nil
}

func (g Ground) ToString() string {
	return utils.ToString(g)
}

// getTargetsFn returns a function to get targets from camp.
func getTargetsFn(camp Camp) func(number int) (targets []Character, err error) {
	return func(number int) (targets []Character, err error) {
		targets = make([]Character, 0, number)
		randIdx := startFromFirst(1, 5, number)
		for _, idx := range randIdx {
			target, err := camp.Get(CampIdx(idx))
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

type RangeNumberPickerFn func(min, max, number int) []int

var _ RangeNumberPickerFn = startFromFirst

func startFromFirst(min, max, number int) []int {
	var res []int
	for i := min; i <= max && len(res) < number; i++ {
		res = append(res, i)
	}
	return res
}
