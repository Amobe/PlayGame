package main

import (
	"errors"
	"fmt"
	"net"
	"sync"

	"golang.org/x/exp/slog"
	"google.golang.org/grpc"

	gamev1 "github.com/Amobe/PlayGame/server/gen/proto/go/game/v1"
	"github.com/Amobe/PlayGame/server/internal/domain/battle"
	"github.com/Amobe/PlayGame/server/internal/domain/stage"
	"github.com/Amobe/PlayGame/server/internal/infra/config"
	"github.com/Amobe/PlayGame/server/internal/infra/database"
	"github.com/Amobe/PlayGame/server/internal/infra/google"
	"github.com/Amobe/PlayGame/server/internal/infra/inmem"
	"github.com/Amobe/PlayGame/server/internal/interfaces/gamegrpc"
	"github.com/Amobe/PlayGame/server/internal/interfaces/gamehttp"
)

func main() {
	if err := run(); err != nil {
		slog.Error("run %s", err)
	}
}

func run() error {
	serverConfig := config.Server{
		GrpcHost: "",
		GrpcPort: 8080,
		HttpHost: "",
		HttpPort: 8081,
	}

	listenOn := fmt.Sprintf("%s:%d", serverConfig.GrpcHost, serverConfig.GrpcPort)
	listener, err := net.Listen("tcp", listenOn)
	if err != nil {
		return fmt.Errorf("listen on %s: %w", listenOn, err)
	}

	httpListenOn := fmt.Sprintf("%s:%d", serverConfig.HttpHost, serverConfig.HttpPort)
	httpListener, err := net.Listen("tcp", httpListenOn)
	if err != nil {
		return fmt.Errorf("listen on %s: %w", httpListenOn, err)
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

	googleAuthConfig, err := config.LoadGoogleAuthConfig()
	if err != nil {
		return fmt.Errorf("load google auth config: %w", err)
	}
	googleClient := google.NewClient()

	deps := deps{
		stageRepo:  inmem.NewInmemStageRepository(),
		battleRepo: battleGormRepo,
	}

	var wg sync.WaitGroup
	wg.Add(2)

	errCh := make(chan error, 2)

	go func() {
		slog.Info("Listening grpc on %s", listenOn)
		grpcServer := grpc.NewServer()
		handler := gamegrpc.NewGameServiceHandler(deps)
		gamev1.RegisterGameServiceServer(grpcServer, handler)
		if err := grpcServer.Serve(listener); err != nil {
			errCh <- fmt.Errorf("serve gRPC server: %w", err)
		}
		wg.Done()
	}()

	go func() {
		slog.Info("Listening http on %s", httpListenOn)
		httpServer := gamehttp.NewFiberServer(serverConfig, googleAuthConfig, googleClient)
		if err := httpServer.Serve(httpListener); err != nil {
			errCh <- fmt.Errorf("serve http server: %w", err)
		}
		wg.Done()
	}()

	wg.Wait()
	close(errCh)

	var serveErrors []error
	for err := range errCh {
		serveErrors = append(serveErrors, err)
	}
	if len(serveErrors) > 0 {
		return errors.Join(serveErrors...)
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
