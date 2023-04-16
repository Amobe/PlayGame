package vo

import (
	"github.com/shopspring/decimal"

	"github.com/Amobe/PlayGame/server/internal/utils"
)

type Character struct {
	CharacterID string
	Basement    AttributeMap
	Equipment   Equipment
	Skill       Skill
}

func NewCharacter(id string, attrs AttributeMap) Character {
	return NewCharacterWithSkill(id, EmptySkill, attrs)
}

func NewCharacterWithSkill(id string, skill Skill, attrs AttributeMap) Character {
	c := Character{
		CharacterID: id,
		Basement:    attrs,
		Skill:       skill,
	}
	return c
}

func RandomCharacter(id string) Character {
	hp := utils.GetRandIntInRange(100, 200)
	atk := utils.GetRandIntInRange(20, 50)
	return NewCharacter(id,
		NewAttributeMap(
			NewAttribute(AttributeTypeHP, decimal.NewFromInt(int64(hp))),
			NewAttribute(AttributeTypeATK, decimal.NewFromInt(int64(atk))),
		),
	)
}

func (c Character) ID() string {
	return c.CharacterID
}

func (c Character) GetAttributeMap() AttributeMap {
	res := NewAttributeMap()
	for _, attr := range c.Basement {
		res.Insert(attr)
	}
	for _, attr := range c.Equipment.GetAttributes() {
		res.Insert(attr)
	}
	return res
}

func (c Character) GetAgi() int {
	return int(c.GetAttributeMap().Get(AttributeTypeAGI).Value.InexactFloat64())
}

func (c Character) GetSkill() Skill {
	return c.Skill
}

func (c Character) IsDead() bool {
	attrMap := c.GetAttributeMap()
	_, ok := attrMap[AttributeTypeDead]
	return ok
}

func (c Character) TakeAffect(attrs AttributeMap) Character {
	c.Basement = c.Basement.Merge(attrs)

	// dead
	am := c.GetAttributeMap()
	if am.Get(AttributeTypeHP).Value.LessThanOrEqual(decimal.Zero) {
		c.Basement = c.Basement.Insert(NewAttribute(AttributeTypeDead, decimal.NewFromInt(1)))
	}

	return Character{
		CharacterID: c.CharacterID,
		Basement:    c.Basement,
		Equipment:   c.Equipment,
	}
}
