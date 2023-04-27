package vo

import "github.com/shopspring/decimal"

var (
	EmptySkill = Skill{
		SkillType: SkillTypeEmpty,
	}
	SkillSlash = newSlashSkill()
)

type Skill struct {
	SkillType    SkillType
	Name         string
	AttributeMap AttributeMap
}

func NewSkill(name string, attrs AttributeMap) Skill {
	return Skill{
		Name:         name,
		SkillType:    SkillTypeNormal,
		AttributeMap: attrs,
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

func newSlashSkill() Skill {
	return Skill{
		SkillType: SkillTypeNormal,
		Name:      "slash",
		AttributeMap: NewAttributeMap(
			NewAttribute(AttributeTypeTarget, decimal.NewFromInt(1)),
			NewAttribute(AttributeTypeSDR, decimal.NewFromInt(1)),
		),
	}
}
