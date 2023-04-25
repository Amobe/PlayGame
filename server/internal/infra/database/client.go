package database

import (
	"fmt"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func GetDSN(config Config) string {
	return fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable TimeZone=Asia/Taipei",
		config.DBHost, config.DBPort, config.DBUser, config.DBPass, config.DBName)
}

func NewClient(config Config) (*gorm.DB, error) {
	client, err := gorm.Open(postgres.Open(GetDSN(config)), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("connect database with gorm: %w", err)
	}
	return client, nil
}
