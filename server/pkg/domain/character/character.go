package character

import (
	"strconv"

	"github.com/Amobe/PlayGame/server/pkg/domain/skill"
	"github.com/Amobe/PlayGame/server/pkg/domain/valueobject"
	"github.com/Amobe/PlayGame/server/pkg/utils"
	"github.com/Amobe/PlayGame/server/pkg/utils/domain"
)

var _ domain.Aggregator = &Character{}

type coreAggregator = domain.CoreAggregator

type Character struct {
	coreAggregator
	CharacterID string
	Basement    valueobject.AttributeTypeMap
	Equipment   Equipment
}

func NewCharacter(id string, attrs ...valueobject.Attribute) Character {
	c := Character{
		CharacterID: id,
		Basement:    valueobject.NewAttributeTypeMap(),
	}
	c.Basement.Insert(attrs...)
	return c
}

func RandomCharacter(id string) Character {
	hp := utils.GetRandIntInRange(100, 200)
	atk := utils.GetRandIntInRange(20, 50)
	return NewCharacter(id,
		valueobject.Attribute{Type: valueobject.AttributeTypeHP, Value: strconv.Itoa(hp)},
		valueobject.Attribute{Type: valueobject.AttributeTypeATK, Value: strconv.Itoa(atk)},
	)
}

func (c Character) ID() string {
	return c.CharacterID
}

func (c Character) AttributeMap() valueobject.AttributeTypeMap {
	res := valueobject.NewAttributeTypeMap()
	for _, attr := range c.Basement {
		res.Insert(attr)
	}
	for _, attr := range c.Equipment.GetAttributes() {
		res.Insert(attr)
	}
	return res
}

func (c Character) GetAgi() int {
	return c.AttributeMap().Get(valueobject.AttributeTypeAGI)
}

func (c Character) Alive() bool {
	attrMap := c.AttributeMap()
	_, ok := attrMap[valueobject.AttributeTypeDead]
	return !ok
}

func (c Character) Affect(attr []valueobject.Attribute) {
	c.Basement.Insert(attr...)

	// dead
	am := c.AttributeMap()
	if am[valueobject.AttributeTypeHP].GetInt() == 0 {
		c.Basement.Insert(valueobject.Attribute{Type: valueobject.AttributeTypeDead})
	}
}

func (c Character) UseSkill(skill skill.Skill, targetAttr valueobject.AttributeTypeMap) (targetAffect []valueobject.Attribute) {
	affect, targetAffect := skill.Use(c.AttributeMap(), targetAttr)
	c.Affect(affect)
	return targetAffect
}

type CharacterEventDefensed struct {
}
