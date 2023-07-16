package account

type Account struct {
	ID    string
	Name  string
	Email string
}

func NewAccount(id string, name string, email string) *Account {
	return &Account{
		ID:    id,
		Name:  name,
		Email: email,
	}
}
