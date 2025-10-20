package config

import (
	"os"
	"strconv"

	"github.com/Coiiap5e/photographer/internal/errors"
	"github.com/joho/godotenv"
)

type DbConfig struct {
	Host     string `json:"host"`
	Port     int    `json:"port"`
	Username string `json:"username"`
	Password string `json:"password"`
	Database string `json:"database"`
}

func LoadDBConfig() (DbConfig, error) {
	if err := godotenv.Load(); err != nil {
		return DbConfig{}, errors.New(
			errors.ErrCodeConfig, "failed to load .env file",
		)
	}

	if os.Getenv("APP_DB_HOST") == "" {
		return DbConfig{}, errors.New(
			errors.ErrCodeConfig, "APP_DB_HOST is required",
		)
	}

	port, err := strconv.Atoi(os.Getenv("APP_DB_PORT"))
	if err != nil {
		return DbConfig{}, errors.New(
			errors.ErrCodeConfig, "APP_DB_PORT must be an integer",
		)
	}

	//Checking required fields

	if os.Getenv("APP_DB_USER") == "" {
		return DbConfig{}, errors.New(
			errors.ErrCodeConfig, "APP_DB_USER is required",
		)
	}

	if os.Getenv("APP_DB_PASS") == "" {
		return DbConfig{}, errors.New(
			errors.ErrCodeConfig, "APP_DB_PASS is required",
		)
	}

	if os.Getenv("APP_DB_NAME") == "" {
		return DbConfig{}, errors.New(
			errors.ErrCodeConfig, "APP_DB_NAME is required",
		)
	}

	return DbConfig{
		Host:     getEnv("APP_DB_HOST", "localhost"),
		Port:     port,
		Username: os.Getenv("APP_DB_USER"),
		Password: os.Getenv("APP_DB_PASS"),
		Database: os.Getenv("APP_DB_NAME"),
	}, nil
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
