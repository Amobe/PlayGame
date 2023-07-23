package usecase

import (
	"context"
	"fmt"

	"github.com/Amobe/PlayGame/server/internal/domain/account"
)

type GetAccountInput struct {
	AccountID string
}

type GetAccountOutput struct {
	Account *account.Account
}

type GetAccountUsecase struct {
	accountRepo account.Repository
}

func NewGetAccountUsecase(accountRepo account.Repository) *GetAccountUsecase {
	return &GetAccountUsecase{
		accountRepo: accountRepo,
	}
}

func (u *GetAccountUsecase) Execute(ctx context.Context, in GetAccountInput) (GetAccountOutput, error) {
	accountEntity, err := u.accountRepo.Get(ctx, in.AccountID)
	if err != nil {
		return GetAccountOutput{}, fmt.Errorf("account repository get: %w", err)
	}
	return GetAccountOutput{
		Account: accountEntity,
	}, nil
}
