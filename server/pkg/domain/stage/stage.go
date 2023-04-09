package stage

import (
	"github.com/Amobe/PlayGame/server/pkg/domain/vo"
	"github.com/Amobe/PlayGame/server/pkg/utils/domain"
)

var _ domain.Aggregator = &Stage{}

type coreAggregator = domain.CoreAggregator

type Stage struct {
	coreAggregator
	StageID    string
	Characters []vo.Character
}

func (s Stage) ID() string {
	return s.StageID
}
