package cache

import "context"

type RateLimiterSlidingCache interface {
	Allow(ctx context.Context, key string) (bool, error)
}
