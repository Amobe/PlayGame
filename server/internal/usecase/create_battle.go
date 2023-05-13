package usecase

import (
	"fmt"

	"github.com/shopspring/decimal"

	"github.com/Amobe/PlayGame/server/internal/domain/battle"
	"github.com/Amobe/PlayGame/server/internal/domain/stage"
	"github.com/Amobe/PlayGame/server/internal/domain/vo"
	"github.com/Amobe/PlayGame/server/internal/utils"
)

type CreateBattleInput struct {
	StageID string
}

type CreateBattleOutput struct {
	Battle *battle.Battle
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

	allyCamp, err := vo.NewCamp(
		vo.NewCharacterWithSkill(1, vo.SkillSlash,
			vo.NewAttributeMap(
				vo.NewAttribute(vo.AttributeTypeATK, decimal.NewFromInt(20)),
				vo.NewAttribute(vo.AttributeTypeHit, decimal.NewFromInt(100)),
			),
		),
		vo.NewCharacter(2, vo.NewAttributeMap()),
		vo.NewCharacter(3, vo.NewAttributeMap()),
		vo.NewCharacter(4, vo.NewAttributeMap()),
		vo.NewCharacter(5, vo.NewAttributeMap()),
		vo.NewCharacter(6, vo.NewAttributeMap()),
	)
	if err != nil {
		err = fmt.Errorf("new ally camp: %w", err)
		return
	}

	enemyCamp, err := vo.NewCamp(s.Characters...)
	if err != nil {
		err = fmt.Errorf("new enemy camp: %w", err)
		return
	}

	ground := vo.NewGround(allyCamp, enemyCamp)
	minionSlot := battle.NewMinionSlot(ground)
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
