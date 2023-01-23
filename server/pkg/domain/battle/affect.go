package battle

import (
	"github.com/Amobe/PlayGame/server/pkg/domain/vo"
)

type Affect struct {
	ActorID    string
	TargetID   string
	ChangerID  string
	Skill      string
	Attributes []vo.Attribute
}

func NewAffect(actorID, targetID, changerID, skillName string, attrs []vo.Attribute) Affect {
	return Affect{
		ActorID:    actorID,
		TargetID:   targetID,
		ChangerID:  changerID,
		Skill:      skillName,
		Attributes: attrs,
	}
}
