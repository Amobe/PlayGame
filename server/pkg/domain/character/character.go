package character

import (
	"github.com/shopspring/decimal"

	"github.com/Amobe/PlayGame/server/pkg/domain/vo"
	"github.com/Amobe/PlayGame/server/pkg/utils"
	"github.com/Amobe/PlayGame/server/pkg/utils/domain"
)

var _ domain.Aggregator = &Character{}

type coreAggregator = domain.CoreAggregator

type Character struct {
	coreAggregator
	CharacterID string
	Basement    vo.AttributeMap
	Equipment   Equipment
}

func NewCharacter(id string, attrs ...vo.Attribute) Character {
	c := Character{
		CharacterID: id,
		Basement:    vo.NewAttributeTypeMap(),
	}
	c.Basement.Insert(attrs...)
	return c
}

func RandomCharacter(id string) Character {
	hp := utils.GetRandIntInRange(100, 200)
	atk := utils.GetRandIntInRange(20, 50)
	return NewCharacter(id,
		vo.Attribute{Type: vo.AttributeTypeHP, Value: decimal.NewFromInt(int64(hp))},
		vo.Attribute{Type: vo.AttributeTypeATK, Value: decimal.NewFromInt(int64(atk))},
	)
}

func (c Character) ID() string {
	return c.CharacterID
}

func (c Character) AttributeMap() vo.AttributeMap {
	res := vo.NewAttributeTypeMap()
	for _, attr := range c.Basement {
		res.Insert(attr)
	}
	for _, attr := range c.Equipment.GetAttributes() {
		res.Insert(attr)
	}
	return res
}

func (c Character) GetAgi() int {
	return int(c.AttributeMap().Get(vo.AttributeTypeAGI).Value.InexactFloat64())
}

func (c Character) Alive() bool {
	attrMap := c.AttributeMap()
	_, ok := attrMap[vo.AttributeTypeDead]
	return !ok
}

func (c Character) Affect(attr []vo.Attribute) {
	c.Basement.Insert(attr...)

	// dead
	am := c.AttributeMap()
	if am[vo.AttributeTypeHP].Value.IsZero() {
		c.Basement.Insert(vo.NewAttribute(vo.AttributeTypeDead, decimal.NewFromInt(1)))
	}
}

func (c Character) UseSkill(skill Skill, targetAttr vo.AttributeMap) (targetAffect []vo.Attribute) {
	affect, targetAffect := skill.Use(c.AttributeMap(), targetAttr)
	c.Affect(affect)
	return targetAffect
}

type CharacterEventDefensed struct {
}
