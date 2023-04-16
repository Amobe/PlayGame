package inmem

import (
	"github.com/shopspring/decimal"

	"github.com/Amobe/PlayGame/server/internal/domain/battle"
	"github.com/Amobe/PlayGame/server/internal/domain/stage"
	"github.com/Amobe/PlayGame/server/internal/domain/vo"
)

var _ stage.Repository = &StageRepository{}

type StageRepository = InmemStorage[*stage.Stage]

func NewInmemStageRepository() *StageRepository {
	s := NewInmemStorage[*stage.Stage]()
	fakeStage := &stage.Stage{
		StageID: "fake",
		Characters: []vo.Character{
			vo.NewCharacter(1, vo.NewAttributeMap(vo.NewAttribute(vo.AttributeTypeHP, decimal.NewFromInt(100)))),
			vo.NewCharacter(2, vo.NewAttributeMap()),
			vo.NewCharacter(3, vo.NewAttributeMap()),
			vo.NewCharacter(4, vo.NewAttributeMap()),
			vo.NewCharacter(5, vo.NewAttributeMap()),
			vo.NewCharacter(6, vo.NewAttributeMap()),
		},
	}
	s.Create(fakeStage)
	return s
}

var _ battle.Repository = &BattleRepository{}

type BattleRepository = inmemEventStorage[*battle.Battle]

func NewInmemBattleRepository() *BattleRepository {
	return newInmemEventStorage[*battle.Battle](battle.AggregatorLoader)
}
