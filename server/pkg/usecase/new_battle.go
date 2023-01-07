package usecase

import (
	"fmt"

	"github.com/Amobe/PlayGame/server/pkg/domain/battle"
	"github.com/Amobe/PlayGame/server/pkg/domain/character"
	"github.com/Amobe/PlayGame/server/pkg/domain/stage"
	"github.com/Amobe/PlayGame/server/pkg/utils"
)

type NewBattleInput struct {
	CharacterID string
	StageID     string
}

type NewBattleOutput struct {
	Battle battle.Battle
}

type NewBattleUsecase struct {
	characterRepo character.Repository
	stageRepo     stage.Repository
	battleRepo    battle.Repository
}

func NewNewBattleUsecase(
	characterRepo character.Repository, stageRepo stage.Repository, battleRepo battle.Repository,
) *NewBattleUsecase {
	return &NewBattleUsecase{
		characterRepo: characterRepo,
		stageRepo:     stageRepo,
		battleRepo:    battleRepo,
	}
}

func (u *NewBattleUsecase) Execute(in NewBattleInput) (out NewBattleOutput, err error) {
	c, err := u.characterRepo.Get(in.CharacterID)
	if err != nil {
		err = fmt.Errorf("character repository get: %w", err)
		return
	}
	s, err := u.stageRepo.Get(in.StageID)
	if err != nil {
		err = fmt.Errorf("stage repository get: %w", err)
		return
	}

	battleID := utils.NewUUID()
	b := battle.NewBattle(battleID, c, s.Mob)
	err = u.battleRepo.Create(b)
	if err != nil {
		err = fmt.Errorf("battle repository create: %w", err)
	}

	out = NewBattleOutput{
		Battle: b,
	}
	return
}
