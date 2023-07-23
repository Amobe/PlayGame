package gamehttp

import (
	"context"
	"fmt"

	"github.com/gofiber/fiber/v2"
	"golang.org/x/exp/slog"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"

	"github.com/Amobe/PlayGame/server/internal/domain/account"
	"github.com/Amobe/PlayGame/server/internal/infra/config"
	"github.com/Amobe/PlayGame/server/internal/usecase"
	"github.com/Amobe/PlayGame/server/internal/utils"
)

type OAuthGoogleHandler struct {
	oAuthConfig     *oauth2.Config
	tokenConfig     config.Token
	googleRepo      GoogleRepository
	createAccountUC *usecase.CreateAccountUseCase
}

func NewOAuthGoogleHandler(
	configDeps FiberServerConfigDeps,
	repoDeps FiberServerRepoDeps,
) *OAuthGoogleHandler {
	googleAuthConfig := configDeps.GoogleAuthConfig()
	googleRepo := repoDeps.GoogleRepo()
	accountRepo := repoDeps.AccountRepo()
	accountProviderRepo := repoDeps.AccountProviderRepo()
	return &OAuthGoogleHandler{
		oAuthConfig: &oauth2.Config{
			ClientID:     googleAuthConfig.ClientID,
			ClientSecret: googleAuthConfig.ClientSecret,
			RedirectURL:  googleAuthConfig.RedirectURL,
			Scopes: []string{
				"https://www.googleapis.com/auth/userinfo.email",
				"https://www.googleapis.com/auth/userinfo.profile",
			},
			Endpoint: google.Endpoint,
		},
		tokenConfig:     configDeps.TokenConfig(),
		googleRepo:      googleRepo,
		createAccountUC: usecase.NewCreateAccountUseCase(accountRepo, accountProviderRepo),
	}
}

func (o *OAuthGoogleHandler) FiberHandleOAuth(ctx *fiber.Ctx) error {
	redirectUrl := o.handleOAuth(ctx.Context())
	return ctx.Redirect(redirectUrl)
}

func (o *OAuthGoogleHandler) handleOAuth(ctx context.Context) string {
	return o.oAuthConfig.AuthCodeURL("not-implemented-yet")
}

func (o *OAuthGoogleHandler) FiberHandleOAuthCallback(ctx *fiber.Ctx) error {
	code := ctx.Query("code")
	jwtToken, err := o.handleOAuthCallback(ctx.Context(), code)
	if err != nil {
		slog.Error("handle oauth callback", err)
		return ctx.SendStatus(fiber.StatusInternalServerError)
	}
	resp := fiber.Map{
		"token": jwtToken,
	}
	return ctx.JSON(resp)
}

func (o *OAuthGoogleHandler) handleOAuthCallback(ctx context.Context, code string) (string, error) {
	oAuthToken, err := o.oAuthConfig.Exchange(ctx, code)
	if err != nil {
		return "", fmt.Errorf("exchange code to token: %w", err)
	}
	profile, err := o.googleRepo.GetUserInformation(oAuthToken.AccessToken)
	if err != nil {
		return "", fmt.Errorf("get user information: %w", err)
	}
	in := usecase.CreateAccountIn{
		Name:         profile.Name,
		Email:        profile.Email,
		ProviderType: account.ProviderTypeGoogle,
		ExternalID:   profile.ID,
	}
	out, err := o.createAccountUC.Execute(ctx, in)
	if err != nil {
		return "", fmt.Errorf("create account: %w", err)
	}
	payload := utils.TokenPayload{
		AccountID: out.AccountID,
		Name:      profile.Name,
	}
	token, err := utils.GenerateToken(o.tokenConfig.ExpiredIn, payload, o.tokenConfig.JWTSecret)
	if err != nil {
		return "", fmt.Errorf("generate token: %w", err)
	}
	return token, nil
}
