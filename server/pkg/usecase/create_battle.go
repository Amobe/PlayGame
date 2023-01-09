package usecase

import (
	"fmt"

	"github.com/Amobe/PlayGame/server/pkg/domain/battle"
	"github.com/Amobe/PlayGame/server/pkg/domain/character"
	"github.com/Amobe/PlayGame/server/pkg/domain/stage"
	"github.com/Amobe/PlayGame/server/pkg/utils"
)

type CreateBattleInput struct {
	CharacterID string
	StageID     string
}

type CreateBattleOutput struct {
	Battle battle.Battle
}

type CreateBattleUsecase struct {
	characterRepo character.Repository
	stageRepo     stage.Repository
	battleRepo    battle.Repository
}

func NewCreateBattleUsecase(
	characterRepo character.Repository, stageRepo stage.Repository, battleRepo battle.Repository,
) *CreateBattleUsecase {
	return &CreateBattleUsecase{
		characterRepo: characterRepo,
		stageRepo:     stageRepo,
		battleRepo:    battleRepo,
	}
}

func (u *CreateBattleUsecase) Execute(in CreateBattleInput) (out CreateBattleOutput, err error) {
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

	var mobs []battle.Fighter
	for _, m := range s.Mobs {
		mobs = append(mobs, m)
	}

	battleID := utils.NewUUID()
	b, err := battle.CreateBattle(battleID, c, mobs)
	if err != nil {
		err = fmt.Errorf("create battle: %w", err)
		return
	}
	err = u.battleRepo.Save(b)
	if err != nil {
		err = fmt.Errorf("battle repository create: %w", err)
		return
	}

	out = CreateBattleOutput{
		Battle: *b,
	}
	return
}
