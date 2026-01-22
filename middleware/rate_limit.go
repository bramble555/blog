package middleware

import (
	"fmt"
	"sync"
	"time"

	config "github.com/bramble555/blog/conf"
	"github.com/bramble555/blog/controller"
	"github.com/bramble555/blog/global"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis"
)

type tokenBucket struct {
	capacity float64
	rate     float64
	tokens   float64
	last     time.Time
	mu       sync.Mutex
}

func newTokenBucket(rate float64, capacity int) *tokenBucket {
	now := time.Now()
	return &tokenBucket{
		capacity: float64(capacity),
		rate:     rate,
		tokens:   float64(capacity),
		last:     now,
	}
}

func (tb *tokenBucket) allow() bool {
	if tb == nil || tb.rate <= 0 || tb.capacity <= 0 {
		return true
	}
	tb.mu.Lock()
	defer tb.mu.Unlock()
	now := time.Now()
	elapsed := now.Sub(tb.last).Seconds()
	tb.tokens = minFloat(tb.capacity, tb.tokens+elapsed*tb.rate)
	tb.last = now
	if tb.tokens < 1 {
		return false
	}
	tb.tokens -= 1
	return true
}

type ipBucket struct {
	bucket   *tokenBucket
	lastSeen time.Time
}

type IPLimitStore interface {
	Allow(ip string) (bool, error)
	Close() error
}

type memoryIPStore struct {
	mu              sync.Mutex
	buckets         map[string]*ipBucket
	rate            float64
	burst           int
	ttl             time.Duration
	cleanupInterval time.Duration
	lastCleanup     time.Time
}

func newMemoryIPStore(rate float64, burst int, ttl, cleanupInterval time.Duration) *memoryIPStore {
	return &memoryIPStore{
		buckets:         make(map[string]*ipBucket),
		rate:            rate,
		burst:           burst,
		ttl:             ttl,
		cleanupInterval: cleanupInterval,
		lastCleanup:     time.Now(),
	}
}

func (s *memoryIPStore) Allow(ip string) (bool, error) {
	if s == nil || s.rate <= 0 || s.burst <= 0 {
		return true, nil
	}
	if ip == "" {
		return true, nil
	}
	now := time.Now()
	s.mu.Lock()
	if s.cleanupInterval > 0 && now.Sub(s.lastCleanup) >= s.cleanupInterval {
		for k, v := range s.buckets {
			if now.Sub(v.lastSeen) >= s.ttl {
				delete(s.buckets, k)
			}
		}
		s.lastCleanup = now
	}
	entry := s.buckets[ip]
	if entry == nil {
		entry = &ipBucket{
			bucket:   newTokenBucket(s.rate, s.burst),
			lastSeen: now,
		}
		s.buckets[ip] = entry
	} else {
		entry.lastSeen = now
	}
	s.mu.Unlock()
	allowed := entry.bucket.allow()
	return allowed, nil
}

func (s *memoryIPStore) Close() error {
	return nil
}

type redisIPStore struct {
	client *redis.Client
	limit  int
	window time.Duration
	script *redis.Script
}

func newRedisIPStore(client *redis.Client, limit int, window time.Duration) *redisIPStore {
	script := redis.NewScript(`
local current = redis.call("INCR", KEYS[1])
if current == 1 then
  redis.call("PEXPIRE", KEYS[1], ARGV[1])
end
if current > tonumber(ARGV[2]) then
  return 0
end
return 1
`)
	return &redisIPStore{
		client: client,
		limit:  limit,
		window: window,
		script: script,
	}
}

func (s *redisIPStore) Allow(ip string) (bool, error) {
	if s == nil || s.client == nil || s.limit <= 0 || s.window <= 0 {
		return true, nil
	}
	if ip == "" {
		return true, nil
	}
	key := fmt.Sprintf("rate:ip:%s", ip)
	res, err := s.script.Run(s.client, []string{key}, int64(s.window/time.Millisecond), s.limit).Result()
	if err != nil {
		return false, err
	}
	allowed, ok := res.(int64)
	if !ok {
		return false, fmt.Errorf("invalid redis result")
	}
	return allowed == 1, nil
}

func (s *redisIPStore) Close() error {
	return nil
}

type rateLimiter struct {
	global  *tokenBucket
	ipStore IPLimitStore
}

func newRateLimiter(cfg config.RateLimit) *rateLimiter {
	normalized := normalizeRateLimitConfig(cfg)
	var globalBucket *tokenBucket
	if normalized.GlobalRPS > 0 && normalized.GlobalBurst > 0 {
		globalBucket = newTokenBucket(normalized.GlobalRPS, normalized.GlobalBurst)
	}
	ipStore := buildIPStore(normalized)
	return &rateLimiter{
		global:  globalBucket,
		ipStore: ipStore,
	}
}

func buildIPStore(cfg config.RateLimit) IPLimitStore {
	if cfg.IPRPS <= 0 || cfg.IPBurst <= 0 {
		return &memoryIPStore{}
	}
	window := time.Second
	if cfg.Store == "redis" && global.Redis != nil {
		return newRedisIPStore(global.Redis, cfg.IPBurst, window)
	}
	ttl := time.Duration(cfg.IPTTLSeconds) * time.Second
	if ttl <= 0 {
		ttl = 60 * time.Second
	}
	cleanupInterval := time.Duration(cfg.CleanupIntervalSeconds) * time.Second
	if cleanupInterval <= 0 {
		cleanupInterval = 30 * time.Second
	}
	return newMemoryIPStore(cfg.IPRPS, cfg.IPBurst, ttl, cleanupInterval)
}

func normalizeRateLimitConfig(cfg config.RateLimit) config.RateLimit {
	if cfg.GlobalRPS <= 0 {
		cfg.GlobalRPS = 5
	}
	if cfg.GlobalBurst <= 0 {
		cfg.GlobalBurst = 5
	}
	if cfg.IPRPS <= 0 {
		cfg.IPRPS = 1
	}
	if cfg.IPBurst <= 0 {
		cfg.IPBurst = 1
	}
	if cfg.Store == "" {
		cfg.Store = "memory"
	}
	if cfg.IPTTLSeconds <= 0 {
		cfg.IPTTLSeconds = 60
	}
	if cfg.CleanupIntervalSeconds <= 0 {
		cfg.CleanupIntervalSeconds = 30
	}
	return cfg
}

func RateLimitMiddleware() gin.HandlerFunc {
	if global.Config == nil {
		return func(c *gin.Context) {
			c.Next()
		}
	}
	cfg := normalizeRateLimitConfig(global.Config.RateLimit)
	if !cfg.Enable {
		return func(c *gin.Context) {
			c.Next()
		}
	}
	limiter := newRateLimiter(cfg)
	return func(c *gin.Context) {
		if limiter.global != nil && !limiter.global.allow() {
			ip := c.ClientIP()
			global.Log.Warnf("rate_limit global reject ip=%s method=%s path=%s", ip, c.Request.Method, c.Request.URL.Path)
			controller.ResponseError(c, controller.CodeTooManyRequests)
			c.Abort()
			return
		}
		if limiter.ipStore != nil {
			ip := c.ClientIP()
			allowed, err := limiter.ipStore.Allow(ip)
			if err != nil {
				global.Log.Warnf("rate_limit ip check error ip=%s err=%s", ip, err.Error())
				controller.ResponseError(c, controller.CodeTooManyRequests)
				c.Abort()
				return
			}
			if !allowed {
				global.Log.Warnf("rate_limit ip reject ip=%s method=%s path=%s", ip, c.Request.Method, c.Request.URL.Path)
				controller.ResponseError(c, controller.CodeTooManyRequests)
				c.Abort()
				return
			}
		}
		c.Next()
	}
}

func minFloat(a, b float64) float64 {
	if a < b {
		return a
	}
	return b
}
