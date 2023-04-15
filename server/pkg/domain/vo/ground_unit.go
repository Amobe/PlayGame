package vo

type GroundUnit struct {
	groundIdx GroundIdx
	Character
}

func NewGroundUnit(groundIdx GroundIdx, c Character) GroundUnit {
	return GroundUnit{
		groundIdx: groundIdx,
		Character: c,
	}
}

func (u GroundUnit) GetGroundIdx() GroundIdx {
	return u.groundIdx
}

func (u GroundUnit) TakeAffect(attributes AttributeMap) GroundUnit {
	return GroundUnit{
		groundIdx: u.groundIdx,
		Character: u.Character.TakeAffect(attributes),
	}
}
