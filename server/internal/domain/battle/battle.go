package battle

import (
	"fmt"

	"github.com/Amobe/PlayGame/server/internal/domain/vo"
	"github.com/Amobe/PlayGame/server/internal/utils/domain"
)

var _ domain.Aggregator = &Battle{}

type coreAggregator = domain.CoreAggregator

type Battle struct {
	coreAggregator
	battleID   string
	status     Status
	minionSlot *MinionSlot
}

func (b Battle) ID() string {
	return b.battleID
}

func (b Battle) Status() Status {
	return b.status
}

func (b Battle) MinionSlot() *MinionSlot {
	return b.minionSlot
}

func newBattle() *Battle {
	return &Battle{}
}

func AggregatorLoader(events []domain.Event) (*Battle, error) {
	b := newBattle()
	if err := b.apply(false, events...); err != nil {
		return nil, fmt.Errorf("apply battle events: %w", err)
	}
	return b, nil
}

func CreateBattle(id string, minionSlot *MinionSlot) (*Battle, error) {
	createdEvent := EventBattleCreated{
		BattleID:   id,
		MinionSlot: minionSlot,
	}
	b := newBattle()
	if err := b.applyNew(createdEvent); err != nil {
		return nil, fmt.Errorf("apply battle created event: %w", err)
	}
	return b, nil
}

func (b *Battle) FightToTheEnd() ([]vo.Affect, error) {
	if b.status != StatusUnspecified {
		return nil, fmt.Errorf("battle is finished")
	}

	var battleAffects []vo.Affect
	const roundLimit = 50
	// Fight until ally or enemy dead.
	for i := 0; i < roundLimit; i++ {
		if b.status != StatusUnspecified {
			break
		}
		affects, err := b.fight()
		if err != nil {
			return nil, fmt.Errorf("fight: %w", err)
		}
		battleAffects = append(battleAffects, affects...)
	}
	if b.status == StatusUnspecified {
		if err := b.applyNew(EventBattleDraw{}); err != nil {
			return nil, fmt.Errorf("apply battle draw event: %w", err)
		}
	}
	return battleAffects, nil
}

func (b *Battle) fight() ([]vo.Affect, error) {
	affects, err := b.minionSlot.PlayOneRound()
	if err != nil {
		return nil, fmt.Errorf("play one round: %w", err)
	}
	if err := b.applyNew(EventBattleFought{Affects: affects}); err != nil {
		return nil, fmt.Errorf("apply battle fought event: %w", err)
	}
	if b.minionSlot.Status == MinionSlotStatusAllyWon {
		if err := b.applyNew(EventBattleWon{}); err != nil {
			return nil, fmt.Errorf("apply battle won event: %w", err)
		}
	}
	if b.minionSlot.Status == MinionSlotStatusEnemyWon {
		if err := b.applyNew(EventBattleLost{}); err != nil {
			return nil, fmt.Errorf("apply battle lost event: %w", err)
		}
	}
	return affects, nil
}

func (b *Battle) applyNew(events ...domain.Event) error {
	return b.apply(true, events...)
}

func (b *Battle) apply(new bool, events ...domain.Event) error {
	for _, event := range events {
		switch ev := event.(type) {
		case EventBattleCreated:
			b.battleID = ev.BattleID
			b.status = StatusUnspecified
			b.minionSlot = ev.MinionSlot
		case EventBattleFought:
			// TODO: the model changed before apply event, so the event is not valid.
			// Take the affects from event and apply to model.
			continue
		case EventBattleWon:
			b.status = StatusWon
		case EventBattleLost:
			b.status = StatusLost
		case EventBattleDraw:
			b.status = StatusDraw
		default:
			return fmt.Errorf("unspecified event type: %v", ev)
		}
	}
	b.coreAggregator.Apply(new, events...)
	return nil
}

type Status string

const (
	StatusUnspecified Status = "unspecified"
	StatusWon         Status = "won"
	StatusLost        Status = "lost"
	StatusDraw        Status = "draw"
)
