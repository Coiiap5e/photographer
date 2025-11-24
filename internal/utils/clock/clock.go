package clock

import (
	"sync"
	"time"
)

type Clock struct {
	mu       sync.RWMutex
	now      time.Time
	location *time.Location
}

func New(initialTime time.Time) *Clock {
	return &Clock{
		now: initialTime,
	}
}

func (c *Clock) Now() time.Time {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.now
}

func (c *Clock) After(d time.Duration) time.Time {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.now = c.now.Add(d)
	return c.now
}

func (c *Clock) Before(d time.Duration) time.Time {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.now = c.now.Add(-d)
	return c.now
}

func NewInMoscow() *Clock {
	location, _ := time.LoadLocation("Europe/Moscow")
	return &Clock{
		now:      time.Now().In(location),
		location: location,
	}
}
