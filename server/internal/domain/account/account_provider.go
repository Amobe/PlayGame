package account

import "fmt"

type AccountProvider struct {
	AccountID  string
	Provider   ProviderType
	ExternalID string
}

func NewAccountProvider(accountID string, provider ProviderType, externalID string) *AccountProvider {
	return &AccountProvider{
		AccountID:  accountID,
		Provider:   provider,
		ExternalID: externalID,
	}
}

type ProviderType string

func GetProviderType(provider string) (ProviderType, error) {
	switch provider {
	case ProviderTypeGoogle.String():
		return ProviderTypeGoogle, nil
	}
	return ProviderTypeUnknown, fmt.Errorf("unknown provider: %s", provider)
}

func (p ProviderType) String() string {
	return string(p)
}

var (
	ProviderTypeUnknown ProviderType = "unknown"
	ProviderTypeGoogle  ProviderType = "google"
)
