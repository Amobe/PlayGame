package inmem

import (
	"github.com/Amobe/PlayGame/server/pkg/domain/battle"
	"github.com/Amobe/PlayGame/server/pkg/domain/character"
	"github.com/Amobe/PlayGame/server/pkg/domain/skill"
	"github.com/Amobe/PlayGame/server/pkg/domain/stage"
	"github.com/Amobe/PlayGame/server/pkg/domain/valueobject"
)

var _ character.Repository = &CharacterRepository{}

type CharacterRepository = inmemStorage[*character.Character]

func NewInmemCharacterRepository() *CharacterRepository {
	s := newInmemStorage[*character.Character]()
	attrs := []valueobject.Attribute{
		{Type: valueobject.AttributeTypeHP, Value: "500"},
		{Type: valueobject.AttributeTypeATK, Value: "10"},
	}
	attrMap := valueobject.NewAttributeTypeMap()
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
		Slot:    battle.NewSlot(skill.NewSkillPoisonHit()),
	}
	s.Create(fakeStage)
	return s
}

var _ battle.Repository = &BattleRepository{}

type BattleRepository = inmemEventStorage[battle.Battle]

func NewInmemBattleRepository() *BattleRepository {
	return newInmemEventStorage[battle.Battle](battle.AggregatorLoader)
}
