package cache

import (
	"context"
	"time"

	"github.com/patrickmn/go-cache"
)

const (
	// DefaultExpiration for items
	DefaultExpiration time.Duration = 24 * time.Hour
	cleanupInterval   time.Duration = 1 * time.Hour
)

// Cache interface
type Cache interface {
	Get(string) (interface{}, bool)
	Set(string, interface{}, time.Duration)
}

type memoryCache struct {
	ctx context.Context
	*cache.Cache
}

// NewStatic cache service
func NewStatic(ctx context.Context) Cache {
	return &memoryCache{
		ctx,
		cache.New(DefaultExpiration, cleanupInterval),
	}
}
