package vo

var SkillEmpty = NewSkill(SkillTypeEmpty)

type SkillUseFn func(actorAttr, targetAttr AttributeMap) (actorAffect, targetAffect []Attribute)

type Skill struct {
	SkillType    SkillType
	AttributeMap AttributeMap
	UseFn        SkillUseFn
}

func NewSkill(skillType SkillType, attrs ...Attribute) Skill {
	return Skill{
		SkillType:    skillType,
		AttributeMap: NewAttributeMap(attrs...),
	}
}

func NewCustomSkill(useFn SkillUseFn) Skill {
	return Skill{
		SkillType: SkillTypeCustom,
		UseFn:     useFn,
	}
}

func (s Skill) Use(actorAttr, targetAttr AttributeMap) (actorAffect, targetAffect []Attribute) {
	if s.UseFn == nil {
		return nil, nil
	}
	return s.UseFn(actorAttr, targetAttr)
}

func (s Skill) IsEmpty() bool {
	return s.SkillType == SkillTypeEmpty
}

type SkillType string

func (s SkillType) String() string {
	return string(s)
}

const (
	SkillTypeCustom SkillType = "custom"
	SkillTypeEmpty  SkillType = "empty"
	SkillTypeSlash  SkillType = "slash"
	SkillTypeBlock  SkillType = "block"
	SkillTypeStab   SkillType = "stab"
	SkillTypeHack   SkillType = "hack"
)
