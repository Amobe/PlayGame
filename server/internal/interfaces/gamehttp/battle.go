package gamehttp

import (
	"context"
	"fmt"

	"github.com/gofiber/fiber/v2"
	"golang.org/x/exp/slog"

	"github.com/Amobe/PlayGame/server/internal/domain/battle"
	"github.com/Amobe/PlayGame/server/internal/domain/stage"
	"github.com/Amobe/PlayGame/server/internal/domain/vo"
	"github.com/Amobe/PlayGame/server/internal/usecase"
)

type BattleHandler struct {
	stageRepo  stage.Repository
	battleRepo battle.Repository
}

func NewBattleHandler(repoDeps FiberServerRepoDeps) *BattleHandler {
	return &BattleHandler{
		stageRepo:  repoDeps.StageRepo(),
		battleRepo: repoDeps.BattleRepo(),
	}
}

type CreateBattleResponse struct {
	BattleID string `json:"battle_id"`
}

func (h *BattleHandler) FiberHandleCreateBattle(ctx *fiber.Ctx) error {
	battleID, err := h.handleCreateBattle(context.Background(), "stage1")
	if err != nil {
		return fmt.Errorf("handle create battle: %w", err)
	}
	resp := CreateBattleResponse{
		BattleID: battleID,
	}
	return ctx.JSON(resp)
}

func (h *BattleHandler) handleCreateBattle(ctx context.Context, stageID string) (string, error) {
	in := usecase.CreateBattleInput{
		StageID: stageID,
	}
	uc := usecase.NewCreateBattleUsecase(h.stageRepo, h.battleRepo)
	out, err := uc.Execute(in)
	if err != nil {
		return "", fmt.Errorf("execute create battle usecase: %w", err)
	}
	return out.Battle.ID(), nil
}

type FightResponse struct {
	Affects []vo.Affect `json:"affects"`
}

func (h *BattleHandler) FiberHandleFight(ctx *fiber.Ctx) error {
	battleID := ctx.Params("battle_id")
	affects, err := h.handleFight(context.Background(), battleID)
	if err != nil {
		return fmt.Errorf("handle fight: %w", err)
	}
	resp := FightResponse{
		Affects: affects,
	}
	return ctx.JSON(resp)
}

func (h *BattleHandler) handleFight(ctx context.Context, battleID string) ([]vo.Affect, error) {
	in := usecase.BattleFightIn{
		BattleID: battleID,
	}
	uc := usecase.NewBattleFightUsecase(h.battleRepo)
	out, err := uc.Execute(in)
	if err != nil {
		return nil, fmt.Errorf("execute battle fight usecase: %w", err)
	}
	return out.Affects, nil
}

func (h *BattleHandler) FiberHandleGetBattle(ctx *fiber.Ctx) error {
	battleID := ctx.Params("battle_id")
	slog.Info("FiberHandleGetBattle", "battle_id", battleID)
	b, err := h.battleRepo.Get(battleID)
	if err != nil {
		return fmt.Errorf("handle get: %w", err)
	}
	return ctx.JSON(b)
}
