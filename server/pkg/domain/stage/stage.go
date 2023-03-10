package stage

import (
	"github.com/Amobe/PlayGame/server/pkg/domain/battle"
	"github.com/Amobe/PlayGame/server/pkg/utils/domain"
)

var _ domain.Aggregator = &Stage{}

type coreAggregator = domain.CoreAggregator

type Stage struct {
	coreAggregator
	StageID string
	Fighter battle.Fighter
	Slot    battle.Slot
}

func (s Stage) ID() string {
	return s.StageID
}
