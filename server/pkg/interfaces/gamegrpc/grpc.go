package gamegrpc

import (
	gamev1 "github.com/Amobe/PlayGame/server/gen/proto/go/game/v1"
	"github.com/Amobe/PlayGame/server/pkg/domain/battle"
	"github.com/Amobe/PlayGame/server/pkg/domain/vo"
)

func BatchAffectToPB(al []battle.Affect) []*gamev1.FightAffect {
	res := make([]*gamev1.FightAffect, 0, len(al))
	for _, a := range al {
		res = append(res, AffectToPB(a))
	}
	return res
}

func AffectToPB(a battle.Affect) *gamev1.FightAffect {
	return &gamev1.FightAffect{
		ActorId:    a.ActorID,
		TargetId:   a.TargetID,
		Skill:      a.Skill,
		Attributes: BatchAttributeToPB(a.Attributes),
	}
}

func BatchAttributeToPB(al []vo.Attribute) []*gamev1.Attribute {
	res := make([]*gamev1.Attribute, 0, len(al))
	for _, a := range al {
		res = append(res, AttributeToPB(a))
	}
	return res
}

func AttributeToPB(a vo.Attribute) *gamev1.Attribute {
	return &gamev1.Attribute{
		Type:  string(a.Type),
		Value: a.Value.String(),
	}
}
