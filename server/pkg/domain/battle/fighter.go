package battle

import "github.com/Amobe/PlayGame/server/pkg/domain/character"

type Fighter interface {
	Affect(attr []character.Attribute)
	UseSkill(skill character.Skill, targetAttr character.AttributeTypeMap) (targetAffect []character.Attribute)

	ID() string
	Alive() bool
	AttributeMap() character.AttributeTypeMap
}
