package currency

import (
	"context"
	"sync"
)

type InMemoryCacheService struct {
	rate float64
	mu   sync.RWMutex
}

func NewInMemoryCacheService() *InMemoryCacheService {
	return &InMemoryCacheService{}
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
}
