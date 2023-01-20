package battle

import (
	"github.com/Amobe/PlayGame/server/pkg/domain/character"
	"github.com/Amobe/PlayGame/server/pkg/domain/skill"
	"github.com/Amobe/PlayGame/server/pkg/utils"
)

type Fighter interface {
	Affect(attr []character.Attribute)
	UseSkill(skill skill.Skill, targetAttr character.AttributeTypeMap) (targetAffect []character.Attribute)

	ID() string
	Alive() bool
	AttributeMap() character.AttributeTypeMap
	GetAgi() int
}

func getFighterOrder(fighters ...Fighter) []string {
	cond := func(current, new Fighter) bool {
		return current.GetAgi() > new.GetAgi()
	}
	ll := utils.NewLinkedList(cond)
	for _, f := range fighters {
		ll.Insert(f)
	}
	var res []string
	iter := ll.Iterator()
	for iter.HasNext() {
		v, _ := iter.Next()
		res = append(res, v.ID())
	}
	return res
}
