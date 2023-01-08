package character

import (
	"strconv"

	"github.com/Amobe/PlayGame/server/pkg/utils"
)

type Character struct {
	CharacterID string
	Basement    AttributeTypeMap
	Equipment   Equipment
}

func NewCharacter() Character {
	c := Character{
		Basement: NewAttributeTypeMap(),
	}
	return c
}

func RandomCharacter() Character {
	c := NewCharacter()
	hp := utils.GetRandIntInRange(100, 200)
	atk := utils.GetRandIntInRange(20, 50)
	c.Basement.Insert(Attribute{Type: AttributeTypeHP, Value: strconv.Itoa(hp)})
	c.Basement.Insert(Attribute{Type: AttributeTypeATK, Value: strconv.Itoa(atk)})
	return c
}

func (c Character) ID() string {
	return c.CharacterID
}

func (c *Character) GetAttributeMap() AttributeTypeMap {
	res := NewAttributeTypeMap()
	for _, attr := range c.Basement {
		res.Insert(attr)
	}
	for _, attr := range c.Equipment.GetAttributes() {
		res.Insert(attr)
	}
	return res
}

func (c *Character) GetAgi() int {
	attrMap := c.GetAttributeMap()
	attr, ok := attrMap[AttributeTypeAGI]
	if !ok {
		return 0
	}
	return attr.GetInt()
}

func (c *Character) Alive() bool {
	attrMap := c.GetAttributeMap()
	_, ok := attrMap[AttributeTypeDead]
	return !ok
}

func (c *Character) Affect(attr []Attribute) {
	c.Basement.Insert(attr...)

	// dead
	am := c.GetAttributeMap()
	if am[AttributeTypeHP].GetInt() == 0 {
		c.Basement.Insert(Attribute{Type: AttributeTypeDead})
	}
}

func (c *Character) UseSkill(skill Skill, targetAttr AttributeTypeMap) (targetAffect []Attribute) {
	affect, targetAffect := skill.Use(c.GetAttributeMap(), targetAttr)
	c.Affect(affect)
	return targetAffect
}

type CharacterEventDefensed struct {
}
