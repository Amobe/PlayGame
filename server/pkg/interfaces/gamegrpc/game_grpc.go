package gamegrpc

import (
	"context"
	"fmt"
	"log"

	"github.com/Amobe/PlayGame/server/pkg/domain/battle"
	"github.com/Amobe/PlayGame/server/pkg/domain/stage"
	"github.com/Amobe/PlayGame/server/pkg/usecase"

	gamev1 "github.com/Amobe/PlayGame/server/gen/proto/go/game/v1"
)

type GameServiceDeps interface {
	StageRepo() stage.Repository
	BattleRepo() battle.Repository
}

type GameServiceHandler struct {
	gamev1.UnimplementedGameServiceServer

	stageRepo  stage.Repository
	battleRepo battle.Repository
}

func NewGameServiceHandler(deps GameServiceDeps) *GameServiceHandler {
	return &GameServiceHandler{
		stageRepo:  deps.StageRepo(),
		battleRepo: deps.BattleRepo(),
	}
}

func (s *GameServiceHandler) NewBattle(ctx context.Context, req *gamev1.NewBattleRequest) (*gamev1.NewBattleResponse, error) {
	log.Println("NewBattle")
	in := usecase.CreateBattleInput{
		StageID: "fake",
	}
	uc := usecase.NewCreateBattleUsecase(s.stageRepo, s.battleRepo)
	out, err := uc.Execute(in)
	if err != nil {
		return nil, fmt.Errorf("execute new battle usecase: %w", err)
	}
	return &gamev1.NewBattleResponse{
		BattleId: out.Battle.ID(),
	}, nil
}

func (s *GameServiceHandler) Fight(ctx context.Context, req *gamev1.FightRequest) (*gamev1.FightResponse, error) {
	log.Println("Fight")
	in := usecase.BattleFightIn{
		BattleID: req.GetBattleId(),
	}
	uc := usecase.NewBattleFightUsecase(s.battleRepo)
	out, err := uc.Execute(in)
	if err != nil {
		return nil, fmt.Errorf("execute battle fight usecase: %w", err)
	}
	return &gamev1.FightResponse{
		Affects: BatchAffectToPB(out.Affects),
	}, nil
}
