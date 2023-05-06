package vo

import (
	"github.com/shopspring/decimal"

	"github.com/Amobe/PlayGame/server/internal/utils"
)

var EmptyCharacter = Character{}

type Character struct {
	GroundIdx GroundIdx
	Basement  AttributeMap
	Equipment Equipment
	Skill     Skill
}

func NewCharacter(groundIdx GroundIdx, attrs AttributeMap) Character {
	return NewCharacterWithSkill(groundIdx, EmptySkill, attrs)
}

func NewCharacterWithSkill(groundIdx GroundIdx, skill Skill, attrs AttributeMap) Character {
	c := Character{
		GroundIdx: groundIdx,
		Basement:  attrs,
		Skill:     skill,
	}
	return c
}

func RandomCharacter(groundIdx GroundIdx) Character {
	hp := utils.GetRandIntInRange(100, 200)
	atk := utils.GetRandIntInRange(20, 50)
	return NewCharacter(groundIdx,
		NewAttributeMap(
			NewAttribute(AttributeTypeHP, decimal.NewFromInt(int64(hp))),
			NewAttribute(AttributeTypeATK, decimal.NewFromInt(int64(atk))),
		),
	)
}

func (c Character) GetGroundIdx() GroundIdx {
	return c.GroundIdx
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
		GroundIdx: c.GroundIdx,
		Basement:  c.Basement,
		Equipment: c.Equipment,
	}
}
