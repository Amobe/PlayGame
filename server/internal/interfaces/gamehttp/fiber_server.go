package gamehttp

import (
	"net"

	"github.com/gofiber/fiber/v2"
	fcookie "github.com/gofiber/fiber/v2/middleware/encryptcookie"
	flogger "github.com/gofiber/fiber/v2/middleware/logger"

	"github.com/Amobe/PlayGame/server/internal/domain/account"
	"github.com/Amobe/PlayGame/server/internal/domain/battle"
	"github.com/Amobe/PlayGame/server/internal/domain/stage"
	"github.com/Amobe/PlayGame/server/internal/infra/config"
)

const (
	FiberCookieKeyToken       = "token"
	FiberLocalKeyTokenPayload = "token_payload"
)

type FiberServerConfigDeps interface {
	GoogleAuthConfig() config.GoogleAuth
	TokenConfig() config.Token
}

type FiberServerRepoDeps interface {
	AccountRepo() account.Repository
	AccountProviderRepo() account.ProviderRepository
	StageRepo() stage.Repository
	BattleRepo() battle.Repository
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
	server.Use(fcookie.New(fcookie.Config{
		Key: "AA07249FBB99ACF4BD5C6877B2C65C7A",
	}))
	server.Use(newFiberMiddlewareJWT(configDeps))

	server.Get("/healthcheck", func(ctx *fiber.Ctx) error {
		return ctx.SendString("OK")
	})

	oAuthGoogleHandler := NewOAuthGoogleHandler(configDeps, repoDeps)
	server.Get("/oauth/google", oAuthGoogleHandler.FiberHandleOAuth)
	server.Get("/auth/google/callback", oAuthGoogleHandler.FiberHandleOAuthCallback)

	currentUserHandler := NewCurrentUserHandler(repoDeps)
	server.Get("/session/user", currentUserHandler.FiberHandleCurrentUser)

	battleHandler := NewBattleHandler(repoDeps)
	server.Get("/battle/:battle_id", battleHandler.FiberHandleGetBattle)
	server.Post("/battle", battleHandler.FiberHandleCreateBattle)
	server.Post("/battle/:battle_id/fight", battleHandler.FiberHandleFight)

	return &FiberServer{
		server: server,
	}
}

func (s *FiberServer) Serve(listener net.Listener) error {
	return s.server.Listener(listener)
}
