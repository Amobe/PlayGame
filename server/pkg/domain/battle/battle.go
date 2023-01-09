package battle

import (
	"fmt"

	"github.com/Amobe/PlayGame/server/pkg/domain/character"
	"github.com/Amobe/PlayGame/server/pkg/utils"
	"github.com/Amobe/PlayGame/server/pkg/utils/domain"
)

var _ domain.Aggregator = &Battle{}

type coreAggregator = domain.CoreAggregator

type Battle struct {
	coreAggregator
	battleID   string
	status     Status
	fighterMap map[string]Fighter
	allyMap    map[string]interface{}
	enemyMap   map[string]interface{}
	order      []string
}

func (b *Battle) ID() string {
	return b.battleID
}

func newBattle() *Battle {
	return &Battle{
		fighterMap: make(map[string]Fighter),
		allyMap:    make(map[string]interface{}),
		enemyMap:   make(map[string]interface{}),
	}
}

func AggregatorLoader(events []domain.Event) (*Battle, error) {
	b := newBattle()
	if err := b.apply(false, events...); err != nil {
		return nil, fmt.Errorf("apply battle events: %w", err)
	}
	return b, nil
}

func CreateBattle(id string, ally Fighter, enemies []Fighter) (*Battle, error) {
	fighterMap := map[string]Fighter{
		ally.ID(): ally,
	}
	allyMap := map[string]interface{}{
		ally.ID(): struct{}{},
	}
	enemyMap := make(map[string]interface{})
	order := []string{ally.ID()}

	for _, enemy := range enemies {
		fighterMap[enemy.ID()] = enemy
		enemyMap[enemy.ID()] = struct{}{}
		order = append(order, enemy.ID())
	}

	createdEvent := EventBattleCreated{
		BattleID:   id,
		FighterMap: fighterMap,
		AllyMap:    allyMap,
		EnemyMap:   enemyMap,
		Order:      order,
	}
	b := newBattle()
	if err := b.applyNew(createdEvent); err != nil {
		return nil, fmt.Errorf("apply battle created event: %w", err)
	}
	return b, nil
}

func (b *Battle) Fight(skills []character.Skill) error {
	var affects []Affect
	for _, id := range b.order {
		actor := b.fighterMap[id]
		for _, s := range skills {
			targets := b.getAliveCharacter(1, b.enemyMap)
			for _, target := range targets {
				attrs := actor.UseSkill(s, target.AttributeMap())
				affect := Affect{
					ActorID:    actor.ID(),
					TargetID:   target.ID(),
					Skill:      s.Name(),
					Attributes: attrs,
				}
				affects = append(affects, affect)
			}
		}
	}
	if err := b.applyNew(EventBattleFought{Affects: affects}); err != nil {
		return fmt.Errorf("apply battle fought event: %w", err)
	}
	if b.isAllDead(b.enemyMap) {
		if err := b.applyNew(EventBattleWon{}); err != nil {
			return fmt.Errorf("apply battle won event: %w", err)
		}
	}
	if b.isAllDead(b.allyMap) {
		if err := b.applyNew(EventBattleLost{}); err != nil {
			return fmt.Errorf("apply battle lost event: %w", err)
		}
	}
	return nil
}

func (b *Battle) getAliveCharacter(n int, ids map[string]interface{}) []Fighter {
	var res []Fighter
	for id := range ids {
		if n == 0 {
			break
		}
		c := b.fighterMap[id]
		if !c.Alive() {
			continue
		}
		res = append(res, c)
		n--
	}
	return res
}

func (b *Battle) isAllDead(ids map[string]interface{}) bool {
	for id := range ids {
		if b.fighterMap[id].Alive() {
			return false
		}
	}
	return true
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
			allyMap := make(map[string]interface{})
			for id, v := range ev.AllyMap {
				allyMap[id] = v
			}
			b.allyMap = allyMap
			enemyMap := make(map[string]interface{})
			for id, v := range ev.EnemyMap {
				enemyMap[id] = v
			}
			b.enemyMap = enemyMap
			order := make([]string, len(ev.Order))
			copy(order, ev.Order)
			b.order = order
			b.status = StatusUnspecified
		case EventBattleFought:
			for _, a := range ev.Affects {
				b.fighterMap[a.TargetID].Affect(a.Attributes)
			}
		case EventBattleWon:
			b.status = StatusWon
		case EventBattleLost:
			b.status = StatusLost
		default:
			return fmt.Errorf("unspecified event type: %v", ev)
		}
	}
	if new {
		b.coreAggregator.Apply(events...)
	}
	return nil
}

type agiGetter interface {
	GetAgi() int
}

type actionOrder = utils.LinkedList[agiGetter]

var actionCond = func(current, next agiGetter) bool {
	return current.GetAgi() < next.GetAgi()
}

func newActionOrder(characters ...agiGetter) *actionOrder {
	list := utils.NewLinkedList(actionCond)
	for _, c := range characters {
		list.Insert(c)
	}
	return list
}

type Status string

const (
	StatusUnspecified Status = "unspecified"
	StatusWon         Status = "won"
	StatusLost        Status = "lost"
)
