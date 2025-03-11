package redis

import (
	"context"
	"errors"
	"strconv"
	"time"

	"github.com/openimsdk/tools/errs"
	"github.com/openimsdk/wiseengage/v1/pkg/common/storage/cache"
	"github.com/openimsdk/wiseengage/v1/pkg/common/storage/cache/cachekey"
	"github.com/redis/go-redis/v9"
)

type RateLimitSlidingCache struct {
	rdb    redis.UniversalClient
	limit  int64
	window time.Duration
	expire time.Duration
}

func NewRateLimitSlidingCache(rdb redis.UniversalClient, limit, window int) cache.RateLimiterSlidingCache {
	r := &RateLimitSlidingCache{
		rdb:    rdb,
		limit:  int64(limit),
		window: time.Duration(window) * time.Second,
	}
	r.expire = time.Duration(float64(r.window) * 1.1)
	return r
}

func (r *RateLimitSlidingCache) Allow(ctx context.Context, k string) (bool, error) {
	//	var script = `
	//local key = KEYS[1]
	//local limit = ARGV[1]
	//local window = ARGV[2]
	//local now = ARGV[3]
	//
	//redis.call('ZREMRANGEBYSCORE', key, 0, now-window)
	//local count = redis.call('ZCARD', key)
	//if count >= limit  then
	//	return -1
	//else
	//	redis.call('ZADD', key, now, now)
	//	redis.call('EXPIRE', key, window)
	//	return 0
	//end
	//`

	key := cachekey.GetRateLimitSlidingKey(k)
	now := time.Now()
	start := now.Add(-r.window).UnixMilli()
	nowm := now.UnixMilli()
	if err := r.rdb.ZRemRangeByScore(ctx, key, "0", strconv.Itoa(int(start))).Err(); err != nil && !errors.Is(err, redis.Nil) {
		return false, errs.WrapMsg(err, "ZRemRangeByScore in Allow failed")
	} else if err == nil {
		res, err := r.rdb.ZCard(ctx, key).Result()
		if err != nil && !errors.Is(err, redis.Nil) {
			return false, errs.WrapMsg(err, "ZCard in Allow failed")
		}
		if res >= r.limit {
			return false, nil
		}
	}

	if err := r.rdb.ZAdd(ctx, key, redis.Z{Member: nowm, Score: float64(nowm)}).Err(); err != nil {
		return false, errs.WrapMsg(err, "ZAdd in Allow failed")
	}
	if err := r.rdb.Expire(ctx, key, r.expire).Err(); err != nil {
		return false, errs.WrapMsg(err, "Expire in Allow failed")
	}
	return true, nil
}
