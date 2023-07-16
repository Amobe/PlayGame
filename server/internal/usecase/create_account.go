package usecase

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"

	"github.com/Amobe/PlayGame/server/internal/domain/account"
)

type CreateAccountIn struct {
	Name         string
	Email        string
	ProviderType account.ProviderType
	ExternalID   string
}

type CreateAccountOut struct {
	AccountID string
}

type CreateAccountUseCase struct {
	accountRepository         account.Repository
	accountProviderRepository account.ProviderRepository
}

func NewCreateAccountUseCase(
	accountRepository account.Repository,
	accountProviderRepository account.ProviderRepository,
) *CreateAccountUseCase {
	return &CreateAccountUseCase{
		accountRepository:         accountRepository,
		accountProviderRepository: accountProviderRepository,
	}
}

func (u *CreateAccountUseCase) Execute(ctx context.Context, in CreateAccountIn) (CreateAccountOut, error) {
	if _, err := u.accountRepository.GetByEmail(ctx, in.Email); err == nil {
		return CreateAccountOut{}, fmt.Errorf("account already exists")
	} else if !errors.Is(err, account.ErrAccountNotFound) {
		return CreateAccountOut{}, fmt.Errorf("account repository get by email: %w", err)
	}

	newAccountID := uuid.New().String()
	newAccount := account.NewAccount(newAccountID, in.Name, in.Email)
	newAccountProvider := account.NewAccountProvider(newAccount.ID, in.ProviderType, in.ExternalID)

	if err := u.accountRepository.Save(ctx, newAccount); err != nil {
		return CreateAccountOut{}, fmt.Errorf("account repository save: %w", err)
	}
	if err := u.accountProviderRepository.Save(ctx, newAccountProvider); err != nil {
		return CreateAccountOut{}, fmt.Errorf("account provider repository save: %w", err)
	}
	return CreateAccountOut{
		AccountID: newAccount.ID,
	}, nil
}
