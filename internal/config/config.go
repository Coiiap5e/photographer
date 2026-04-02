package config

import (
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/Coiiap5e/photographer/internal/errors"
	"github.com/joho/godotenv"
)

type Config struct {
	Server  ServerConfig
	DB      DbConfig
	Kafka   KafkaConfig
	Metrics MetricsConfig
}

type ServerConfig struct {
	Host         string
	Port         int
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
	IdleTimeout  time.Duration
}

type DbConfig struct {
	Host     string
	Port     int
	Username string
	Password string
	Database string
}

type KafkaConfig struct {
	BrokerURLs        []string
	NotificationTopic string
}

type MetricsConfig struct {
	FilePath       string
	ExportInterval time.Duration
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
		MaxOpenConns:    10,
		MaxIdleConns:    3,
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

func Load() (*Config, error) {
	if err := godotenv.Load(); err != nil {
		return nil, errors.Wrap(err, errors.ErrCodeConfig, "error loading .env file")
	}

	dbConfig, err := loadDBConfig()
	if err != nil {
		return nil, errors.Wrap(err, errors.ErrCodeConfig, "failed to load DB config")
	}

	serverConfig, err := loadServerConfig()
	if err != nil {
		return nil, errors.Wrap(err, errors.ErrCodeConfig, "failed to load server config")
	}

	kafkaConfig, err := loadKafkaConfig()
	if err != nil {
		return nil, errors.Wrap(err, errors.ErrCodeConfig, "failed to load kafka config")
	}

	metricsConfig, err := loadMetricsConfig()
	if err != nil {
		return nil, errors.Wrap(err, errors.ErrCodeConfig, "failed to load metrics config")
	}

	return &Config{
		Server:  serverConfig,
		DB:      *dbConfig,
		Kafka:   *kafkaConfig,
		Metrics: *metricsConfig,
	}, nil
}

func loadMetricsConfig() (*MetricsConfig, error) {
	filePath := getEnv("METRICS_FILE_PATH", "metrics.log")
	exportIntervalStr := getEnv("METRICS_EXPORT_INTERVAL_SECONDS", "10")

	exportIntervalSeconds, err := strconv.Atoi(exportIntervalStr)
	if err != nil {
		return nil, errors.New(errors.ErrCodeConfig, "invalid METRICS_EXPORT_INTERVAL_SECONDS")
	}

	return &MetricsConfig{
		FilePath:       filePath,
		ExportInterval: time.Duration(exportIntervalSeconds) * time.Second,
	}, nil
}

func loadServerConfig() (ServerConfig, error) {
	port, err := strconv.Atoi(os.Getenv("APP_SERVER_PORT"))
	if err != nil {
		return ServerConfig{}, errors.New(errors.ErrCodeConfig, "invalid APP_SERVER_PORT")
	}

	readTimeout, _ := strconv.Atoi(os.Getenv("APP_SERVER_READ_TIMEOUT"))
	writeTimeout, _ := strconv.Atoi(os.Getenv("APP_SERVER_WRITE_TIMEOUT"))
	idleTimeout, _ := strconv.Atoi(os.Getenv("APP_SERVER_IDLE_TIMEOUT"))

	return ServerConfig{
		Host:         getEnv("APP_SERVER_HOST", "0.0.0.0"),
		Port:         port,
		ReadTimeout:  time.Duration(readTimeout) * time.Second,
		WriteTimeout: time.Duration(writeTimeout) * time.Second,
		IdleTimeout:  time.Duration(idleTimeout) * time.Second,
	}, nil
}

func loadDBConfig() (*DbConfig, error) {
	if err := godotenv.Load(); err != nil {
		return &DbConfig{}, errors.New(
			errors.ErrCodeConfig, "failed to load .env file",
		)
	}

	if os.Getenv("APP_DB_HOST") == "" {
		return &DbConfig{}, errors.New(
			errors.ErrCodeConfig, "APP_DB_HOST is required",
		)
	}

	port, err := strconv.Atoi(os.Getenv("APP_DB_PORT"))
	if err != nil {
		return &DbConfig{}, errors.New(
			errors.ErrCodeConfig, "APP_DB_PORT must be an integer",
		)
	}

	//Checking required fields

	if os.Getenv("APP_DB_USER") == "" {
		return &DbConfig{}, errors.New(
			errors.ErrCodeConfig, "APP_DB_USER is required",
		)
	}

	if os.Getenv("APP_DB_PASS") == "" {
		return &DbConfig{}, errors.New(
			errors.ErrCodeConfig, "APP_DB_PASS is required",
		)
	}

	if os.Getenv("APP_DB_NAME") == "" {
		return &DbConfig{}, errors.New(
			errors.ErrCodeConfig, "APP_DB_NAME is required",
		)
	}

	return &DbConfig{
		Host:     getEnv("APP_DB_HOST", "localhost"),
		Port:     port,
		Username: os.Getenv("APP_DB_USER"),
		Password: os.Getenv("APP_DB_PASS"),
		Database: os.Getenv("APP_DB_NAME"),
	}, nil
}

func loadKafkaConfig() (*KafkaConfig, error) {
	brokerURLsStr := os.Getenv("KAFKA_BROKER_URLS")
	if brokerURLsStr == "" {
		return nil, errors.New(errors.ErrCodeConfig, "KAFKA_BROKER_URLS is required")
	}
	brokerURLs := strings.Split(brokerURLsStr, ",")

	notificationTopic := os.Getenv("KAFKA_NOTIFICATION_TOPIC")
	if notificationTopic == "" {
		return nil, errors.New(errors.ErrCodeConfig, "KAFKA_NOTIFICATION_TOPIC is required")
	}

	return &KafkaConfig{
		BrokerURLs:        brokerURLs,
		NotificationTopic: notificationTopic,
	}, nil
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
