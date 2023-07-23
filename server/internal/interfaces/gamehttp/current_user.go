package gamehttp

import (
	"context"
	"fmt"

	"github.com/gofiber/fiber/v2"

	"github.com/Amobe/PlayGame/server/internal/domain/account"
	"github.com/Amobe/PlayGame/server/internal/utils"
)

type CurrentUserHandler struct {
	accountRepo account.Repository
}

func NewCurrentUserHandler(repoDeps FiberServerRepoDeps) *CurrentUserHandler {
	return &CurrentUserHandler{
		accountRepo: repoDeps.AccountRepo(),
	}
}

func (h *CurrentUserHandler) FiberHandleCurrentUser(ctx *fiber.Ctx) error {
	tokenPayload := ctx.Locals(FiberLocalKeyTokenPayload).(utils.TokenPayload)
	accountEntity, err := h.handleCurrentUser(ctx.Context(), tokenPayload.AccountID)
	if err != nil {
		return fmt.Errorf("handle current user: %w", err)
	}
	return ctx.JSON(accountEntity)
}

func (h *CurrentUserHandler) handleCurrentUser(ctx context.Context, accountID string) (*account.Account, error) {
	accountEntity, err := h.accountRepo.Get(ctx, accountID)
	if err != nil {
		return nil, fmt.Errorf("account repository get: %w", err)
	}
	return accountEntity, nil
}
