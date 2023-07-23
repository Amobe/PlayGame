package google

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/Amobe/PlayGame/server/internal/interfaces/gamehttp"
)

type Client struct{}

func NewClient() *Client {
	return &Client{}
}

type userInformationResponse struct {
	SUB           string `json:"sub"`
	Name          string `json:"name"`
	GivenName     string `json:"given_name"`
	FamilyName    string `json:"family_name"`
	Picture       string `json:"picture"`
	Email         string `json:"email"`
	EmailVerified bool   `json:"email_verified"`
	Locale        string `json:"locale"`
}

func (a *Client) GetUserInformation(accessToken string) (*gamehttp.UserInformation, error) {
	getUserInfoUrl := fmt.Sprintf("https://www.googleapis.com/oauth2/v3/userinfo?access_token=%s", accessToken)

	req, err := http.NewRequest(http.MethodGet, getUserInfoUrl, nil)
	if err != nil {
		return nil, fmt.Errorf("new user info request: %w", err)
	}

	client := http.Client{
		Timeout: time.Second * 30,
	}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("get request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("get request not ok: %s", resp.Status)
	}
	var body bytes.Buffer
	_, err = body.ReadFrom(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("read body: %w", err)
	}

	var userInfoResp map[string]interface{}
	if err := json.Unmarshal(body.Bytes(), &userInfoResp); err != nil {
		return nil, fmt.Errorf("unmarshal body: %w", err)
	}
	if userInfoResp["error"] != nil {
		return nil, fmt.Errorf("user info error: %s", userInfoResp["error"])
	}

	var userInfo userInformationResponse
	if err := json.Unmarshal(body.Bytes(), &userInfo); err != nil {
		return nil, fmt.Errorf("unmarshal user info: %w", err)
	}

	return &gamehttp.UserInformation{
		ID:            userInfo.SUB,
		Email:         userInfo.Email,
		VerifiedEmail: userInfo.EmailVerified,
		Name:          userInfo.Name,
		GivenName:     userInfo.GivenName,
		FamilyName:    userInfo.FamilyName,
		Locale:        userInfo.Locale,
	}, nil
}
