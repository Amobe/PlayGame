package account

import (
	"context"
	"fmt"
)

var (
	ErrAccountNotFound = fmt.Errorf("account not found")
)

type Repository interface {
	Get(ctx context.Context, id string) (*Account, error)
	GetByEmail(ctx context.Context, email string) (*Account, error)
	Save(ctx context.Context, a *Account) error
}

type ProviderRepository interface {
	List(ctx context.Context, accountID string) ([]*AccountProvider, error)
	Save(ctx context.Context, a *AccountProvider) error
}
