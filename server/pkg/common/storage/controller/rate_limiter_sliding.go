package controller

import (
	"context"

	"github.com/openimsdk/wiseengage/v1/pkg/common/storage/cache"
)

type RateLimiterDatabase interface {
	// Allow check if access is allowed
	Allow(ctx context.Context, key string) (bool, error)
}

type rateLimiterDatabase struct {
	ca cache.RateLimiterSlidingCache
}

func NewRateLimiterDatabase(ca cache.RateLimiterSlidingCache) RateLimiterDatabase {
	return &rateLimiterDatabase{ca: ca}
}

func (r *rateLimiterDatabase) Allow(ctx context.Context, key string) (bool, error) {
	return r.ca.Allow(ctx, key)
}
