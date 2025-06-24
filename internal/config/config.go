package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	DBHost     string
	DBPort     string
	DBUser     string
	DBPassword string
	DBName     string

	APIGenderURL string
	APIAgeURL    string
	APINationURL string

	LogLevel string
}

func LoadConfig() Config {
	err := godotenv.Load()
	if err != nil {
		log.Println("Warning: .env file not found, reading config from environment variables")
	}

	return Config{
		DBHost:       getEnv("DB_HOST", "localhost"),
		DBPort:       getEnv("DB_PORT", "5432"),
		DBUser:       getEnv("DB_USER", "postgres"),
		DBPassword:   getEnv("DB_PASSWORD", ""),
		DBName:       getEnv("DB_NAME", "peopledb"),
		APIGenderURL: getEnv("API_GENDER_URL", "https://api.genderapi.io"),
		APIAgeURL:    getEnv("API_AGE_URL", "https://api.agify.io"),
		APINationURL: getEnv("API_NATION_URL", "https://api.nationalize.io"),
		LogLevel:     getEnv("LOG_LEVEL", "debug"),
	}
}

func getEnv(key, fallback string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return fallback
}
