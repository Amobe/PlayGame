package vo

import (
	"github.com/shopspring/decimal"

	"github.com/Amobe/PlayGame/server/pkg/utils"
)

type Character struct {
	CharacterID string
	Basement    AttributeMap
	Equipment   Equipment
}

func NewCharacter(id string, attrs ...Attribute) Character {
	c := Character{
		CharacterID: id,
		Basement:    NewAttributeMap(),
	}
	c.Basement.Insert(attrs...)
	return c
}

func RandomCharacter(id string) Character {
	hp := utils.GetRandIntInRange(100, 200)
	atk := utils.GetRandIntInRange(20, 50)
	return NewCharacter(id,
		Attribute{Type: AttributeTypeHP, Value: decimal.NewFromInt(int64(hp))},
		Attribute{Type: AttributeTypeATK, Value: decimal.NewFromInt(int64(atk))},
	)
}

func (c Character) ID() string {
	return c.CharacterID
}

func (c Character) AttributeMap() AttributeMap {
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
	return int(c.AttributeMap().Get(AttributeTypeAGI).Value.InexactFloat64())
}

func (c Character) Alive() bool {
	attrMap := c.AttributeMap()
	_, ok := attrMap[AttributeTypeDead]
	return !ok
}

func (c Character) Affect(attr []Attribute) Character {
	c.Basement = c.Basement.Insert(attr...)

	// dead
	am := c.AttributeMap()
	if am.Get(AttributeTypeHP).Value.LessThanOrEqual(decimal.Zero) {
		c.Basement = c.Basement.Insert(NewAttribute(AttributeTypeDead, decimal.NewFromInt(1)))
	}

	return Character{
		CharacterID: c.CharacterID,
		Basement:    c.Basement,
		Equipment:   c.Equipment,
	}
}
