package usecase

import (
	"fmt"

	"github.com/shopspring/decimal"

	"github.com/Amobe/PlayGame/server/pkg/domain/battle"
	"github.com/Amobe/PlayGame/server/pkg/domain/stage"
	"github.com/Amobe/PlayGame/server/pkg/domain/vo"
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
	allyMinions := battle.NewAllyMinions([]vo.Character{
		vo.NewCharacterWithSkill("a1", vo.SkillSlash,
			vo.NewAttribute(vo.AttributeTypeATK, decimal.NewFromInt(10)),
			vo.NewAttribute(vo.AttributeTypeHit, decimal.NewFromInt(100)),
		),
		vo.NewCharacter("a2"),
		vo.NewCharacter("a3"),
		vo.NewCharacter("a4"),
		vo.NewCharacter("a5"),
		vo.NewCharacter("a6"),
	})
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
