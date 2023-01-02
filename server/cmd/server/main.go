package main

import (
	"context"
	"fmt"
	"log"
	"net"

	"google.golang.org/grpc"

	gamev1 "github.com/Amobe/PlayGame/server/gen/proto/go/game/v1"
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

	server := grpc.NewServer()
	gamev1.RegisterGameServiceServer(server, &GameServiceServer{})
	log.Println("Listening on", listenOn)
	if err := server.Serve(listener); err != nil {
		return fmt.Errorf("serve gRPC server: %w", err)
	}

	return nil
}

type GameServiceServer struct {
	gamev1.UnimplementedGameServiceServer
}

func (s *GameServiceServer) NewBattle(ctx context.Context, req *gamev1.NewBattleRequest) (*gamev1.NewBattleResponse, error) {
	log.Println("NewBattle")
	return &gamev1.NewBattleResponse{}, nil
}

func (s *GameServiceServer) Fight(ctx context.Context, req *gamev1.FightRequest) (*gamev1.FightResponse, error) {
	log.Println("Fight")
	return &gamev1.FightResponse{}, nil
}
