package inmem

import (
	"github.com/Amobe/PlayGame/server/pkg/domain/battle"
	"github.com/Amobe/PlayGame/server/pkg/domain/character"
	"github.com/Amobe/PlayGame/server/pkg/domain/stage"
)

var _ character.Repository = &CharacterRepository{}

type CharacterRepository = inmemStorage[character.Character]

func NewInmemCharacterRepository() *CharacterRepository {
	s := newInmemStorage[character.Character]()
	fakeCharacter := character.Character{
		CharacterID: "hero",
		Basement: []character.Attribute{
			{Type: character.AttributeTypeHP, Value: "500"},
			{Type: character.AttributeTypeATK, Value: "10"},
		},
	}
	s.Create(fakeCharacter)
	return s
}

var _ stage.Repository = &StageRepository{}

type StageRepository = inmemStorage[stage.Stage]

func NewInmemStageRepository() *StageRepository {
	s := newInmemStorage[stage.Stage]()
	fakeStage := stage.Stage{
		StageID: "fake",
		Mob:     battle.NewMob(character.NewCharacter(), character.NewSkillPoisonHit()),
	}
	fakeStage.Mob.CharacterID = "monster"
	fakeStage.Mob.Basement = append(fakeStage.Mob.Basement, character.Attribute{Type: character.AttributeTypeATK, Value: "5"})
	s.Create(fakeStage)
	return s
}

var _ battle.Repository = &BattleRepository{}

type BattleRepository = inmemStorage[battle.Battle]

func NewInmemBattleRepository() *BattleRepository {
	return newInmemStorage[battle.Battle]()
}
