package battle

import "github.com/Amobe/PlayGame/server/pkg/domain/character"

// Mob is the enemy with a fixed skill to attack the user.
type Mob struct {
	character.Character
	s character.Skill
}

func NewMob(c character.Character, s character.Skill) Mob {
	return Mob{
		Character: c,
		s:         s,
	}
}

func (m *Mob) UseSkill(targetAttr character.AttributeTypeMap) (targetAffect []character.Attribute) {
	return m.Character.UseSkill(m.s, targetAttr)
}