package vo

type Affect struct {
	ActorIdx   GroundIdx
	TargetIdx  GroundIdx
	Skill      string
	Attributes AttributeMap
}

func NewAffect(actorIdx, targetIdx GroundIdx, skillType string, attrs AttributeMap) Affect {
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

// EqualTo compares two Affect objects.
func (a Affect) EqualTo(other Affect) bool {
	if a.ActorIdx.EqualTo(other.ActorIdx) {
		return false
	}
	if a.TargetIdx.EqualTo(other.TargetIdx) {
		return false
	}
	if a.Skill != other.Skill {
		return false
	}
	if len(a.Attributes) != len(other.Attributes) {
		return false
	}
	for i := range a.Attributes {
		if a.Attributes[i] != other.Attributes[i] {
			return false
		}
	}
	return true
}
