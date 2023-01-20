package battle

import (
	"github.com/Amobe/PlayGame/server/pkg/domain/skill"
)

type Slot struct {
	Skills [5]skill.Skill
}

func NewSlot(skills ...skill.Skill) Slot {
	slot := Slot{}
	for i, s := range skills {
		if s == nil {
			s = skill.NewSkillEmpty()
		}
		slot.Skills[i] = s
	}
	return slot
}
