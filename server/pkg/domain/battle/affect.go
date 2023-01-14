package battle

import (
	"github.com/Amobe/PlayGame/server/pkg/domain/character"
)

type Affect struct {
	ActorID    string
	TargetID   string
	Skill      string
	Attributes []character.Attribute
}

func NewAffect(actorID, targetID, skillName string, attrs []character.Attribute) Affect {
	return Affect{
		ActorID:    actorID,
		TargetID:   targetID,
		Skill:      skillName,
		Attributes: attrs,
	}
}
