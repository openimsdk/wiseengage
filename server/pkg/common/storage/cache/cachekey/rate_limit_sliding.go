package cachekey

const (
	buildRateLimitSlidingKey = "RATE_LIMIT_Sliding:"
)

func GetRateLimitSlidingKey(key string) string {
	return buildRateLimitSlidingKey + key
}
