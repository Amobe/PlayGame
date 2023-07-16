package database

import (
	"context"
	"errors"
	"fmt"

	"gorm.io/gorm"

	"github.com/Amobe/PlayGame/server/internal/domain/account"
)

var _ account.Repository = (*AccountRepository)(nil)

type AccountGorm struct {
	gorm.Model
	AccountID string `gorm:"column:account_id;type:varchar(64);index"`
	Name      string `gorm:"column:name;type:varchar(64);not null"`
	Email     string `gorm:"column:email;type:varchar(256);not null;uniqueIndex"`
}

func (AccountGorm) TableName() string {
	return "accounts"
}

type AccountRepository struct {
	client *gorm.DB
}

func NewAccountRepository(client *gorm.DB) (*AccountRepository, error) {
	if err := client.AutoMigrate(&AccountGorm{}); err != nil {
		return nil, fmt.Errorf("migrate account gorm: %w", err)
	}
	return &AccountRepository{
		client: client,
	}, nil
}

// Get returns account by id.
func (a *AccountRepository) Get(ctx context.Context, id string) (*account.Account, error) {
	accountGorm := AccountGorm{}
	if err := a.client.WithContext(ctx).
		Where("account_id = ?", id).
		First(&accountGorm).Error; errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, account.ErrAccountNotFound
	} else if err != nil {
		return nil, fmt.Errorf("select account by id: %w", err)
	}
	return account.NewAccount(accountGorm.AccountID, accountGorm.Name, accountGorm.Email), nil
}

// GetByEmail returns account by email.
func (a *AccountRepository) GetByEmail(ctx context.Context, email string) (*account.Account, error) {
	accountGorm := AccountGorm{}
	if err := a.client.WithContext(ctx).
		Where("email = ?", email).
		First(&accountGorm).Error; errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, account.ErrAccountNotFound
	} else if err != nil {
		return nil, fmt.Errorf("select account by email: %w", err)
	}
	return account.NewAccount(accountGorm.AccountID, accountGorm.Name, accountGorm.Email), nil
}

func (a *AccountRepository) Save(ctx context.Context, account *account.Account) error {
	accountGorm := AccountGorm{
		AccountID: account.ID,
		Name:      account.Name,
		Email:     account.Email,
	}
	if err := a.client.WithContext(ctx).
		Save(&accountGorm).Error; err != nil {
		return fmt.Errorf("save account: %w", err)
	}
	return nil
}

type AccountProviderGorm struct {
	gorm.Model
	AccountID  string `gorm:"column:account_id;type:varchar(64);index;not null"`
	Provider   string `gorm:"column:provider;type:varchar(64);not null"`
	ExternalID string `gorm:"column:external_id;type:varchar(64);index;not null"`
}

func (AccountProviderGorm) TableName() string {
	return "account_providers"
}

type AccountProviderRepository struct {
	client *gorm.DB
}

func NewAccountProviderRepository(client *gorm.DB) (*AccountProviderRepository, error) {
	if err := client.AutoMigrate(&AccountProviderGorm{}); err != nil {
		return nil, fmt.Errorf("migrate account provider gorm: %w", err)
	}
	return &AccountProviderRepository{
		client: client,
	}, nil
}

func (a *AccountProviderRepository) List(ctx context.Context, accountID string) ([]*account.AccountProvider, error) {
	var accountProviderGorms []*AccountProviderGorm
	if err := a.client.WithContext(ctx).
		Where("account_id = ?", accountID).
		Find(&accountProviderGorms).Error; err != nil {
		return nil, fmt.Errorf("select account providers: %w", err)
	}
	accountProviders := make([]*account.AccountProvider, len(accountProviderGorms))
	for i, accountProviderGorm := range accountProviderGorms {
		providerType, err := account.GetProviderType(accountProviderGorm.Provider)
		if err != nil {
			return nil, fmt.Errorf("get provider type: %w", err)
		}
		accountProviders[i] = account.NewAccountProvider(accountProviderGorm.AccountID, providerType, accountProviderGorm.ExternalID)
	}
	return accountProviders, nil
}

func (a *AccountProviderRepository) Save(ctx context.Context, accountProvider *account.AccountProvider) error {
	accountProviderGorm := AccountProviderGorm{
		AccountID:  accountProvider.AccountID,
		Provider:   accountProvider.Provider.String(),
		ExternalID: accountProvider.ExternalID,
	}
	if err := a.client.WithContext(ctx).
		Save(&accountProviderGorm).Error; err != nil {
		return fmt.Errorf("save account provider: %w", err)
	}
	return nil
}
