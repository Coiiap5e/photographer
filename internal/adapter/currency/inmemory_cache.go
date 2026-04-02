package currency

import (
	"context"
	"log/slog"
	"os"
	"path/filepath"
	"strconv"
	"sync"
)

type InMemoryCacheService struct {
	rate     float64
	mu       sync.RWMutex
	filePath string
	logger   *slog.Logger
}

func NewInMemoryCacheService(filePath string, logger *slog.Logger) *InMemoryCacheService {
	s := &InMemoryCacheService{
		filePath: filePath,
		logger:   logger,
	}
	s.loadRateFromFile()
	return s
}

func (s *InMemoryCacheService) GetUSDRate(ctx context.Context) (float64, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.rate, nil
}

func (s *InMemoryCacheService) UpdateRate(newRate float64) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.rate = newRate
	s.saveRateToFile()
}

func (s *InMemoryCacheService) loadRateFromFile() {
	s.mu.Lock()
	defer s.mu.Unlock()

	data, err := os.ReadFile(s.filePath)
	if err != nil {
		s.logger.Warn("failed to read cache file, starting with empty cache", "path", s.filePath, "error", err)
		return
	}

	rate, err := strconv.ParseFloat(string(data), 64)
	if err != nil {
		s.logger.Error("failed to parse rate from cache file", "path", s.filePath, "error", err)
		return
	}

	s.rate = rate
	s.logger.Info("successfully loaded currency rate from cache file", "path", s.filePath, "rate", rate)
}

func (s *InMemoryCacheService) saveRateToFile() {
	data := []byte(strconv.FormatFloat(s.rate, 'f', -1, 64))

	// Ensure the directory exists
	dir := filepath.Dir(s.filePath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		s.logger.Error("failed to create cache directory", "path", dir, "error", err)
		return
	}

	if err := os.WriteFile(s.filePath, data, 0644); err != nil {
		s.logger.Error("failed to write rate to cache file", "path", s.filePath, "error", err)
	}
}
