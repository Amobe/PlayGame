package battle

import (
	"github.com/Amobe/PlayGame/server/pkg/domain/vo"
)

type Affect struct {
	ActorID    string
	ActorIdx   GroundIdx
	TargetID   string
	TargetIdx  GroundIdx
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

func NewAffectV2(actorIdx, targetIdx GroundIdx, skillType string, attrs []vo.Attribute) Affect {
	return Affect{
		ActorIdx:   actorIdx,
		TargetIdx:  targetIdx,
		Skill:      skillType,
		Attributes: attrs,
	}
}

func NewMissAffect(actorIdx, targetIdx GroundIdx) Affect {
	return Affect{
		ActorIdx:  actorIdx,
		TargetIdx: targetIdx,
		Skill:     "miss",
	}
}
