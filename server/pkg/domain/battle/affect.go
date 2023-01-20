package battle

import (
	"github.com/Amobe/PlayGame/server/pkg/domain/valueobject"
)

type Affect struct {
	ActorID    string
	TargetID   string
	Skill      string
	Attributes []valueobject.Attribute
}

func NewAffect(actorID, targetID, skillName string, attrs []valueobject.Attribute) Affect {
	return Affect{
		ActorID:    actorID,
		TargetID:   targetID,
		Skill:      skillName,
		Attributes: attrs,
	}
}
