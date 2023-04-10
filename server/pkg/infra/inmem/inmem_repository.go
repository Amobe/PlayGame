package inmem

import (
	"github.com/shopspring/decimal"

	"github.com/Amobe/PlayGame/server/pkg/domain/battle"
	"github.com/Amobe/PlayGame/server/pkg/domain/stage"
	"github.com/Amobe/PlayGame/server/pkg/domain/vo"
)

var _ stage.Repository = &StageRepository{}

type StageRepository = InmemStorage[*stage.Stage]

func NewInmemStageRepository() *StageRepository {
	s := NewInmemStorage[*stage.Stage]()
	fakeStage := &stage.Stage{
		StageID: "fake",
		Characters: []vo.Character{
			vo.NewCharacter("e1", vo.NewAttribute(vo.AttributeTypeHP, decimal.NewFromInt(100))),
			vo.NewCharacter("e2"),
			vo.NewCharacter("e3"),
			vo.NewCharacter("e4"),
			vo.NewCharacter("e5"),
			vo.NewCharacter("e6"),
		},
	}
	s.Create(fakeStage)
	return s
}

var _ battle.Repository = &BattleRepository{}

type BattleRepository = inmemEventStorage[battle.Battle]

func NewInmemBattleRepository() *BattleRepository {
	return newInmemEventStorage[battle.Battle](battle.AggregatorLoader)
}
