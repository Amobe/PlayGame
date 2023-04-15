package vo

import "github.com/shopspring/decimal"

var (
	SkillEmpty = Skill{
		SkillType: SkillTypeEmpty,
	}
	SkillSlash = NewSkill("slash", NewAttribute(AttributeTypeTarget, decimal.NewFromInt(1)))
)

type Skill struct {
	SkillType    SkillType
	Name         string
	AttributeMap AttributeMap
}

func NewSkill(name string, attrs ...Attribute) Skill {
	return Skill{
		Name:         name,
		SkillType:    SkillTypeNormal,
		AttributeMap: NewAttributeMap(attrs...),
	}
}

func (s Skill) IsEmpty() bool {
	return s.SkillType == SkillTypeEmpty
}

// EqualTo compares two Skill objects.
func (s Skill) EqualTo(other Skill) bool {
	if s.SkillType != other.SkillType {
		return false
	}
	if s.Name != other.Name {
		return false
	}
	if !s.AttributeMap.EqualTo(other.AttributeMap) {
		return false
	}
	return true
}

type SkillType string

func (s SkillType) String() string {
	return string(s)
}

const (
	SkillTypeNormal SkillType = "normal"
	SkillTypeEmpty  SkillType = "empty"
)
