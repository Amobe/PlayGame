package vo

type Affect struct {
	ActorIdx   GroundIdx
	TargetIdx  GroundIdx
	Skill      string
	Attributes []Attribute
}

func NewAffect(actorIdx, targetIdx GroundIdx, skillType string, attrs []Attribute) Affect {
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
