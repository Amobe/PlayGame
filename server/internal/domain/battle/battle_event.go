package battle

import "github.com/Amobe/PlayGame/server/internal/utils/domain"

type EventBattleCreated struct {
	domain.CoreEvent
	BattleID   string
	MinionSlot *MinionSlot
}

func (EventBattleCreated) Name() string {
	return "battle_created"
}

type EventBattleWon struct {
	domain.CoreEvent
}

func (EventBattleWon) Name() string {
	return "battle_won"
}

type EventBattleLost struct {
	domain.CoreEvent
}

func (EventBattleLost) Name() string {
	return "battle_lost"
}

type EventBattleDraw struct {
	domain.CoreEvent
}

func (EventBattleDraw) Name() string {
	return "battle_draw"
}
