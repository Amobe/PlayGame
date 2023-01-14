package battle

import "github.com/Amobe/PlayGame/server/pkg/utils/domain"

type coreEvent = domain.CoreEvent

type EventBattleCreated struct {
	coreEvent
	BattleID   string
	FighterMap map[string]Fighter
	TargetMap  map[string]string
	AllyMap    map[string]interface{}
	EnemyMap   map[string]interface{}
	EnemySlot  Slot
	Order      []string // the order of the action, contains a list of character ID
}

func (EventBattleCreated) Name() string {
	return "battle_created"
}

type EventBattleAllySlotSet struct {
	coreEvent
	AllySlot Slot
}

func (EventBattleAllySlotSet) Name() string {
	return "battle_ally_slot_Set"
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

type EventBattleDraw struct {
	coreEvent
}

func (EventBattleDraw) Name() string {
	return "battle_draw"
}
