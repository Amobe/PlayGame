package database

import (
	"fmt"
	"os"
	"strconv"
)

type Config struct {
	DBHost string
	DBPort int
	DBUser string
	DBPass string
	DBName string
}

func LoadConfig() (Config, error) {
	dbHost := os.Getenv("DB_HOST")
	if len(dbHost) == 0 {
		dbHost = "localhost"
	}

	dbPortStr := os.Getenv("DB_PORT")
	dbPort, err := strconv.Atoi(dbPortStr)
	if err != nil {
		return Config{}, fmt.Errorf("parse db port: %w", err)
	}

	dbUser := os.Getenv("DB_USER")
	if len(dbUser) == 0 {
		dbUser = "gorm"
	}

	dbPass := os.Getenv("DB_PASS")
	if len(dbPass) == 0 {
		dbPass = "gorm"
	}

	dbName := os.Getenv("DB_NAME")
	if len(dbName) == 0 {
		dbName = "game_db"
	}

	return Config{
		DBHost: dbHost,
		DBPort: dbPort,
		DBUser: dbUser,
		DBPass: dbPass,
		DBName: dbName,
	}, nil
}
