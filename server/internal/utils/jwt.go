package utils

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type TokenPayload struct {
	AccountID string
	Name      string
}

func GenerateToken(ttl time.Duration, payload TokenPayload, secretJWTKey string) (string, error) {
	now := time.Now().UTC()
	claims := jwt.MapClaims{
		"id":   payload.AccountID,
		"name": payload.Name,
		"exp":  now.Add(ttl).Unix(),
		"iat":  now.Unix(),
		"nbf":  now.Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenString, err := token.SignedString([]byte(secretJWTKey))
	if err != nil {
		return "", fmt.Errorf("signed token string: %w", err)
	}
	return tokenString, nil
}

func ValidateToken(token string, secretJWTKey string) (TokenPayload, error) {
	jwtToken, err := jwt.Parse(token, func(jwtToken *jwt.Token) (interface{}, error) {
		if _, ok := jwtToken.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected method: %s", jwtToken.Header["alg"])
		}
		return []byte(secretJWTKey), nil
	})
	if err != nil {
		return TokenPayload{}, fmt.Errorf("parse token: %w", err)
	}
	return RetrieveTokenPayload(jwtToken)
}

func RetrieveTokenPayload(token interface{}) (TokenPayload, error) {
	jwtToken, ok := token.(*jwt.Token)
	if !ok {
		return TokenPayload{}, fmt.Errorf("invalid token")
	}
	claims, ok := jwtToken.Claims.(jwt.MapClaims)
	if !ok || !jwtToken.Valid {
		return TokenPayload{}, fmt.Errorf("invalid claims")
	}
	return TokenPayload{
		AccountID: claims["id"].(string),
		Name:      claims["name"].(string),
	}, nil
}
