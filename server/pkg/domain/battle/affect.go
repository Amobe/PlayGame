package battle

import (
	"github.com/Amobe/PlayGame/server/pkg/domain/vo"
)

type Affect struct {
	ActorIdx   GroundIdx
	TargetIdx  GroundIdx
	Skill      string
	Attributes []vo.Attribute
}

func NewAffect(actorIdx, targetIdx GroundIdx, skillType string, attrs []vo.Attribute) Affect {
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
