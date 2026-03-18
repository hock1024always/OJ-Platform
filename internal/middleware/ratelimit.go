package middleware

import (
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
)

// tokenBucket 令牌桶，每个 IP 独立一个
type tokenBucket struct {
	tokens     float64
	capacity   float64
	refillRate float64 // tokens per second
	lastRefill time.Time
	mu         sync.Mutex
}

func (b *tokenBucket) allow() bool {
	b.mu.Lock()
	defer b.mu.Unlock()

	now := time.Now()
	elapsed := now.Sub(b.lastRefill).Seconds()
	b.lastRefill = now

	// 补充令牌
	b.tokens += elapsed * b.refillRate
	if b.tokens > b.capacity {
		b.tokens = b.capacity
	}

	if b.tokens >= 1 {
		b.tokens--
		return true
	}
	return false
}

// ipLimiter 全局 IP 令牌桶集合
type ipLimiter struct {
	mu       sync.Mutex
	buckets  map[string]*tokenBucket
	capacity float64
	rate     float64
}

func newIPLimiter(capacity float64, rate float64) *ipLimiter {
	l := &ipLimiter{
		buckets:  make(map[string]*tokenBucket),
		capacity: capacity,
		rate:     rate,
	}
	// 定期清理长期不活跃的 bucket（每 10 分钟）
	go func() {
		ticker := time.NewTicker(10 * time.Minute)
		for range ticker.C {
			l.cleanup()
		}
	}()
	return l
}

func (l *ipLimiter) get(ip string) *tokenBucket {
	l.mu.Lock()
	defer l.mu.Unlock()

	if b, ok := l.buckets[ip]; ok {
		return b
	}
	b := &tokenBucket{
		tokens:     l.capacity,
		capacity:   l.capacity,
		refillRate: l.rate,
		lastRefill: time.Now(),
	}
	l.buckets[ip] = b
	return b
}

func (l *ipLimiter) cleanup() {
	l.mu.Lock()
	defer l.mu.Unlock()
	cutoff := time.Now().Add(-30 * time.Minute)
	for ip, b := range l.buckets {
		b.mu.Lock()
		idle := b.lastRefill.Before(cutoff)
		b.mu.Unlock()
		if idle {
			delete(l.buckets, ip)
		}
	}
}

// 全局通用限流器：100 req/s，峰值 200
var globalLimiter = newIPLimiter(200, 100)

// 提交专用限流器：每 IP 5 req/s，峰值 10（防止恶意刷题）
var submitLimiter = newIPLimiter(10, 5)

// RateLimit 全局 API 限流中间件（令牌桶，100 req/s per IP）
func RateLimit() gin.HandlerFunc {
	return func(c *gin.Context) {
		ip := c.ClientIP()
		if !globalLimiter.get(ip).allow() {
			c.JSON(http.StatusTooManyRequests, gin.H{
				"code":    429,
				"message": "Too Many Requests — please slow down",
			})
			c.Abort()
			return
		}
		c.Next()
	}
}

// SubmitRateLimit 提交/测试接口专用限流（5 req/s per IP，防止代码炸弹刷满 worker）
func SubmitRateLimit() gin.HandlerFunc {
	return func(c *gin.Context) {
		ip := c.ClientIP()
		if !submitLimiter.get(ip).allow() {
			c.JSON(http.StatusTooManyRequests, gin.H{
				"code":    429,
				"message": "Submit rate limit exceeded — max 5 submissions/s per IP",
			})
			c.Abort()
			return
		}
		c.Next()
	}
}
