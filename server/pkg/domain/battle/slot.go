package battle

import (
	"github.com/Amobe/PlayGame/server/pkg/domain/vo"
)

type Slot struct {
	Skills [5]vo.Skill
}

func NewSlot(skills ...vo.Skill) Slot {
	slot := Slot{}
	for i := 0; i < 5; i++ {
		slot.Skills[i] = vo.SkillEmpty
	}
	for i, s := range skills {
		slot.Skills[i] = s
	}
	return slot
}
