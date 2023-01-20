package character

import (
	"strconv"

	"github.com/Amobe/PlayGame/server/pkg/domain/skill"
	"github.com/Amobe/PlayGame/server/pkg/utils"
	"github.com/Amobe/PlayGame/server/pkg/utils/domain"
)

var _ domain.Aggregator = &Character{}

type coreAggregator = domain.CoreAggregator

type Character struct {
	coreAggregator
	CharacterID string
	Basement    AttributeTypeMap
	Equipment   Equipment
}

func NewCharacter(id string, attrs ...Attribute) Character {
	c := Character{
		CharacterID: id,
		Basement:    NewAttributeTypeMap(),
	}
	c.Basement.Insert(attrs...)
	return c
}

func RandomCharacter(id string) Character {
	hp := utils.GetRandIntInRange(100, 200)
	atk := utils.GetRandIntInRange(20, 50)
	return NewCharacter(id,
		Attribute{Type: AttributeTypeHP, Value: strconv.Itoa(hp)},
		Attribute{Type: AttributeTypeATK, Value: strconv.Itoa(atk)},
	)
}

func (c Character) ID() string {
	return c.CharacterID
}

func (c Character) AttributeMap() AttributeTypeMap {
	res := NewAttributeTypeMap()
	for _, attr := range c.Basement {
		res.Insert(attr)
	}
	for _, attr := range c.Equipment.GetAttributes() {
		res.Insert(attr)
	}
	return res
}

func (c Character) GetAgi() int {
	return c.AttributeMap().Get(AttributeTypeAGI)
}

func (c Character) Alive() bool {
	attrMap := c.AttributeMap()
	_, ok := attrMap[AttributeTypeDead]
	return !ok
}

func (c Character) Affect(attr []Attribute) {
	c.Basement.Insert(attr...)

	// dead
	am := c.AttributeMap()
	if am[AttributeTypeHP].GetInt() == 0 {
		c.Basement.Insert(Attribute{Type: AttributeTypeDead})
	}
}

func (c Character) UseSkill(skill skill.Skill, targetAttr AttributeTypeMap) (targetAffect []Attribute) {
	affect, targetAffect := skill.Use(c.AttributeMap(), targetAttr)
	c.Affect(affect)
	return targetAffect
}

type CharacterEventDefensed struct {
}
