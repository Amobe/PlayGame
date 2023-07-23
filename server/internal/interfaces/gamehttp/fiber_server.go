package gamehttp

import (
	"net"

	fjwt "github.com/gofiber/contrib/jwt"
	"github.com/gofiber/fiber/v2"
	flogger "github.com/gofiber/fiber/v2/middleware/logger"

	"github.com/Amobe/PlayGame/server/internal/domain/account"
	"github.com/Amobe/PlayGame/server/internal/infra/config"
)

const (
	FiberContextKeyToken = "token"
)

type FiberServerConfigDeps interface {
	GoogleAuthConfig() config.GoogleAuth
	TokenConfig() config.Token
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
	server.Use(fjwt.New(fjwt.Config{
		ContextKey: FiberContextKeyToken,
		Filter:     jwtRouteFilter,
		SigningKey: fjwt.SigningKey{
			Key: []byte(configDeps.TokenConfig().JWTSecret),
		},
	}))

	server.Get("/healthcheck", func(ctx *fiber.Ctx) error {
		return ctx.SendString("OK")
	})

	oAuthGoogleHandler := NewOAuthGoogleHandler(configDeps, repoDeps)
	server.Get("/oauth/google", oAuthGoogleHandler.FiberHandleOAuth)
	server.Get("/auth/google/callback", oAuthGoogleHandler.FiberHandleOAuthCallback)

	currentUserHandler := NewCurrentUserHandler(repoDeps)
	server.Get("/session/user", currentUserHandler.FiberHandleCurrentUser)

	return &FiberServer{
		server: server,
	}
}

func (s *FiberServer) Serve(listener net.Listener) error {
	return s.server.Listener(listener)
}

var routePathSkipJWT = map[string]interface{}{
	"/healthcheck":          nil,
	"/oauth/google":         nil,
	"/auth/google/callback": nil,
}

func jwtRouteFilter(ctx *fiber.Ctx) bool {
	// Skip jwt authentication if route is in routePathSkipJWT.
	if _, ok := routePathSkipJWT[ctx.Path()]; ok {
		return true
	}
	return false
}
