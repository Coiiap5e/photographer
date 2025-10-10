package config

import (
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type DbConfig struct {
	Host     string `json:"host"`
	Port     int    `json:"port"`
	Username string `json:"username"`
	Password string `json:"password"`
	Database string `json:"database"`
}

func LoadDBConfig() DbConfig {
	if err := godotenv.Load(); err != nil {
		log.Println("Error loading .env file")
	}

	port, err := strconv.Atoi(os.Getenv("DB_PORT"))
	if err != nil {
		log.Fatal("APP_DB_PORT must be an integer")
	}

	return DbConfig{
		Host:     os.Getenv("APP_DB_HOST"),
		Port:     port,
		Username: os.Getenv("APP_DB_USER"),
		Password: os.Getenv("APP_DB_PASS"),
		Database: os.Getenv("APP_DB_NAME"),
	}
}
