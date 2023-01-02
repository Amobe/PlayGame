package battle

import (
	"fmt"
	"github.com/Amobe/PlayGame/server/pkg/domain/character"
	"github.com/Amobe/PlayGame/server/pkg/utils"
)

type Battle struct {
	ally  *character.Character
	enemy *Mob
	order *actionOrder
}

func NewBattle(ally character.Character, enemy Mob) Battle {
	order := newActionOrder(&ally, &enemy)
	return Battle{
		ally:  &ally,
		enemy: &enemy,
		order: order,
	}
}

func (b *Battle) Fight(skills []character.Skill) error {
	iter := b.order.Iterator()
	for iter.HasNext() {
		next, err := iter.Next()
		if err != nil {
			return fmt.Errorf("get next actor: %w", err)
		}
		switch actor := next.(type) {
		case *character.Character:
			for _, s := range skills {
				target := b.enemy
				ta := actor.UseSkill(s, target.GetAttributeMap())
				target.Affect(ta)
			}
		case *Mob:
			target := b.ally
			ta := actor.UseSkill(target.GetAttributeMap())
			target.Affect(ta)
		default:
			return fmt.Errorf("unknown actor type to fight")
		}
	}
	return nil
}

type agiGetter interface {
	GetAgi() int
}

type actionOrder = utils.LinkedList[agiGetter]

var actionCond = func(current, next agiGetter) bool {
	return current.GetAgi() > next.GetAgi()
}

func newActionOrder(characters ...agiGetter) *actionOrder {
	list := utils.NewLinkedList(actionCond)
	for _, c := range characters {
		list.Insert(c)
	}
	return list
}
