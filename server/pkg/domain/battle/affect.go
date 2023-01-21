package battle

import (
	"github.com/Amobe/PlayGame/server/pkg/domain/vo"
)

type Affect struct {
	ActorID    string
	TargetID   string
	Skill      string
	Attributes []vo.Attribute
}

func NewAffect(actorID, targetID, skillName string, attrs []vo.Attribute) Affect {
	return Affect{
		ActorID:    actorID,
		TargetID:   targetID,
		Skill:      skillName,
		Attributes: attrs,
	}
}
