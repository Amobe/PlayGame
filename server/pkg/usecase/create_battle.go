package usecase

import (
	"fmt"

	"github.com/Amobe/PlayGame/server/pkg/domain/battle"
	"github.com/Amobe/PlayGame/server/pkg/domain/stage"
	"github.com/Amobe/PlayGame/server/pkg/utils"
)

type CreateBattleInput struct {
	StageID string
}

type CreateBattleOutput struct {
	Battle battle.Battle
}

type CreateBattleUsecase struct {
	stageRepo  stage.Repository
	battleRepo battle.Repository
}

func NewCreateBattleUsecase(
	stageRepo stage.Repository, battleRepo battle.Repository,
) *CreateBattleUsecase {
	return &CreateBattleUsecase{
		stageRepo:  stageRepo,
		battleRepo: battleRepo,
	}
}

func (u *CreateBattleUsecase) Execute(in CreateBattleInput) (out CreateBattleOutput, err error) {
	s, err := u.stageRepo.Get(in.StageID)
	if err != nil {
		err = fmt.Errorf("stage repository get: %w", err)
		return
	}

	battleID := utils.NewUUID()
	allyMinions := battle.NewAllyMinions(nil)
	enemyMinions := battle.NewEnemyMinions(s.Characters)
	minionSlot := battle.NewMinionSlot(allyMinions, enemyMinions)
	b, err := battle.CreateBattle(battleID, minionSlot)
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
		Battle: b,
	}
	return
}
