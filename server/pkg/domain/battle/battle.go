package battle

import (
	"fmt"

	"github.com/Amobe/PlayGame/server/pkg/domain/vo"
	"github.com/Amobe/PlayGame/server/pkg/utils"
	"github.com/Amobe/PlayGame/server/pkg/utils/domain"
)

var _ domain.Aggregator = Battle{}

type coreAggregator = domain.CoreAggregator

type Battle struct {
	coreAggregator
	battleID   string
	status     Status
	fighterMap map[string]Fighter
	targetMap  map[string]string // bi-direction map between ally and enemy
	allyMap    map[string]interface{}
	enemyMap   map[string]interface{}
	allySlot   Slot
	enemySlot  Slot
	order      []string
}

func (b Battle) ID() string {
	return b.battleID
}

func (b Battle) Status() Status {
	return b.status
}

func (b Battle) Fighter(id string) Fighter {
	return b.fighterMap[id]
}

func newBattle() Battle {
	return Battle{
		fighterMap: make(map[string]Fighter),
		allyMap:    make(map[string]interface{}),
		enemyMap:   make(map[string]interface{}),
		targetMap:  make(map[string]string),
	}
}

func AggregatorLoader(events []domain.Event) (Battle, error) {
	b := newBattle()
	if err := b.apply(false, events...); err != nil {
		return Battle{}, fmt.Errorf("apply battle events: %w", err)
	}
	return b, nil
}

func CreateBattle(id string, ally Fighter, enemy Fighter, enemySlot Slot) (Battle, error) {
	fighterMap := map[string]Fighter{
		ally.ID():  ally,
		enemy.ID(): enemy,
	}
	targetMap := map[string]string{
		ally.ID():  enemy.ID(),
		enemy.ID(): ally.ID(),
	}
	allyMap := map[string]interface{}{
		ally.ID(): struct{}{},
	}
	enemyMap := map[string]interface{}{
		enemy.ID(): struct{}{},
	}
	fighters := []Fighter{
		ally,
		enemy,
	}
	createdEvent := EventBattleCreated{
		BattleID:   id,
		FighterMap: fighterMap,
		TargetMap:  targetMap,
		AllyMap:    allyMap,
		EnemyMap:   enemyMap,
		EnemySlot:  enemySlot,
		Order:      getFighterOrder(fighters...),
	}
	b := newBattle()
	if err := b.applyNew(createdEvent); err != nil {
		return Battle{}, fmt.Errorf("apply battle created event: %w", err)
	}
	return b, nil
}

func (b *Battle) SetAllySlot(slot Slot) error {
	allySlotSetEvent := EventBattleAllySlotSet{
		AllySlot: slot,
	}
	if err := b.applyNew(allySlotSetEvent); err != nil {
		return fmt.Errorf("apply battle ally slot set event: %w", err)
	}
	return nil
}

func (b *Battle) FightToTheEnd() error {
	const roundLimit = 50
	// Fight until ally or enemy dead.
	for i := 0; i < roundLimit; i++ {
		if err := b.Fight(); err != nil {
			return fmt.Errorf("fight: %w", err)
		}
	}
	if err := b.applyNew(EventBattleDraw{}); err != nil {
		return fmt.Errorf("apply battle draw event: %w", err)
	}
	return nil
}

func (b *Battle) Fight() error {
	for _, actorID := range b.order {
		actor := b.fighterMap[actorID]
		target, err := b.getTarget(actorID)
		if err != nil {
			return fmt.Errorf("get target: %w", err)
		}
		slot, err := b.getSlot(actorID)
		if err != nil {
			return fmt.Errorf("get slot: %w", err)
		}
		var affects []Affect
		for _, s := range slot.Skills {
			if s.IsEmpty() {
				continue
			}
			aa, ta := b.useSkill(s, actor, target)
			if len(aa.Attributes) > 0 {
				affects = append(affects, aa)
			}
			if len(ta.Attributes) > 0 {
				affects = append(affects, ta)
			}
		}
		if len(affects) > 0 {
			if err := b.applyNew(EventBattleFought{Affects: affects}); err != nil {
				return fmt.Errorf("apply battle fought event: %w", err)
			}
		}
		if b.isAllDead(b.enemyMap) {
			if err := b.applyNew(EventBattleWon{}); err != nil {
				return fmt.Errorf("apply battle won event: %w", err)
			}
			return nil
		}
		if b.isAllDead(b.allyMap) {
			if err := b.applyNew(EventBattleLost{}); err != nil {
				return fmt.Errorf("apply battle lost event: %w", err)
			}
			return nil
		}
	}
	return nil
}

func (b *Battle) getSlot(actorID string) (slot Slot, err error) {
	if _, ok := b.allyMap[actorID]; ok {
		slot = b.allySlot
		return
	}
	if _, ok := b.enemyMap[actorID]; ok {
		slot = b.enemySlot
		return
	}
	err = fmt.Errorf("actor id is missing in ally and enemy map")
	return
}

func (b *Battle) getTarget(actorID string) (target Fighter, err error) {
	targetID, ok := b.targetMap[actorID]
	if !ok {
		err = fmt.Errorf("actor id is missing in target map")
		return
	}
	target, ok = b.fighterMap[targetID]
	if !ok {
		err = fmt.Errorf("target id is missing in fighter map")
		return
	}
	return
}

func (b *Battle) isAllDead(ids map[string]interface{}) bool {
	for id := range ids {
		if b.fighterMap[id].Alive() {
			return false
		}
	}
	return true
}

func (b *Battle) useSkill(skill vo.Skill, actor, target Fighter) (actorAffect, targetAffect Affect) {
	aa, ta := skill.Use(actor.AttributeMap(), target.AttributeMap())
	skillName := skill.SkillType.String()
	actorAffect = NewAffect(actor.ID(), target.ID(), actor.ID(), skillName, aa)
	targetAffect = NewAffect(actor.ID(), target.ID(), target.ID(), skillName, ta)
	return actorAffect, targetAffect
}

func (b *Battle) applyNew(events ...domain.Event) error {
	return b.apply(true, events...)
}

func (b *Battle) apply(new bool, events ...domain.Event) error {
	for _, event := range events {
		switch ev := event.(type) {
		case EventBattleCreated:
			b.battleID = ev.BattleID
			for id, c := range ev.FighterMap {
				b.fighterMap[id] = c
			}
			b.fighterMap = map[string]Fighter{}
			utils.CopyMap(b.fighterMap, ev.FighterMap)
			b.targetMap = map[string]string{}
			utils.CopyMap(b.targetMap, ev.TargetMap)
			b.allyMap = make(map[string]interface{})
			utils.CopyMap(b.allyMap, ev.AllyMap)
			b.enemyMap = make(map[string]interface{})
			utils.CopyMap(b.enemyMap, ev.EnemyMap)
			b.enemySlot = ev.EnemySlot
			order := make([]string, len(ev.Order))
			copy(order, ev.Order)
			b.order = order
			b.status = StatusUnspecified
		case EventBattleAllySlotSet:
			b.allySlot = ev.AllySlot
		case EventBattleFought:
			for _, a := range ev.Affects {
				b.fighterMap[a.ChangerID] = b.fighterMap[a.ChangerID].Affect(a.Attributes)
			}
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
