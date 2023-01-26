package vo

var SkillEmpty = Skill{
	SkillType: SkillTypeEmpty,
}

type SkillUseFn func(actorAttr, targetAttr AttributeMap) (actorAffect, targetAffect []Attribute)

type Skill struct {
	SkillID      string
	SkillType    SkillType
	Name         string
	AttributeMap AttributeMap
	UseFn        SkillUseFn
}

func NewSkill(skillID string, name string, attrs ...Attribute) Skill {
	return Skill{
		SkillID:      skillID,
		Name:         name,
		SkillType:    SkillTypeNormal,
		AttributeMap: NewAttributeMap(attrs...),
	}
}

func NewCustomSkill(useFn SkillUseFn) Skill {
	return Skill{
		Name:      "custom_skill",
		SkillType: SkillTypeNormal,
		UseFn:     useFn,
	}
}

func (s Skill) ID() string {
	return s.SkillID
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
	SkillTypeNormal SkillType = "normal"
	SkillTypeEmpty  SkillType = "empty"
)
