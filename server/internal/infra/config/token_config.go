package config

import (
	"fmt"
	"os"
	"strconv"
	"time"
)

type Token struct {
	ExpiredIn   time.Duration
	MaxAgeInMin int
	JWTSecret   string
}

func LoadTokenConfig() (Token, error) {
	jwtSecret := os.Getenv("TOKEN_JWT_SECRET")
	if len(jwtSecret) == 0 {
		return Token{}, fmt.Errorf("token jwt secret not provided")
	}

	expiredInStr := os.Getenv("TOKEN_EXPIRED_IN")
	if len(expiredInStr) == 0 {
		return Token{}, fmt.Errorf("token expired in not provided")
	}
	expiredIn, err := time.ParseDuration(expiredInStr)
	if err != nil {
		return Token{}, fmt.Errorf("parse expired in duration: %w", err)
	}

	maxAgeInMinStr := os.Getenv("TOKEN_MAX_AGE_IN_MIN")
	if len(maxAgeInMinStr) == 0 {
		return Token{}, fmt.Errorf("token max age in min not provided")
	}
	maxAgeInMin, err := strconv.Atoi(maxAgeInMinStr)
	if err != nil {
		return Token{}, fmt.Errorf("parse max age in min: %w", err)
	}

	return Token{
		ExpiredIn:   expiredIn,
		MaxAgeInMin: maxAgeInMin,
		JWTSecret:   jwtSecret,
	}, nil
}
