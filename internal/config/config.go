package config

import (
	"os"
	"strconv"
	"time"

	"github.com/Coiiap5e/photographer/internal/errors"
	"github.com/joho/godotenv"
)

type DbConfig struct {
	Host     string
	Port     int
	Username string
	Password string
	Database string
}

type PoolConfig struct {
	MaxOpenConns    int
	MaxIdleConns    int
	MaxConnLifetime time.Duration
	MaxConnIdleTime time.Duration
}

type PoolConfigBuilder struct {
	config PoolConfig
}

func NewPoolConfigBuilder() *PoolConfigBuilder {
	return &PoolConfigBuilder{config: PoolConfig{
		MaxOpenConns:    25,
		MaxIdleConns:    5,
		MaxConnLifetime: 20 * time.Minute,
		MaxConnIdleTime: 1 * time.Minute,
	}}
}

func (b *PoolConfigBuilder) WithMaxOpenConns(MaxOpenConns int) *PoolConfigBuilder {
	if MaxOpenConns > 0 {
		b.config.MaxOpenConns = MaxOpenConns
	}
	return b
}
func (b *PoolConfigBuilder) WithMaxIdleConns(MaxIdleConns int) *PoolConfigBuilder {
	if MaxIdleConns > 0 {
		b.config.MaxIdleConns = MaxIdleConns
	}
	return b
}

func (b *PoolConfigBuilder) WithMaxConnLifetime(MaxConnLifetime time.Duration) *PoolConfigBuilder {
	if MaxConnLifetime > 0 {
		b.config.MaxConnLifetime = MaxConnLifetime
	}
	return b
}

func (b *PoolConfigBuilder) WithMaxConnIdleTime(MaxConnIdleTime time.Duration) *PoolConfigBuilder {
	if MaxConnIdleTime > 0 {
		b.config.MaxConnIdleTime = MaxConnIdleTime
	}
	return b
}

func (b *PoolConfigBuilder) Build() PoolConfig {
	return b.config
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
