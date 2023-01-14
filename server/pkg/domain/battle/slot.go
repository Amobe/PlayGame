package battle

import "github.com/Amobe/PlayGame/server/pkg/domain/character"

type Slot struct {
	Skills [5]character.Skill
}

func NewSlot(skills ...character.Skill) Slot {
	slot := Slot{}
	for i, s := range skills {
		if s == nil {
			s = character.NewSkillEmpty()
		}
		slot.Skills[i] = s
	}
	return slot
}
