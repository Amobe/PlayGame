package main

import (
	"fmt"
	"log"
	"net"

	"google.golang.org/grpc"

	gamev1 "github.com/Amobe/PlayGame/server/gen/proto/go/game/v1"
	"github.com/Amobe/PlayGame/server/pkg/domain/battle"
	"github.com/Amobe/PlayGame/server/pkg/domain/character"
	"github.com/Amobe/PlayGame/server/pkg/domain/stage"
	"github.com/Amobe/PlayGame/server/pkg/infra/inmem"
	"github.com/Amobe/PlayGame/server/pkg/interfaces/gamegrpc"
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

	deps := deps{
		characterRepo: inmem.NewInmemCharacterRepository(),
		stageRepo:     inmem.NewInmemStageRepository(),
		battleRepo:    inmem.NewInmemBattleRepository(),
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
	characterRepo character.Repository
	stageRepo     stage.Repository
	battleRepo    battle.Repository
}

func (d deps) CharacterRepo() character.Repository {
	return d.characterRepo
}

func (d deps) StageRepo() stage.Repository {
	return d.stageRepo
}

func (d deps) BattleRepo() battle.Repository {
	return d.battleRepo
}
