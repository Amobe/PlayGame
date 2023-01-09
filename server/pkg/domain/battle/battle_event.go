package battle

import "github.com/Amobe/PlayGame/server/pkg/utils/domain"

type coreEvent = domain.CoreEvent

type EventBattleCreated struct {
	coreEvent
	BattleID   string
	FighterMap map[string]Fighter
	AllyMap    map[string]interface{}
	EnemyMap   map[string]interface{}
	Order      []string // the order of the action, contains a list of character ID
}

func (EventBattleCreated) Name() string {
	return "battle_created"
}

type EventBattleFought struct {
	coreEvent
	Affects []Affect
}

func (EventBattleFought) Name() string {
	return "battle_fought"
}

type EventBattleWon struct {
	coreEvent
}

func (EventBattleWon) Name() string {
	return "battle_won"
}

type EventBattleLost struct {
	coreEvent
}

func (EventBattleLost) Name() string {
	return "battle_lost"
}
