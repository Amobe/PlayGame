package vo

import "github.com/shopspring/decimal"

var (
	SkillEmpty = Skill{
		SkillType: SkillTypeEmpty,
	}
	SkillSlash = NewSkill("slash", "slash", NewAttribute(AttributeTypeTarget, decimal.NewFromInt(1)))
)

type Skill struct {
	SkillID      string
	SkillType    SkillType
	Name         string
	AttributeMap AttributeMap
}

func NewSkill(skillID string, name string, attrs ...Attribute) Skill {
	return Skill{
		SkillID:      skillID,
		Name:         name,
		SkillType:    SkillTypeNormal,
		AttributeMap: NewAttributeMap(attrs...),
	}
}

func (s Skill) ID() string {
	return s.SkillID
}

func (s Skill) IsEmpty() bool {
	return s.SkillType == SkillTypeEmpty
}

type SkillType string

func (s SkillType) String() string {
	return string(s)
}

const (
	SkillTypeNormal SkillType = "normal"
	SkillTypeEmpty  SkillType = "empty"
)
