package database

import "fmt"

func GetDSN(config Config) string {
	return fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable TimeZone=Asia/Taipei",
		config.DBHost, config.DBPort, config.DBUser, config.DBPass, config.DBName)
}
