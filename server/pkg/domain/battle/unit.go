package battle

import "github.com/Amobe/PlayGame/server/pkg/domain/vo"

//go:generate mockery --name Unit --inpackage
type Unit interface {
	GetGroundIdx() GroundIdx
	GetAttributeMap() vo.AttributeMap
	GetAgi() int
	IsDead() bool
	TakeAffect(affects []vo.Attribute) Unit
	GetSkill() vo.Skill
}

var _ Unit = unit{}

type unit struct {
	groundIdx GroundIdx
	vo.Character
}

func NewUnit(groundIdx GroundIdx, c vo.Character) *unit {
	return &unit{
		groundIdx: groundIdx,
		Character: c,
	}
}

func (u unit) GetGroundIdx() GroundIdx {
	return u.groundIdx
}

func (u unit) TakeAffect(affects []vo.Attribute) Unit {
	c := u.Character.TakeAffect(affects)
	return &unit{
		groundIdx: u.groundIdx,
		Character: c,
	}
}
