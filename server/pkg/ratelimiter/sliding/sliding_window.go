package sliding

import (
	"context"

	"github.com/openimsdk/wiseengage/v1/pkg/common/storage/controller"
	"github.com/openimsdk/wiseengage/v1/pkg/ratelimiter"
)

type RateLimiter struct {
	keyPrefix string
	db        controller.RateLimiterDatabase
}

func NewRateLimiter(keyPrefix string, db controller.RateLimiterDatabase) ratelimiter.RateLimiter {
	return &RateLimiter{
		keyPrefix: keyPrefix,
		db:        db,
	}
}

func (r *RateLimiter) Allow(ctx context.Context, key string) (bool, error) {
	return r.db.Allow(ctx, r.keyPrefix+":"+key)
}
