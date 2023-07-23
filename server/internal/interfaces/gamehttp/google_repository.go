package gamehttp

type UserInformation struct {
	ID            string
	Email         string
	VerifiedEmail bool
	Name          string
	GivenName     string
	FamilyName    string
	Locale        string
}

type GoogleRepository interface {
	GetUserInformation(accessToken string) (*UserInformation, error)
}
