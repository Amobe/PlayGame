package battle

import "github.com/Amobe/PlayGame/server/pkg/domain/vo"

//go:generate mockery --name Unit --inpackage
type Unit interface {
	GetGroundIdx() vo.GroundIdx
	GetAttributeMap() vo.AttributeMap
	GetAgi() int
	IsDead() bool
	TakeAffect(affects []vo.Attribute) vo.GroundUnit
	GetSkill() vo.Skill
}
