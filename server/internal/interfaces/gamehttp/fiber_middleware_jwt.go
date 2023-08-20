package gamehttp

import (
	"github.com/gofiber/fiber/v2"

	"github.com/Amobe/PlayGame/server/internal/utils"
)

var routePathSkipJWT = map[string]interface{}{
	"/health":               nil,
	"/health/":              nil,
	"/oauth/google":         nil,
	"/auth/google/callback": nil,
}

type fiberMiddlewareJWT struct {
	jwtSecretKey string
}

func newFiberMiddlewareJWT(configDeps FiberServerConfigDeps) fiber.Handler {
	handler := &fiberMiddlewareJWT{
		jwtSecretKey: configDeps.TokenConfig().JWTSecret,
	}
	return handler.handle
}

func (m *fiberMiddlewareJWT) handle(ctx *fiber.Ctx) error {
	if _, ok := routePathSkipJWT[ctx.Path()]; ok {
		return ctx.Next()
	}

	token := ctx.Cookies(FiberCookieKeyToken)
	tokenPayload, err := utils.ValidateToken(token, m.jwtSecretKey)
	if err != nil {
		return fiber.NewError(fiber.StatusUnauthorized, "validate token")
	}
	ctx.Locals(FiberLocalKeyTokenPayload, tokenPayload)
	return ctx.Next()
}
