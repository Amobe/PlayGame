package main

import (
	"fmt"
	"log"
	"net"

	"google.golang.org/grpc"

	gamev1 "github.com/Amobe/PlayGame/server/gen/proto/go/game/v1"
	"github.com/Amobe/PlayGame/server/internal/domain/battle"
	"github.com/Amobe/PlayGame/server/internal/domain/stage"
	"github.com/Amobe/PlayGame/server/internal/infra/database"
	"github.com/Amobe/PlayGame/server/internal/infra/inmem"
	"github.com/Amobe/PlayGame/server/internal/interfaces/gamegrpc"
)

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}

func run() error {
	listenOn := ":8080"
	listener, err := net.Listen("tcp", listenOn)
	if err != nil {
		return fmt.Errorf("listen on %s: %w", listenOn, err)
	}

	dbConfig, err := database.LoadConfig()
	if err != nil {
		return fmt.Errorf("load database config: %w", err)
	}
	dbClient, err := database.NewClient(dbConfig)
	if err != nil {
		return fmt.Errorf("new database client: %w", err)
	}

	battleGormRepo, err := database.NewBattleGormRepository(dbClient)
	if err != nil {
		return fmt.Errorf("new battle gorm repository: %w", err)
	}

	deps := deps{
		stageRepo:  inmem.NewInmemStageRepository(),
		battleRepo: battleGormRepo,
	}

	server := grpc.NewServer()
	handler := gamegrpc.NewGameServiceHandler(deps)
	gamev1.RegisterGameServiceServer(server, handler)
	log.Println("Listening on", listenOn)
	if err := server.Serve(listener); err != nil {
		return fmt.Errorf("serve gRPC server: %w", err)
	}

	return nil
}

type deps struct {
	stageRepo  stage.Repository
	battleRepo battle.Repository
}

func (d deps) StageRepo() stage.Repository {
	return d.stageRepo
}

func (d deps) BattleRepo() battle.Repository {
	return d.battleRepo
}
