package battle

import "github.com/Amobe/PlayGame/server/internal/utils/domain"

type coreEvent = domain.CoreEvent

type EventBattleCreated struct {
	coreEvent
	BattleID   string
	MinionSlot *MinionSlot
}

func (EventBattleCreated) Name() string {
	return "battle_created"
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
