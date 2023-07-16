package usecase

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/Amobe/PlayGame/server/internal/domain/account"
)

func TestCreateAccountUseCase_Execute(t *testing.T) {
	t.Run("DuplicatedEmail", testCreateAccountUseCase_Execute_DuplicatedEmail)
	t.Run("Success", testCreateAccountUseCase_Execute_Success)
}

type fakeAccountRepository struct {
	getFn        func(ctx context.Context, id string) (*account.Account, error)
	getByEmailFn func(ctx context.Context, email string) (*account.Account, error)
	saveFn       func(ctx context.Context, a *account.Account) error
}

func (f *fakeAccountRepository) Get(ctx context.Context, id string) (*account.Account, error) {
	if f.getFn != nil {
		return f.getFn(ctx, id)
	}
	return nil, nil
}

func (f *fakeAccountRepository) GetByEmail(ctx context.Context, email string) (*account.Account, error) {
	if f.getByEmailFn != nil {
		return f.getByEmailFn(ctx, email)
	}
	return nil, nil
}

func (f *fakeAccountRepository) Save(ctx context.Context, a *account.Account) error {
	if f.saveFn != nil {
		return f.saveFn(ctx, a)
	}
	return nil
}

type fakeAccountProviderRepository struct {
	listFn func(ctx context.Context, accountID string) ([]*account.AccountProvider, error)
	saveFn func(ctx context.Context, a *account.AccountProvider) error
}

func (f *fakeAccountProviderRepository) List(ctx context.Context, accountID string) ([]*account.AccountProvider, error) {
	if f.listFn != nil {
		return f.listFn(ctx, accountID)
	}
	return nil, nil
}

func (f *fakeAccountProviderRepository) Save(ctx context.Context, a *account.AccountProvider) error {
	if f.saveFn != nil {
		return f.saveFn(ctx, a)
	}
	return nil
}

func testCreateAccountUseCase_Execute_DuplicatedEmail(t *testing.T) {
	// create a use case with fake account repository
	// the repository will return an account without error
	// in this case, the use case should return an account existed error
	uc := CreateAccountUseCase{
		accountRepository: &fakeAccountRepository{
			getByEmailFn: func(ctx context.Context, email string) (*account.Account, error) {
				return &account.Account{
					ID: "fake-id",
				}, nil
			},
		},
		accountProviderRepository: nil,
	}
	_, err := uc.Execute(context.Background(), CreateAccountIn{})
	assert.EqualError(t, err, "account already exists")
}

func testCreateAccountUseCase_Execute_Success(t *testing.T) {
	uc := CreateAccountUseCase{
		accountRepository: &fakeAccountRepository{
			getByEmailFn: func(ctx context.Context, email string) (*account.Account, error) {
				return nil, account.ErrAccountNotFound
			},
		},
		accountProviderRepository: &fakeAccountProviderRepository{
			saveFn: func(ctx context.Context, ap *account.AccountProvider) error {
				return nil
			},
		},
	}
	out, err := uc.Execute(context.Background(), CreateAccountIn{
		Name:  "fake-name",
		Email: "fake-email",
	})
	assert.NoError(t, err)
	assert.True(t, len(out.AccountID) > 0)
}
