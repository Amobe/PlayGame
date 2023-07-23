package gamehttp

import (
	"net"

	"github.com/gofiber/fiber/v2"
	flogger "github.com/gofiber/fiber/v2/middleware/logger"

	"github.com/Amobe/PlayGame/server/internal/domain/account"
	"github.com/Amobe/PlayGame/server/internal/infra/config"
)

type FiberServerConfigDeps interface {
	GoogleAuthConfig() config.GoogleAuth
}

type FiberServerRepoDeps interface {
	AccountRepo() account.Repository
	AccountProviderRepo() account.ProviderRepository
	GoogleRepo() GoogleRepository
}

type FiberServer struct {
	server *fiber.App
}

func NewFiberServer(
	configDeps FiberServerConfigDeps,
	repoDeps FiberServerRepoDeps,
) *FiberServer {
	server := fiber.New()
	server.Use(flogger.New())

	server.Get("/healthcheck", func(ctx *fiber.Ctx) error {
		return ctx.SendString("OK")
	})

	oAuthGoogleHandler := NewOAuthGoogleHandler(configDeps, repoDeps)
	server.Get("/oauth/google", oAuthGoogleHandler.FiberHandleOAuth)
	server.Get("/auth/google/callback", oAuthGoogleHandler.FiberHandleOAuthCallback)

	return &FiberServer{
		server: server,
	}
}

func (s *FiberServer) Serve(listener net.Listener) error {
	return s.server.Listener(listener)
}
