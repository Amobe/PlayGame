package stage

import "github.com/Amobe/PlayGame/server/pkg/domain/battle"

type Stage struct {
	StageID string
	Mob     battle.Mob
}

func (s Stage) ID() string {
	return s.StageID
}
