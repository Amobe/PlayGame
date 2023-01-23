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
	CharacterID string
	StageID     string
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
	c := vo.NewCharacter(in.CharacterID, []vo.Attribute{
		vo.NewAttribute(vo.AttributeTypeHP, decimal.NewFromInt(500)),
		vo.NewAttribute(vo.AttributeTypeATK, decimal.NewFromInt(10)),
	}...)
	s, err := u.stageRepo.Get(in.StageID)
	if err != nil {
		err = fmt.Errorf("stage repository get: %w", err)
		return
	}

	battleID := utils.NewUUID()
	b, err := battle.CreateBattle(battleID, c, s.Fighter, s.Slot)
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
