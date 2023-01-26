package inmem

import (
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
		Fighter: vo.RandomCharacter("fake_character"),
		Slot:    battle.NewSlot(),
	}
	s.Create(fakeStage)
	return s
}

var _ battle.Repository = &BattleRepository{}

type BattleRepository = inmemEventStorage[battle.Battle]

func NewInmemBattleRepository() *BattleRepository {
	return newInmemEventStorage[battle.Battle](battle.AggregatorLoader)
}
