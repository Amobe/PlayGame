package character

type Character struct {
	CharacterID string
	Basement    []Attribute
	Equipment   Equipment
}

func NewCharacter() Character {
	return Character{
		Basement: []Attribute{
			{Type: AttributeTypeHP, Value: "100"},
		},
	}
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
	c.Basement = append(c.Basement, attr...)

	// dead
	am := c.GetAttributeMap()
	if am[AttributeTypeHP].GetInt() == 0 {
		c.Basement = append(c.Basement, Attribute{Type: AttributeTypeDead})
	}
}

func (c *Character) UseSkill(skill Skill, targetAttr AttributeTypeMap) (targetAffect []Attribute) {
	affect, targetAffect := skill.Use(c.GetAttributeMap(), targetAttr)
	c.Affect(affect)
	return targetAffect
}

type CharacterEventDefensed struct {
}
