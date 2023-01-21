package inmem

import (
	"github.com/shopspring/decimal"

	"github.com/Amobe/PlayGame/server/pkg/domain/battle"
	"github.com/Amobe/PlayGame/server/pkg/domain/character"
	"github.com/Amobe/PlayGame/server/pkg/domain/stage"
	"github.com/Amobe/PlayGame/server/pkg/domain/vo"
)

var _ character.Repository = &CharacterRepository{}

type CharacterRepository = inmemStorage[*character.Character]

func NewInmemCharacterRepository() *CharacterRepository {
	s := newInmemStorage[*character.Character]()
	attrs := []vo.Attribute{
		vo.NewAttribute(vo.AttributeTypeHP, decimal.NewFromInt(500)),
		vo.NewAttribute(vo.AttributeTypeATK, decimal.NewFromInt(10)),
	}
	attrMap := vo.NewAttributeTypeMap()
	attrMap.Insert(attrs...)
	fakeCharacter := &character.Character{
		CharacterID: "hero",
		Basement:    attrMap,
	}
	s.Create(fakeCharacter)
	return s
}

var _ stage.Repository = &StageRepository{}

type StageRepository = inmemStorage[*stage.Stage]

func NewInmemStageRepository() *StageRepository {
	s := newInmemStorage[*stage.Stage]()
	fakeStage := &stage.Stage{
		StageID: "fake",
		Fighter: character.RandomCharacter("fake_character"),
		Slot:    battle.NewSlot(character.NewSkillPoisonHit()),
	}
	s.Create(fakeStage)
	return s
}

var _ battle.Repository = &BattleRepository{}

type BattleRepository = inmemEventStorage[battle.Battle]

func NewInmemBattleRepository() *BattleRepository {
	return newInmemEventStorage[battle.Battle](battle.AggregatorLoader)
}
